package mhandlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/authority"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type authorityData struct {
	Id                 bson.ObjectId `json:"id"`
	Name               string        `json:"name"`
	Type               string        `json:"type"`
	Expire             int           `json:"expire"`
	HostExpire         int           `json:"host_expire"`
	MatchRoles         bool          `json:"match_roles"`
	Roles              []string      `json:"roles"`
	HostDomain         string        `json:"host_domain"`
	HostProxy          string        `json:"host_proxy"`
	HostCertificates   bool          `json:"host_certificates"`
	StrictHostChecking bool          `json:"strict_host_checking"`
}

func authorityPut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &authorityData{}

	authrId, ok := utils.ParseObjectId(c.Param("authr_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	authr, err := authority.Get(db, authrId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	authr.Name = data.Name
	authr.Type = data.Type
	authr.Expire = data.Expire
	authr.HostExpire = data.HostExpire
	authr.MatchRoles = data.MatchRoles
	authr.Roles = data.Roles

	if !authr.HostCertificates && data.HostCertificates {
		err = authr.TokenNew()
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}
	}
	authr.HostDomain = data.HostDomain
	authr.HostProxy = data.HostProxy
	authr.HostCertificates = data.HostCertificates
	authr.StrictHostChecking = data.StrictHostChecking

	fields := set.NewSet(
		"name",
		"type",
		"expire",
		"host_expire",
		"info",
		"match_roles",
		"roles",
		"host_domain",
		"host_tokens",
		"host_proxy",
		"host_certificates",
		"strict_host_checking",
	)

	errData, err := authr.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = authr.CommitFields(db, fields)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "authority.change")

	c.JSON(200, authr)
}

func authorityPost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &authorityData{
		Name: "New Authority",
	}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	authr := &authority.Authority{
		Name:               data.Name,
		Type:               data.Type,
		Expire:             data.Expire,
		HostExpire:         data.HostExpire,
		MatchRoles:         data.MatchRoles,
		Roles:              data.Roles,
		HostDomain:         data.HostDomain,
		StrictHostChecking: data.StrictHostChecking,
	}

	err = authr.GenerateRsaPrivateKey()
	if err != nil {
		return
	}

	errData, err := authr.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = authr.Insert(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "authority.change")

	c.JSON(200, authr)
}

func authorityDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)

	authrId, ok := utils.ParseObjectId(c.Param("authr_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := authority.Remove(db, authrId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "authority.change")

	c.JSON(200, nil)
}

func authorityGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	authrId, ok := utils.ParseObjectId(c.Param("authr_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	authr, err := authority.Get(db, authrId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if demo.IsDemo() {
		for i := range authr.HostTokens {
			authr.HostTokens[i] = "demo"
		}
	}

	c.JSON(200, authr)
}

func authoritysGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	authrs, err := authority.GetAll(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if demo.IsDemo() {
		for _, authr := range authrs {
			for i := range authr.HostTokens {
				authr.HostTokens[i] = "demo"
			}
		}
	}

	c.JSON(200, authrs)
}

func authorityPublicKeyGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	authrIdsStr := strings.Split(c.Param("authr_ids"), ",")
	authrIds := []bson.ObjectId{}

	for _, authrIdStr := range authrIdsStr {
		if authrIdStr == "" {
			continue
		}

		authrId, ok := utils.ParseObjectId(authrIdStr)
		if !ok {
			utils.AbortWithStatus(c, 400)
			return
		}

		authrIds = append(authrIds, authrId)
	}

	if len(authrIds) == 0 {
		utils.AbortWithStatus(c, 400)
		return
	}

	publicKeys := ""

	authrs, err := authority.GetMulti(db, authrIds)
	if err != nil {
		return
	}

	for _, authr := range authrs {
		publicKeys += strings.TrimSpace(authr.PublicKey) + "\n"
	}

	c.String(200, publicKeys)
}

func authorityTokenPost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)

	authrId, ok := utils.ParseObjectId(c.Param("authr_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	authr, err := authority.Get(db, authrId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	err = authr.TokenNew()
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	err = authr.CommitFields(db, set.NewSet("host_tokens"))
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "authority.change")

	c.Status(200)
}

func authorityTokenDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)

	authrId, ok := utils.ParseObjectId(c.Param("authr_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	token := c.Param("token")
	if token == "" {
		utils.AbortWithStatus(c, 400)
		return
	}

	authr, err := authority.Get(db, authrId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	err = authr.TokenDelete(token)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	err = authr.CommitFields(db, set.NewSet("host_tokens"))
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "authority.change")

	c.Status(200)
}
