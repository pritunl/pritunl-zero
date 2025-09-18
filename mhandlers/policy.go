package mhandlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/utils"
)

type policyData struct {
	Id                        primitive.ObjectID      `json:"id"`
	Name                      string                  `json:"name"`
	Disabled                  bool                    `json:"disabled"`
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

type policiesData struct {
	Policies []*policy.Policy `json:"policies"`
	Count    int64            `json:"count"`
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
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	polcy, err := policy.Get(db, polcyId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	polcy.Name = data.Name
	polcy.Disabled = data.Disabled
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
		"disabled",
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

	_ = event.PublishDispatch(db, "policy.change")

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
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	polcy := &policy.Policy{
		Name:                     data.Name,
		Disabled:                 data.Disabled,
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

	_ = event.PublishDispatch(db, "policy.change")

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

	_ = event.PublishDispatch(db, "policy.change")

	c.JSON(200, nil)
}

func policiesDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := []primitive.ObjectID{}

	err := c.Bind(&data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	err = policy.RemoveMulti(db, data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "policy.change")

	c.JSON(200, nil)
}

func policyGet(c *gin.Context) {
	if demo.IsDemo() {
		polcy := demo.Policies[0]
		c.JSON(200, polcy)
		return
	}

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
	if demo.IsDemo() {
		data := &policiesData{
			Policies: demo.Policies,
			Count:    int64(len(demo.Policies)),
		}

		c.JSON(200, data)
		return
	}

	db := c.MustGet("db").(*database.Database)

	page, _ := strconv.ParseInt(c.Query("page"), 10, 0)
	pageCount, _ := strconv.ParseInt(c.Query("page_count"), 10, 0)

	query := bson.M{}

	policyId, ok := utils.ParseObjectId(c.Query("id"))
	if ok {
		query["_id"] = policyId
	}

	name := strings.TrimSpace(c.Query("name"))
	if name != "" {
		query["name"] = &bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", regexp.QuoteMeta(name)),
			"$options": "i",
		}
	}

	organization, ok := utils.ParseObjectId(c.Query("organization"))
	if ok {
		query["organization"] = organization
	}

	description := strings.TrimSpace(c.Query("description"))
	if description != "" {
		query["description"] = &bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", regexp.QuoteMeta(description)),
			"$options": "i",
		}
	}

	policies, count, err := policy.GetAllPaged(db, &query, page, pageCount)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	data := &policiesData{
		Policies: policies,
		Count:    count,
	}

	c.JSON(200, data)
}
