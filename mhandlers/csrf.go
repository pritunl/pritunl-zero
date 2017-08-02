package mhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/csrf"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/utils"
)

type csrfData struct {
	Token string `json:"token"`
	Theme string `json:"theme"`
}

func csrfGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	usr, err := sess.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	token, err := csrf.NewToken(db, sess.Id)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	data := &csrfData{
		Token: token,
		Theme: usr.Theme,
	}
	c.JSON(200, data)
}
