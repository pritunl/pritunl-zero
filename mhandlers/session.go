package mhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/utils"
)

func sessionsGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	userId, ok := utils.ParseObjectId(c.Param("user_id"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	sessions, err := session.GetAll(db, userId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, sessions)
}

func sessionDelete(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	sessionId := c.Param("session_id")
	if sessionId == "" {
		c.AbortWithStatus(400)
		return
	}

	err := session.Remove(db, sessionId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, nil)
}
