package mhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/csrf"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/session"
)

type csrfData struct {
	Token string `json:"token"`
}

func csrfGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	token, err := csrf.NewToken(db, sess.Id)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	data := &csrfData{
		Token: token,
	}
	c.JSON(200, data)
}
