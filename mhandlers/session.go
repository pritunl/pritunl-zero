package mhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/utils"
)

func sessionsGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	userId, ok := utils.ParseObjectId(c.Param("user_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	sessions, err := session.GetAll(db, userId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.JSON(200, sessions)
}

func sessionDelete(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	sessionId := c.Param("session_id")
	if sessionId == "" {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := session.Remove(db, sessionId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "session.change")

	c.JSON(200, nil)
}
