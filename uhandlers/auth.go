package uhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/auth"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/cookie"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/session"
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

	errData, err = validator.ValidateUser(db, usr, false, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.LoginFailed,
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
		audit.Login,
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

func logoutAllGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)

	usr, err := authr.GetUser(db)
	if err != nil {
		return
	}

	sessions, err := session.GetAll(db, usr.Id, false)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	for _, sess := range sessions {
		println(sess.Id)
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

	c.Redirect(302, "/")
}

func authRequestGet(c *gin.Context) {
	auth.Request(c)
}

func authCallbackGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sig := c.Query("sig")
	query := strings.Split(c.Request.URL.RawQuery, "&sig=")[0]

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

	errData, err = validator.ValidateUser(db, usr, false, c.Request)
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
		audit.Login,
		audit.Fields{
			"method": "sso",
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
