package uhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/device"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/utils"
)

func deviceDelete(c *gin.Context) {
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

	devcId, ok := utils.ParseObjectId(c.Param("device_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err = device.RemoveUser(db, devcId, usr.Id)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "device.change")

	c.JSON(200, nil)
}

func devicesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)

	usr, err := authr.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	devices, err := device.GetAll(db, usr.Id)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.JSON(200, devices)
}
