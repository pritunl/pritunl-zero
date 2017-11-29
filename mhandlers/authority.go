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
	Id         bson.ObjectId `json:"id"`
	Name       string        `json:"name"`
	Type       string        `json:"type"`
	Expire     int           `json:"expire"`
	MatchRoles bool          `json:"match_roles"`
	Roles      []string      `json:"roles"`
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
	authr.MatchRoles = data.MatchRoles
	authr.Roles = data.Roles

	fields := set.NewSet(
		"name",
		"type",
		"expire",
		"info",
		"match_roles",
		"roles",
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
		Name:       data.Name,
		Type:       data.Type,
		Expire:     data.Expire,
		MatchRoles: data.MatchRoles,
		Roles:      data.Roles,
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

	c.JSON(200, authr)
}

func authoritysGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	authrs, err := authority.GetAll(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
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
