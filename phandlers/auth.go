package phandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/auth"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/cookie"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/secondary"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/u2flib"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/pritunl/pritunl-zero/validator"
	"strings"
)

func authStateGet(c *gin.Context) {
	data := auth.GetState()

	if demo.IsDemo() {
		provider := &auth.StateProvider{
			Id:    "demo",
			Type:  "demo",
			Label: "demo",
		}
		data.Providers = append(data.Providers, provider)
	}

	c.JSON(200, data)
}

type authData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func authSessionPost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	srvc := c.MustGet("service").(*service.Service)
	data := &authData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if srvc == nil {
		utils.AbortWithStatus(c, 404)
		return
	}

	usr, errData, err := auth.Local(db, data.Username, data.Password)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(401, errData)
		return
	}

	deviceAuth, secProviderId, errData, err := validator.ValidateProxy(
		db, usr, false, srvc, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.ProxyLoginFailed,
			audit.Fields{
				"error":   errData.Error,
				"message": errData.Message,
			},
		)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(401, errData)
		return
	}

	if deviceAuth {
		secd, err := secondary.New(db, usr.Id, secondary.Proxy,
			secondary.DeviceProvider)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		data, err := secd.GetData()
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(201, data)
		return
	} else if secProviderId != "" {
		secd, err := secondary.New(db, usr.Id, secondary.Proxy, secProviderId)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		data, err := secd.GetData()
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(201, data)
		return
	}

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.ProxyLogin,
		audit.Fields{
			"method": "local",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewProxy(srvc, c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.Proxy)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	redirectJson(c, c.Request.URL.Query().Get("redirect_url"))
}

type secondaryData struct {
	Token    string `json:"token"`
	Factor   string `json:"factor"`
	Passcode string `json:"passcode"`
}

func authSecondaryPost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	srvc := c.MustGet("service").(*service.Service)
	data := &secondaryData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	secd, err := secondary.Get(db, data.Token, secondary.Proxy)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			errData := &errortypes.ErrorData{
				Error:   "secondary_expired",
				Message: "Two-factor authentication has expired",
			}
			c.JSON(401, errData)
		} else {
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	errData, err := secd.Handle(db, c.Request, data.Factor, data.Passcode)
	if err != nil {
		if _, ok := err.(*secondary.IncompleteError); ok {
			c.Status(201)
		} else {
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	if errData != nil {
		c.JSON(401, errData)
		return
	}

	usr, err := secd.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_, _, errData, err = validator.ValidateProxy(
		db, usr, false, srvc, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.ProxyLoginFailed,
			audit.Fields{
				"error":   errData.Error,
				"message": errData.Message,
			},
		)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(401, errData)
		return
	}

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.ProxyLogin,
		audit.Fields{
			"method": "secondary",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewProxy(srvc, c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.Proxy)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	redirectJson(c, c.Request.URL.Query().Get("redirect_url"))
}

func logoutGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)

	if authr.IsValid() {
		err := authr.Remove(db)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}
	}

	c.Redirect(302, "/")
}

func authRequestGet(c *gin.Context) {
	auth.Request(c)
}

func authCallbackGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	srvc := c.MustGet("service").(*service.Service)
	sig := c.Query("sig")
	query := strings.Split(c.Request.URL.RawQuery, "&sig=")[0]

	if srvc == nil {
		utils.AbortWithStatus(c, 404)
		return
	}

	usr, tokn, errData, err := auth.Callback(db, sig, query)
	if err != nil {
		switch err.(type) {
		case *auth.InvalidState:
			c.Redirect(302, "/")
			break
		default:
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	if errData != nil {
		c.JSON(401, errData)
		return
	}

	deviceAuth, secProviderId, errData, err := validator.ValidateProxy(
		db, usr, false, srvc, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.ProxyLoginFailed,
			audit.Fields{
				"error":   errData.Error,
				"message": errData.Message,
			},
		)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(401, errData)
		return
	}

	if deviceAuth {
		secd, err := secondary.New(db, usr.Id, secondary.Proxy,
			secondary.DeviceProvider)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		urlQuery, err := secd.GetQuery()
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		if tokn.Query != "" {
			urlQuery += "&" + tokn.Query
		}

		c.Redirect(302, "/login?"+urlQuery)
		return
	} else if secProviderId != "" {
		secd, err := secondary.New(db, usr.Id, secondary.Proxy, secProviderId)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		urlQuery, err := secd.GetQuery()
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		if tokn.Query != "" {
			urlQuery += "&" + tokn.Query
		}

		c.Redirect(302, "/login?"+urlQuery)
		return
	}

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.ProxyLogin,
		audit.Fields{
			"method": "sso",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewProxy(srvc, c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.Proxy)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	redirectQuery(c, tokn.Query)
}

func authU2fSignGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	token := c.Query("token")

	secd, err := secondary.Get(db, token, secondary.Proxy)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			errData := &errortypes.ErrorData{
				Error:   "secondary_expired",
				Message: "Two-factor authentication has expired",
			}
			c.JSON(401, errData)
		} else {
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	resp, errData, err := secd.DeviceSignRequest(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(401, errData)
		return
	}

	c.JSON(200, resp)
}

type u2fSignData struct {
	Token    string               `json:"token"`
	Response *u2flib.SignResponse `json:"response"`
}

func authU2fSignPost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	srvc := c.MustGet("service").(*service.Service)
	data := &u2fSignData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	secd, err := secondary.Get(db, data.Token, secondary.Proxy)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			errData := &errortypes.ErrorData{
				Error:   "secondary_expired",
				Message: "Two-factor authentication has expired",
			}
			c.JSON(401, errData)
		} else {
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	errData, err := secd.DeviceSignResponse(db, data.Response)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	usr, err := secd.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_, secProviderId, errData, err := validator.ValidateProxy(
		db, usr, false, srvc, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.ProxyLoginFailed,
			audit.Fields{
				"error":   errData.Error,
				"message": errData.Message,
			},
		)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(401, errData)
		return
	}

	if secProviderId != "" {
		secd, err := secondary.New(db, usr.Id, secondary.Proxy, secProviderId)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		data, err := secd.GetData()
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(201, data)
		return
	}

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.ProxyLogin,
		audit.Fields{
			"method": "secondary",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewProxy(srvc, c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.Proxy)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.Status(200)
}
