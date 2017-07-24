package phandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/auth"
	"github.com/pritunl/pritunl-zero/cookie"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/session"
	"strings"
)

func authStateGet(c *gin.Context) {
	data := auth.GetState()
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
		c.AbortWithError(500, err)
		return
	}

	usr, errData, err := auth.Local(db, data.Username, data.Password)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if errData != nil {
		c.JSON(401, errData)
		return
	}

	errData, err = auth.Validate(db, usr, srvc)
	if err != nil {
		return
	}

	if errData != nil {
		c.JSON(401, errData)
		return
	}

	cook := cookie.NewProxy(srvc, c.Writer, c.Request)

	_, err = cook.NewSession(db, usr.Id, true)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Status(200)
}

func logoutGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	if sess != nil {
		err := sess.Remove(db)
		if err != nil {
			c.AbortWithError(500, err)
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

	usr, errData, err := auth.Local(db, sig, query)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if errData != nil {
		c.JSON(401, errData)
		return
	}

	errData, err = auth.Validate(db, usr, srvc)
	if err != nil {
		return
	}

	if errData != nil {
		c.JSON(401, errData)
		return
	}

	cook := cookie.New(c.Writer, c.Request)

	_, err = cook.NewSession(db, usr.Id, true)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Redirect(302, "/")
}
