package mhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/auth"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/cookie"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/device"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/secondary"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/u2flib"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/pritunl/pritunl-zero/validator"
	"gopkg.in/mgo.v2/bson"
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
	data := &authData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
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

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.AdminPrimaryApprove,
		audit.Fields{
			"method": "local",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	deviceAuth, secProviderId, errData, err := validator.ValidateAdmin(
		db, usr, false, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.AdminLoginFailed,
			audit.Fields{
				"method":  "local",
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
		deviceCount, err := device.Count(db, usr.Id)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		secType := ""
		var secProvider bson.ObjectId
		if deviceCount == 0 {
			if secProviderId == "" {
				secType = secondary.AdminDeviceRegister
				secProvider = secondary.DeviceProvider
			} else {
				secType = secondary.Admin
				secProvider = secProviderId
			}
		} else {
			secType = secondary.AdminDevice
			secProvider = secondary.DeviceProvider
		}

		secd, err := secondary.New(db, usr.Id, secType, secProvider)
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
		secd, err := secondary.New(db, usr.Id, secondary.Admin, secProviderId)
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
		audit.AdminLogin,
		audit.Fields{
			"method": "local",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewAdmin(c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.Admin)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.Status(200)
}

type secondaryData struct {
	Token    string `json:"token"`
	Factor   string `json:"factor"`
	Passcode string `json:"passcode"`
}

func authSecondaryPost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &secondaryData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	secd, err := secondary.Get(db, data.Token, secondary.Admin)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			errData := &errortypes.ErrorData{
				Error:   "secondary_expired",
				Message: "Secondary authentication has expired",
			}
			c.JSON(401, errData)
		} else {
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	usr, err := secd.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
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
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.AdminLoginFailed,
			audit.Fields{
				"method":      "secondary",
				"provider_id": secd.ProviderId,
				"error":       errData.Error,
				"message":     errData.Message,
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
		audit.AdminSecondaryApprove,
		audit.Fields{
			"provider_id": secd.ProviderId,
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	deviceAuth, _, errData, err := validator.ValidateAdmin(
		db, usr, false, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.AdminLoginFailed,
			audit.Fields{
				"method":  "secondary",
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
		deviceCount, err := device.Count(db, usr.Id)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		if deviceCount == 0 {
			secd, err := secondary.New(db, usr.Id,
				secondary.AdminDeviceRegister, secondary.DeviceProvider)
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
	}

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.AdminLogin,
		audit.Fields{
			"method": "secondary",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewAdmin(c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.Admin)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.Status(200)
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

	usr, _ := authr.GetUser(db)
	if usr != nil {
		err := audit.New(
			db,
			c.Request,
			usr.Id,
			audit.AdminLogout,
			audit.Fields{},
		)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}
	}

	c.Redirect(302, "/login")
}

func authRequestGet(c *gin.Context) {
	auth.Request(c)
}

func authCallbackGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sig := c.Query("sig")
	query := strings.Split(c.Request.URL.RawQuery, "&sig=")[0]

	usr, _, errData, err := auth.Callback(db, sig, query)
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

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.AdminPrimaryApprove,
		audit.Fields{
			"method": "callback",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	deviceAuth, secProviderId, errData, err := validator.ValidateAdmin(
		db, usr, false, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.AdminLoginFailed,
			audit.Fields{
				"method":  "callback",
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
		deviceCount, err := device.Count(db, usr.Id)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		secType := ""
		var secProvider bson.ObjectId
		if deviceCount == 0 {
			if secProviderId == "" {
				secType = secondary.AdminDeviceRegister
				secProvider = secondary.DeviceProvider
			} else {
				secType = secondary.Admin
				secProvider = secProviderId
			}
		} else {
			secType = secondary.AdminDevice
			secProvider = secondary.DeviceProvider
		}

		secd, err := secondary.New(db, usr.Id, secType, secProvider)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		urlQuery, err := secd.GetQuery()
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.Redirect(302, "/login?"+urlQuery)
		return
	} else if secProviderId != "" {
		secd, err := secondary.New(db, usr.Id, secondary.Admin, secProviderId)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		urlQuery, err := secd.GetQuery()
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.Redirect(302, "/login?"+urlQuery)
		return
	}

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.AdminLogin,
		audit.Fields{
			"method": "callback",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewAdmin(c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.Admin)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.Redirect(302, "/")
}

func authU2fRegisterGet(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	token := c.Query("token")

	secd, err := secondary.Get(db, token, secondary.AdminDeviceRegister)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			errData := &errortypes.ErrorData{
				Error:   "secondary_expired",
				Message: "Secondary authentication has expired",
			}
			c.JSON(401, errData)
		} else {
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	usr, err := secd.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.AdminDeviceRegisterRequest,
		audit.Fields{},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	resp, errData, err := secd.DeviceRegisterRequest(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.AdminLoginFailed,
			audit.Fields{
				"method":  "device_register",
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

	c.JSON(200, resp)
}

type u2fRegisterData struct {
	Token    string                   `json:"token"`
	Name     string                   `json:"name"`
	Response *u2flib.RegisterResponse `json:"response"`
}

func authU2fRegisterPost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)

	data := &u2fRegisterData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	secd, err := secondary.Get(db, data.Token, secondary.AdminDeviceRegister)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			errData := &errortypes.ErrorData{
				Error:   "secondary_expired",
				Message: "Secondary authentication has expired",
			}
			c.JSON(401, errData)
		} else {
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	usr, err := secd.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_, _, errData, err := validator.ValidateAdmin(
		db, usr, false, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.AdminLoginFailed,
			audit.Fields{
				"method":  "device_register",
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

	devc, errData, err := secd.DeviceRegisterResponse(
		db, data.Response, data.Name)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.DeviceRegisterFailed,
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
		audit.AdminDeviceRegister,
		audit.Fields{
			"device_id": devc.Id,
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "device.change")

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.AdminLogin,
		audit.Fields{
			"method": "device_register",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewAdmin(c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.Admin)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.Status(200)
}

func authU2fSignGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	token := c.Query("token")

	secd, err := secondary.Get(db, token, secondary.AdminDevice)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			errData := &errortypes.ErrorData{
				Error:   "secondary_expired",
				Message: "Secondary authentication has expired",
			}
			c.JSON(401, errData)
		} else {
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	usr, err := secd.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	resp, errData, err := secd.DeviceSignRequest(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.AdminLoginFailed,
			audit.Fields{
				"method":  "device",
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

	c.JSON(200, resp)
}

type u2fSignData struct {
	Token    string               `json:"token"`
	Response *u2flib.SignResponse `json:"response"`
}

func authU2fSignPost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &u2fSignData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	secd, err := secondary.Get(db, data.Token, secondary.AdminDevice)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			errData := &errortypes.ErrorData{
				Error:   "secondary_expired",
				Message: "Secondary authentication has expired",
			}
			c.JSON(401, errData)
		} else {
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	usr, err := secd.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_, secProviderId, errData, err := validator.ValidateAdmin(
		db, usr, false, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.AdminLoginFailed,
			audit.Fields{
				"method":  "device",
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

	errData, err = secd.DeviceSignResponse(db, data.Response)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.AdminLoginFailed,
			audit.Fields{
				"method":  "device",
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
		audit.AdminDeviceApprove,
		audit.Fields{},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if secProviderId != "" {
		secd, err := secondary.New(db, usr.Id, secondary.Admin, secProviderId)
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
		audit.AdminLogin,
		audit.Fields{
			"method": "device",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewAdmin(c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.Admin)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.Status(200)
}
