package uhandlers

import (
	"strings"

	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
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
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/u2flib"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/pritunl/pritunl-zero/validator"
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
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
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
		audit.UserPrimaryApprove,
		audit.Fields{
			"method": "local",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	devAuth, secProviderId, errAudit, errData, err := validator.ValidateUser(
		db, usr, false, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		if errAudit == nil {
			errAudit = audit.Fields{
				"error":   errData.Error,
				"message": errData.Message,
			}
		}
		errAudit["method"] = "local"

		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.UserLoginFailed,
			errAudit,
		)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(401, errData)
		return
	}

	if devAuth {
		deviceCount, err := device.CountSecondary(db, usr.Id)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		secType := ""
		var secProvider primitive.ObjectID
		if deviceCount == 0 {
			if secProviderId.IsZero() {
				secType = secondary.UserDeviceRegister
				secProvider = secondary.DeviceProvider
			} else {
				secType = secondary.User
				secProvider = secProviderId
			}
		} else {
			secType = secondary.UserDevice
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
	} else if !secProviderId.IsZero() {
		secd, err := secondary.New(db, usr.Id, secondary.User, secProviderId)
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
		audit.UserLogin,
		audit.Fields{
			"method": "local",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewUser(c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.User)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	redirectQueryJson(c, c.Request.URL.RawQuery)
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
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	secd, err := secondary.Get(db, data.Token, secondary.User)
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
			c.Status(206)
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
			audit.UserLoginFailed,
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
		audit.UserSecondaryApprove,
		audit.Fields{
			"provider_id": secd.ProviderId,
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	deviceAuth, _, errAudit, errData, err := validator.ValidateUser(
		db, usr, false, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		if errAudit == nil {
			errAudit = audit.Fields{
				"error":   errData.Error,
				"message": errData.Message,
			}
		}
		errAudit["method"] = "secondary"
		errAudit["provider_id"] = secd.ProviderId

		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.UserLoginFailed,
			errAudit,
		)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(401, errData)
		return
	}

	if deviceAuth {
		deviceCount, err := device.CountSecondary(db, usr.Id)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		if deviceCount == 0 {
			secd, err := secondary.New(db, usr.Id,
				secondary.UserDeviceRegister, secondary.DeviceProvider)
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
		audit.UserLogin,
		audit.Fields{
			"method":      "secondary",
			"provider_id": secd.ProviderId,
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewUser(c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.User)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	redirectQueryJson(c, c.Request.URL.RawQuery)
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
			audit.UserLogout,
			audit.Fields{},
		)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}
	}

	c.Redirect(302, "/")
}

func logoutAllGet(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)

	usr, err := authr.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	sessions, err := session.GetAll(db, usr.Id, false)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	for _, sess := range sessions {
		err = sess.Remove(db)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}
	}

	if authr.IsValid() {
		err := authr.Remove(db)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}
	}

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.UserLogoutAll,
		audit.Fields{},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.Redirect(302, "/")
}

func authRequestGet(c *gin.Context) {
	auth.Request(c)
}

func authCallbackGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sig := c.Query("sig")
	query := strings.Split(c.Request.URL.RawQuery, "&sig=")[0]

	usr, tokn, errAudit, errData, err := auth.Callback(db, sig, query)
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
		if usr != nil {
			if errAudit == nil {
				errAudit = audit.Fields{
					"error":   errData.Error,
					"message": errData.Message,
				}
			}
			errAudit["method"] = "callback"

			err = audit.New(
				db,
				c.Request,
				usr.Id,
				audit.UserLoginFailed,
				errAudit,
			)
			if err != nil {
				utils.AbortWithError(c, 500, err)
				return
			}
		}

		c.JSON(401, errData)
		return
	}

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.UserPrimaryApprove,
		audit.Fields{
			"method": "callback",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	devAuth, secProviderId, errAudit, errData, err := validator.ValidateUser(
		db, usr, false, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		if errAudit == nil {
			errAudit = audit.Fields{
				"error":   errData.Error,
				"message": errData.Message,
			}
		}
		errAudit["method"] = "callback"

		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.UserLoginFailed,
			errAudit,
		)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(401, errData)
		return
	}

	if devAuth {
		deviceCount, err := device.CountSecondary(db, usr.Id)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		secType := ""
		var secProvider primitive.ObjectID
		if deviceCount == 0 {
			if secProviderId.IsZero() {
				secType = secondary.UserDeviceRegister
				secProvider = secondary.DeviceProvider
			} else {
				secType = secondary.User
				secProvider = secProviderId
			}
		} else {
			secType = secondary.UserDevice
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

		if tokn.Query != "" {
			urlQuery += "&" + tokn.Query
		}

		c.Redirect(302, "/login?"+urlQuery)
		return
	} else if !secProviderId.IsZero() {
		secd, err := secondary.New(db, usr.Id, secondary.User, secProviderId)
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
		audit.UserLogin,
		audit.Fields{
			"method": "callback",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewUser(c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.User)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	redirectQuery(c, tokn.Query)
}

func authU2fAppGet(c *gin.Context) {
	c.JSON(200, device.GetFacets())
}

func authU2fRegisterGet(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	token := c.Query("token")

	if settings.Local.AppId == "" {
		errData := &errortypes.ErrorData{
			Error: "user_node_unavailable",
			Message: "At least one node must have a user domain configured " +
				"to use secondary device authentication",
		}
		c.JSON(400, errData)
		return
	}

	secd, err := secondary.Get(db, token, secondary.UserDeviceRegister)
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
			audit.UserLoginFailed,
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
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	secd, err := secondary.Get(db, data.Token, secondary.UserDeviceRegister)
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
		audit.UserDeviceRegisterRequest,
		audit.Fields{},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_, _, errAudit, errData, err := validator.ValidateUser(
		db, usr, false, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		if errAudit == nil {
			errAudit = audit.Fields{
				"error":   errData.Error,
				"message": errData.Message,
			}
		}
		errAudit["method"] = "device_register"

		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.UserLoginFailed,
			errAudit,
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
		audit.UserDeviceRegister,
		audit.Fields{
			"device_id": devc.Id,
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "device.change")

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.UserLogin,
		audit.Fields{
			"method": "device_register",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewUser(c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.User)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	redirectQueryJson(c, c.Request.URL.RawQuery)
}

func authU2fSignGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	token := c.Query("token")

	secd, err := secondary.Get(db, token, secondary.UserDevice)
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
			audit.UserLoginFailed,
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
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	secd, err := secondary.Get(db, data.Token, secondary.UserDevice)
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

	_, secProviderId, errAudit, errData, err := validator.ValidateUser(
		db, usr, false, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		if errAudit == nil {
			errAudit = audit.Fields{
				"error":   errData.Error,
				"message": errData.Message,
			}
		}
		errAudit["method"] = "device"

		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.UserLoginFailed,
			errAudit,
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
			audit.UserLoginFailed,
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
		audit.UserDeviceApprove,
		audit.Fields{},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if !secProviderId.IsZero() {
		secd, err := secondary.New(db, usr.Id, secondary.User, secProviderId)
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
		audit.UserLogin,
		audit.Fields{
			"method": "device",
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cook := cookie.NewUser(c.Writer, c.Request)

	_, err = cook.NewSession(db, c.Request, usr.Id, true, session.User)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	redirectQueryJson(c, c.Request.URL.RawQuery)
}
