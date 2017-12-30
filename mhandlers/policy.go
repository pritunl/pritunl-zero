package mhandlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
)

type policyData struct {
	Id                 bson.ObjectId           `json:"id"`
	Name               string                  `json:"name"`
	Services           []bson.ObjectId         `json:"services"`
	Authorities        []bson.ObjectId         `json:"authorities"`
	Roles              []string                `json:"roles"`
	Rules              map[string]*policy.Rule `json:"rules"`
	KeybaseMode        string                  `json:"keybase_mode"`
	AdminSecondary     bson.ObjectId           `json:"admin_secondary"`
	UserSecondary      bson.ObjectId           `json:"user_secondary"`
	ProxySecondary     bson.ObjectId           `json:"proxy_secondary"`
	AuthoritySecondary bson.ObjectId           `json:"authority_secondary"`
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
	polcy.KeybaseMode = data.KeybaseMode
	polcy.AdminSecondary = data.AdminSecondary
	polcy.UserSecondary = data.UserSecondary
	polcy.ProxySecondary = data.ProxySecondary
	polcy.AuthoritySecondary = data.AuthoritySecondary

	fields := set.NewSet(
		"name",
		"services",
		"authorities",
		"roles",
		"rules",
		"keybase_mode",
		"admin_secondary",
		"user_secondary",
		"proxy_secondary",
		"authority_secondary",
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
		Name:               data.Name,
		Services:           data.Services,
		Authorities:        data.Authorities,
		Roles:              data.Roles,
		Rules:              data.Rules,
		KeybaseMode:        data.KeybaseMode,
		AdminSecondary:     data.AdminSecondary,
		UserSecondary:      data.UserSecondary,
		ProxySecondary:     data.ProxySecondary,
		AuthoritySecondary: data.AuthoritySecondary,
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
