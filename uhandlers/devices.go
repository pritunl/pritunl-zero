package uhandlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/device"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/utils"
)

type deviceData struct {
	Name string `json:"name"`
}

func devicePut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)
	data := &deviceData{}

	devcId, ok := utils.ParseObjectId(c.Param("device_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	usr, err := authr.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	devc, err := device.GetUser(db, devcId, usr.Id)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	devc.Name = data.Name

	fields := set.NewSet(
		"name",
	)

	errData, err := devc.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = devc.CommitFields(db, fields)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "device.change")

	c.JSON(200, devc)
}

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
