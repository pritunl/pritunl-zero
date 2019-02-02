package uhandlers

import (
	"encoding/base64"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/device"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/secondary"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/u2flib"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/pritunl/pritunl-zero/validator"
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

	count, err := device.Count(db, usr.Id)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if count <= 1 {
		usr.Disabled = true
		err = usr.CommitFields(db, set.NewSet("disabled"))
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}
	}

	err = device.RemoveUser(db, devcId, usr.Id)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	count, err = device.Count(db, usr.Id)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if count == 0 {
		if !usr.Disabled {
			usr.Disabled = true
			err = usr.CommitFields(db, set.NewSet("disabled"))
			if err != nil {
				utils.AbortWithError(c, 500, err)
				return
			}
		}

		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.UserAccountDisable,
			audit.Fields{
				"reason": "All authentication devices removed",
			},
		)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		errData := &errortypes.ErrorData{
			Error:   "device_empty",
			Message: "Account disabled contact an administrator",
		}
		c.JSON(401, errData)
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

	devices, err := device.GetAllSorted(db, usr.Id)
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
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)
	deviceType := c.Query("type")

	if settings.Local.AppId == "" {
		errData := &errortypes.ErrorData{
			Error: "user_node_unavailable",
			Message: "At least one node must have a user domain configured " +
				"to use secondary device authentication",
		}
		c.JSON(400, errData)
		return
	}

	usr, err := authr.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_, secProviderId, errAudit, errData, err := validator.ValidateUser(
		db, usr, false, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		if errAudit == nil {
			errAudit = audit.Fields{
				"error":   errData.Error,
				"message": errData.Message,
			}
		}
		errAudit["method"] = "add_device_register"

		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.UserAuthFailed,
			errAudit,
		)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(400, errData)
		return
	}

	deviceCount, err := device.Count(db, usr.Id)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if deviceCount > 0 || !secProviderId.IsZero() {
		secType := ""
		var secProvider primitive.ObjectID

		if deviceCount == 0 {
			if deviceType == device.SmartCard {
				errData := &errortypes.ErrorData{
					Error:   "no_devices",
					Message: "Cannot register Smart Card without a U2F device",
				}
				c.JSON(401, errData)
				return
			}

			secType = secondary.UserManage
			secProvider = secProviderId
		} else {
			secType = secondary.UserManageDevice
			secProvider = secondary.DeviceProvider
		}

		secd, err := secondary.New(db, usr.Id, secType, secProvider)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		data, err := secd.GetData()
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(201, data)
		return
	}

	if deviceType == device.SmartCard {
		errData := &errortypes.ErrorData{
			Error:   "no_devices",
			Message: "Cannot register Smart Card without a U2F device",
		}
		c.JSON(401, errData)
		return
	}

	secd, err := secondary.New(db, usr.Id, secondary.UserManageDeviceRegister,
		secondary.DeviceProvider)
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
	Type         string                   `json:"type"`
	Token        string                   `json:"token"`
	Name         string                   `json:"name"`
	Response     *u2flib.RegisterResponse `json:"response"`
	SshPublicKey string                   `json:"ssh_public_key"`
}

func deviceU2fRegisterPost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)

	data := &devicesU2fRegisterData{}

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
		secondary.UserManageDeviceRegister)
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

	var devc *device.Device
	if data.Type == device.SmartCard {
		deviceCount, err := device.Count(db, usr.Id)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}
		if deviceCount == 0 {
			errData := &errortypes.ErrorData{
				Error:   "no_devices",
				Message: "Cannot register Smart Card without a U2F device",
			}
			c.JSON(401, errData)
			return
		}

		sshPubKey, err := base64.URLEncoding.DecodeString(data.SshPublicKey)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err,
					"uhandlers: Failed to decode SSH public key"),
			}
			utils.AbortWithError(c, 500, err)
			return
		}

		dvc, errData, err := secd.DeviceRegisterSmartCard(
			db, string(sshPubKey), data.Name)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		if errData != nil {
			c.JSON(400, errData)
			return
		}

		devc = dvc
	} else {
		dvc, errData, err := secd.DeviceRegisterResponse(
			db, data.Response, data.Name)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		if errData != nil {
			c.JSON(400, errData)
			return
		}

		devc = dvc
	}

	err = audit.New(
		db,
		c.Request,
		usr.Id,
		audit.DeviceRegister,
		audit.Fields{
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

type deviceU2fSecondaryData struct {
	Type     string `json:"type"`
	Token    string `json:"token"`
	Factor   string `json:"factor"`
	Passcode string `json:"passcode"`
}

func deviceU2fSecondaryPut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)
	data := &deviceU2fSecondaryData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	secd, err := secondary.Get(db, data.Token, secondary.UserManage)
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

	errData, err := secd.Handle(db, c.Request, data.Factor, data.Passcode)
	if err != nil {
		if _, ok := err.(*secondary.IncompleteError); ok {
			c.Status(206)
		} else {
			utils.AbortWithError(c, 500, err)
		}
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	usr, err := authr.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	secd, err = secondary.New(db, usr.Id, secondary.UserManageDeviceRegister,
		secondary.DeviceProvider)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	var jsonResp interface{}
	if data.Type != device.SmartCard {
		jsonResp, errData, err = secd.DeviceRegisterRequest(db)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		if errData != nil {
			c.JSON(400, errData)
			return
		}
	}

	resp := &devicesU2fRegisterRespData{
		Token:   secd.Id,
		Request: jsonResp,
	}

	c.JSON(200, resp)
}

func deviceU2fSignGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	token := c.Query("token")

	secd, err := secondary.Get(db, token, secondary.UserManageDevice)
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

	resp, errData, err := secd.DeviceSignRequest(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	c.JSON(200, resp)
}

type deviceU2fSignData struct {
	Type     string               `json:"type"`
	Token    string               `json:"token"`
	Response *u2flib.SignResponse `json:"response"`
}

func deviceU2fSignPost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authorizer").(*authorizer.Authorizer)
	data := &deviceU2fSignData{}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	secd, err := secondary.Get(db, data.Token, secondary.UserManageDevice)
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

	errData, err := secd.DeviceSignResponse(db, data.Response)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	usr, err := authr.GetUser(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_, secProviderId, errAudit, errData, err := validator.ValidateUser(
		db, usr, false, c.Request)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		if errAudit == nil {
			errAudit = audit.Fields{
				"error":   errData.Error,
				"message": errData.Message,
			}
		}
		errAudit["method"] = "add_device_register"

		err = audit.New(
			db,
			c.Request,
			usr.Id,
			audit.UserAuthFailed,
			errAudit,
		)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(400, errData)
		return
	}

	if !secProviderId.IsZero() {
		secd, err := secondary.New(db, usr.Id,
			secondary.UserManage, secProviderId)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		data, err := secd.GetData()
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(201, data)
		return
	}

	secd, err = secondary.New(db, usr.Id, secondary.UserManageDeviceRegister,
		secondary.DeviceProvider)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	var jsonResp interface{}
	if data.Type != device.SmartCard {
		jsonResp, errData, err = secd.DeviceRegisterRequest(db)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		if errData != nil {
			c.JSON(400, errData)
			return
		}
	}

	resp := &devicesU2fRegisterRespData{
		Token:   secd.Id,
		Request: jsonResp,
	}

	c.JSON(200, resp)
}
