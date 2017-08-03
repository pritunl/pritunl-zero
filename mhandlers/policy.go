package mhandlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
)

type policyData struct {
	Id       bson.ObjectId           `json:"id"`
	Name     string                  `json:"name"`
	Services []bson.ObjectId         `json:"services"`
	Roles    []string                `json:"roles"`
	Rules    map[string]*policy.Rule `json:"rules"`
}

func policyPut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &policyData{}

	certId, ok := utils.ParseObjectId(c.Param("policy_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cert, err := policy.Get(db, certId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cert.Name = data.Name
	cert.Services = data.Services
	cert.Roles = data.Roles
	cert.Rules = data.Rules

	fields := set.NewSet(
		"name",
		"services",
		"roles",
		"rules",
	)

	errData, err := cert.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = cert.CommitFields(db, fields)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "policy.change")

	c.JSON(200, cert)
}

func policyPost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &policyData{
		Name: "New Policy",
	}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cert := &policy.Policy{
		Name:     data.Name,
		Services: data.Services,
		Roles:    data.Roles,
		Rules:    data.Rules,
	}

	errData, err := cert.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = cert.Insert(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "policy.change")

	c.JSON(200, cert)
}

func policyDelete(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	certId, ok := utils.ParseObjectId(c.Param("policy_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := policy.Remove(db, certId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "policy.change")

	c.JSON(200, nil)
}

func policyGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	certId, ok := utils.ParseObjectId(c.Param("policy_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	cert, err := policy.Get(db, certId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.JSON(200, cert)
}

func policiesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	certs, err := policy.GetAll(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.JSON(200, certs)
}
