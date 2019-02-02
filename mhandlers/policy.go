package mhandlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/utils"
)

type policyData struct {
	Id                        primitive.ObjectID      `json:"id"`
	Name                      string                  `json:"name"`
	Services                  []primitive.ObjectID    `json:"services"`
	Authorities               []primitive.ObjectID    `json:"authorities"`
	Roles                     []string                `json:"roles"`
	Rules                     map[string]*policy.Rule `json:"rules"`
	AdminSecondary            primitive.ObjectID      `json:"admin_secondary"`
	UserSecondary             primitive.ObjectID      `json:"user_secondary"`
	ProxySecondary            primitive.ObjectID      `json:"proxy_secondary"`
	AuthoritySecondary        primitive.ObjectID      `json:"authority_secondary"`
	AdminDeviceSecondary      bool                    `json:"admin_device_secondary"`
	UserDeviceSecondary       bool                    `json:"user_device_secondary"`
	ProxyDeviceSecondary      bool                    `json:"proxy_device_secondary"`
	AuthorityDeviceSecondary  bool                    `json:"authority_device_secondary"`
	AuthorityRequireSmartCard bool                    `json:"authority_require_smart_card"`
}

func policyPut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &policyData{}

	polcyId, ok := utils.ParseObjectId(c.Param("policy_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	polcy, err := policy.Get(db, polcyId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	polcy.Name = data.Name
	polcy.Services = data.Services
	polcy.Authorities = data.Authorities
	polcy.Roles = data.Roles
	polcy.Rules = data.Rules
	polcy.AdminSecondary = data.AdminSecondary
	polcy.UserSecondary = data.UserSecondary
	polcy.ProxySecondary = data.ProxySecondary
	polcy.AuthoritySecondary = data.AuthoritySecondary
	polcy.AdminDeviceSecondary = data.AdminDeviceSecondary
	polcy.UserDeviceSecondary = data.UserDeviceSecondary
	polcy.ProxyDeviceSecondary = data.ProxyDeviceSecondary
	polcy.AuthorityDeviceSecondary = data.AuthorityDeviceSecondary
	polcy.AuthorityRequireSmartCard = data.AuthorityRequireSmartCard

	fields := set.NewSet(
		"name",
		"services",
		"authorities",
		"roles",
		"rules",
		"admin_secondary",
		"user_secondary",
		"proxy_secondary",
		"authority_secondary",
		"admin_device_secondary",
		"user_device_secondary",
		"proxy_device_secondary",
		"authority_device_secondary",
		"authority_require_smart_card",
	)

	errData, err := polcy.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = polcy.CommitFields(db, fields)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "policy.change")

	c.JSON(200, polcy)
}

func policyPost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &policyData{
		Name: "New Policy",
	}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	polcy := &policy.Policy{
		Name:                     data.Name,
		Services:                 data.Services,
		Authorities:              data.Authorities,
		Roles:                    data.Roles,
		Rules:                    data.Rules,
		AdminSecondary:           data.AdminSecondary,
		UserSecondary:            data.UserSecondary,
		ProxySecondary:           data.ProxySecondary,
		AuthoritySecondary:       data.AuthoritySecondary,
		AdminDeviceSecondary:     data.AdminDeviceSecondary,
		UserDeviceSecondary:      data.UserDeviceSecondary,
		ProxyDeviceSecondary:     data.ProxyDeviceSecondary,
		AuthorityDeviceSecondary: data.AuthorityDeviceSecondary,
	}

	errData, err := polcy.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = polcy.Insert(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "policy.change")

	c.JSON(200, polcy)
}

func policyDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)

	polcyId, ok := utils.ParseObjectId(c.Param("policy_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := policy.Remove(db, polcyId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "policy.change")

	c.JSON(200, nil)
}

func policyGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	polcyId, ok := utils.ParseObjectId(c.Param("policy_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	polcy, err := policy.Get(db, polcyId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.JSON(200, polcy)
}

func policiesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	policies, err := policy.GetAll(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.JSON(200, policies)
}
