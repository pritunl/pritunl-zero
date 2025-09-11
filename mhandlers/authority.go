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
	"github.com/pritunl/pritunl-zero/authority"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/utils"
)

type authorityData struct {
	Id                 primitive.ObjectID `json:"id"`
	Name               string             `json:"name"`
	Type               string             `json:"type"`
	Algorithm          string             `json:"algorithm"`
	KeyIdFormat        string             `json:"key_id_format"`
	Expire             int                `json:"expire"`
	HostExpire         int                `json:"host_expire"`
	MatchRoles         bool               `json:"match_roles"`
	Roles              []string           `json:"roles"`
	ProxyHosting       bool               `json:"proxy_hosting"`
	ProxyHostname      string             `json:"proxy_hostname"`
	ProxyPort          int                `json:"proxy_port"`
	HostDomain         string             `json:"host_domain"`
	HostMatches        []string           `json:"host_matches"`
	HostSubnets        []string           `json:"host_subnets"`
	HostProxy          string             `json:"host_proxy"`
	HostCertificates   bool               `json:"host_certificates"`
	StrictHostChecking bool               `json:"strict_host_checking"`
	HsmToken           string             `json:"hsm_token"`
	HsmSecret          string             `json:"hsm_secret"`
	HsmSerial          string             `json:"hsm_serial"`
	HsmGenerateSecret  bool               `json:"hsm_generate_secret"`
	ResetProxyHostKey  bool               `json:"reset_proxy_host_key"`
}

type authoritiesData struct {
	Authorities []*authority.Authority `json:"authorities"`
	Count       int64                  `json:"count"`
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
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	authr, err := authority.Get(db, authrId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	showSecret := false
	if authr.Type != data.Type {
		if data.Type == authority.PritunlHsm {
			err = authr.GenerateHsmToken()
			if err != nil {
				utils.AbortWithError(c, 500, err)
				return
			}
			showSecret = true
		}
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
	authr.KeyIdFormat = data.KeyIdFormat
	authr.ProxyHosting = data.ProxyHosting
	authr.ProxyHostname = data.ProxyHostname
	authr.ProxyPort = data.ProxyPort
	authr.HostMatches = data.HostMatches
	authr.HostSubnets = data.HostSubnets
	authr.HostDomain = data.HostDomain
	authr.HostProxy = data.HostProxy
	authr.HostCertificates = data.HostCertificates
	authr.StrictHostChecking = data.StrictHostChecking
	authr.HsmSerial = data.HsmSerial

	if authr.Type == authority.PritunlHsm && data.HsmGenerateSecret {
		err = authr.GenerateHsmToken()
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}
		showSecret = true
	}

	fields := set.NewSet(
		"name",
		"type",
		"expire",
		"host_expire",
		"public_key",
		"public_key_pem",
		"root_certificate",
		"private_key",
		"info",
		"match_roles",
		"roles",
		"key_id_format",
		"proxy_hosting",
		"proxy_hostname",
		"proxy_port",
		"host_domain",
		"host_matches",
		"host_subnets",
		"host_tokens",
		"host_proxy",
		"host_certificates",
		"strict_host_checking",
		"hsm_token",
		"hsm_secret",
		"hsm_serial",
	)

	if data.ResetProxyHostKey {
		err = authr.GenerateProxyKey()
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		fields.Add("proxy_private_key")
		fields.Add("proxy_public_key")
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

	err = authr.CommitFields(db, fields)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "authority.change")
	_ = event.PublishDispatch(db, "node.change")

	authr.Json()

	if !showSecret {
		authr.HsmSecret = ""
	}

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
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	authr := &authority.Authority{
		Name:               data.Name,
		Type:               data.Type,
		Algorithm:          data.Algorithm,
		Expire:             data.Expire,
		HostExpire:         data.HostExpire,
		MatchRoles:         data.MatchRoles,
		Roles:              data.Roles,
		KeyIdFormat:        data.KeyIdFormat,
		ProxyHosting:       data.ProxyHosting,
		ProxyHostname:      data.ProxyHostname,
		ProxyPort:          data.ProxyPort,
		HostDomain:         data.HostDomain,
		HostMatches:        data.HostMatches,
		HostSubnets:        data.HostSubnets,
		StrictHostChecking: data.StrictHostChecking,
	}

	err = authr.GeneratePrivateKey()
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

	_ = event.PublishDispatch(db, "authority.change")

	authr.Json()

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

	errData, err := authority.Remove(db, authrId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	_ = event.PublishDispatch(db, "authority.change")
	_ = event.PublishDispatch(db, "node.change")

	c.JSON(200, nil)
}

func authoritiesDelete(c *gin.Context) {
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

	err = authority.RemoveMulti(db, data)
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

	authr.Json()

	authr.HsmSecret = ""

	c.JSON(200, authr)
}

func authoritiesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	if c.Query("names") == "true" {
		authrs, err := authority.GetAllNames(db, &bson.M{})
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(200, authrs)
		return
	}

	page, _ := strconv.ParseInt(c.Query("page"), 10, 0)
	pageCount, _ := strconv.ParseInt(c.Query("page_count"), 10, 0)

	query := bson.M{}

	authorityId, ok := utils.ParseObjectId(c.Query("id"))
	if ok {
		query["_id"] = authorityId
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

	authorities, count, err := authority.GetAllPaged(db, &query,
		page, pageCount)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	data := &authoritiesData{
		Authorities: authorities,
		Count:       count,
	}

	c.JSON(200, data)
}

func authorityPublicKeyGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	authrIdsStr := strings.Split(c.Param("authr_ids"), ",")
	authrIds := []primitive.ObjectID{}

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

	var publicKeys strings.Builder

	authrs, err := authority.GetMulti(db, authrIds)
	if err != nil {
		return
	}

	for _, authr := range authrs {
		publicKeys.WriteString(strings.TrimSpace(authr.PublicKey) + "\n")
	}

	c.String(200, publicKeys.String())
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

	_ = event.PublishDispatch(db, "authority.change")

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

	_ = event.PublishDispatch(db, "authority.change")

	c.Status(200)
}
