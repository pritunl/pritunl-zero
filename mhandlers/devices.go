package mhandlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/alertevent"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/device"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/secondary"
	"github.com/pritunl/pritunl-zero/utils"
)

type deviceData struct {
	User         bson.ObjectID `json:"user"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Mode         string        `json:"mode"`
	Number       string        `json:"number"`
	AlertLevels  []int         `json:"alert_levels"`
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
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	devc, err := device.Get(db, devcId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	devc.Name = data.Name
	devc.AlertLevels = data.AlertLevels

	fields := set.NewSet(
		"name",
		"alert_levels",
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

	_ = event.PublishDispatch(db, "device.change")

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
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	devc := device.New(data.User, data.Type, data.Mode)

	devc.Name = data.Name
	devc.Number = data.Number
	devc.AlertLevels = data.AlertLevels
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

	_ = event.PublishDispatch(db, "device.change")

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

	_ = event.PublishDispatch(db, "device.change")

	c.JSON(200, nil)
}

func devicesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	usrId, ok := utils.ParseObjectId(c.Param("user_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	devices, err := device.GetAllSorted(db, usrId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.JSON(200, devices)
}

func deviceAlertPost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &deviceData{}

	devcId, ok := utils.ParseObjectId(c.Param("resource_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := c.Bind(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	devc, err := device.Get(db, devcId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	errData, err := alertevent.SendTest(db, devc)
	if errData != nil {
		c.JSON(400, errData)
		return
	}

	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "device.change")

	c.JSON(200, devc)
}

func deviceMethodPost(c *gin.Context) {
	switch c.Param("method") {
	case "alert":
		deviceAlertPost(c)
		return
	default:
		utils.AbortWithStatus(c, 404)
		return
	}

	return
}

type devicesWanRegisterRespData struct {
	Token   string      `json:"token"`
	Options interface{} `json:"options"`
}

func deviceWanRegisterGet(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)

	if node.Self.WebauthnDomain == "" {
		errData := &errortypes.ErrorData{
			Error:   "webauthn_domain_unavailable",
			Message: "WebAuthn domain must be configured",
		}
		c.JSON(400, errData)
		return
	}

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

	jsonResp, errData, err := secd.DeviceRegisterRequest(db,
		utils.GetOrigin(c.Request))
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	resp := &devicesWanRegisterRespData{
		Token:   secd.Id,
		Options: jsonResp,
	}

	c.JSON(200, resp)
}

type devicesWanRegisterData struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}

func deviceWanRegisterPost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)
	data := &devicesWanRegisterData{}

	body, err := utils.CopyBody(c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	err = c.Bind(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	usrId, ok := utils.ParseObjectId(c.Param("resource_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
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
		db, utils.GetOrigin(c.Request), body, data.Name)
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

	_ = event.PublishDispatch(db, "device.change")

	c.JSON(200, nil)
}
