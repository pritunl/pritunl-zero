package mhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/dropbox/godropbox/container/set"
)

type themeData struct {
	Theme string `json:"theme"`
}

func themePut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)
	data := &themeData{}

	err := c.Bind(&data)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	usr, err := sess.GetUser(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	usr.Theme = data.Theme

	err = usr.CommitFields(db, set.NewSet("theme"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	return
}
