package mhandlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/utils"
)

func sessionsGet(c *gin.Context) {
	if demo.IsDemo() {
		demo.Sessions[0].LastActive = time.Now()
		c.JSON(200, demo.Sessions)
		return
	}

	db := c.MustGet("db").(*database.Database)

	showRemoved, _ := strconv.ParseBool(c.Query("show_removed"))

	userId, ok := utils.ParseObjectId(c.Param("user_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	sessions, err := session.GetAll(db, userId, showRemoved)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.JSON(200, sessions)
}

func sessionDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

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

	_ = event.PublishDispatch(db, "session.change")

	c.JSON(200, nil)
}
