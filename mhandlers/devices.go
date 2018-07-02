package mhandlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/device"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/secondary"
	"github.com/pritunl/pritunl-zero/u2flib"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
)

type deviceData struct {
	User         bson.ObjectId `json:"user"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	SshPublicKey string        `json:"ssh_public_key"`
}

func devicePut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
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

	devc, err := device.Get(db, devcId)
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

func devicePost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &deviceData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	devc := device.New(data.User, data.Type, device.Ssh)

	devc.Name = data.Name
	devc.SshPublicKey = data.SshPublicKey

	errData, err := devc.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = devc.Insert(db)
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

	devcId, ok := utils.ParseObjectId(c.Param("device_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := device.Remove(db, devcId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "device.change")

	c.JSON(200, nil)
}

func devicesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	usrId, ok := utils.ParseObjectId(c.Param("user_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	devices, err := device.GetAll(db, usrId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.JSON(200, devices)
}

type devicesU2fRegisterRespData struct {
	Token   string      `json:"token"`
	Request interface{} `json:"request"`
}

func deviceU2fRegisterGet(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)

	usrId, ok := utils.ParseObjectId(c.Param("user_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	secd, err := secondary.New(db, usrId,
		secondary.AdminDeviceRegister, secondary.DeviceProvider)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	jsonResp, errData, err := secd.DeviceRegisterRequest(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	resp := &devicesU2fRegisterRespData{
		Token:   secd.Id,
		Request: jsonResp,
	}

	c.JSON(200, resp)
}

type devicesU2fRegisterData struct {
	Token    string                   `json:"token"`
	Name     string                   `json:"name"`
	Response *u2flib.RegisterResponse `json:"response"`
}

func deviceU2fRegisterPost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)
	data := &devicesU2fRegisterData{}

	usrId, ok := utils.ParseObjectId(c.Param("user_id"))
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

	secd, err := secondary.Get(db, data.Token,
		secondary.AdminDeviceRegister)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			errData := &errortypes.ErrorData{
				Error:   "secondary_expired",
				Message: "Secondary authentication has expired",
			}
			c.JSON(400, errData)
		} else {
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	devc, errData, err := secd.DeviceRegisterResponse(
		db, data.Response, data.Name)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = audit.New(
		db,
		c.Request,
		usrId,
		audit.AdminDeviceRegister,
		audit.Fields{
			"admin_id":  usr.Id,
			"device_id": devc.Id,
		},
	)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "device.change")

	c.JSON(200, nil)
}
