package mhandlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/secret"
	"github.com/pritunl/pritunl-zero/utils"
)

type secretData struct {
	Id      primitive.ObjectID `json:"id"`
	Name    string             `json:"name"`
	Comment string             `json:"comment"`
	Type    string             `json:"type"`
	Key     string             `json:"key"`
	Value   string             `json:"value"`
	Region  string             `json:"region"`
}

func secretPut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &secretData{}

	secrId, ok := utils.ParseObjectId(c.Param("secr_id"))
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

	secr, err := secret.Get(db, secrId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	secr.Name = data.Name
	secr.Comment = data.Comment
	secr.Type = data.Type
	secr.Key = data.Key
	secr.Value = data.Value
	secr.Region = data.Region

	fields := set.NewSet(
		"name",
		"comment",
		"type",
		"key",
		"value",
		"region",
		"public_key",
		"private_key",
	)

	errData, err := secr.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = secr.CommitFields(db, fields)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "secret.change")

	c.JSON(200, secr)
}

func secretPost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &secretData{
		Name: "New Secret",
	}

	err := c.Bind(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	secr := &secret.Secret{
		Name:    data.Name,
		Comment: data.Comment,
		Type:    data.Type,
		Key:     data.Key,
		Value:   data.Value,
		Region:  data.Region,
	}

	errData, err := secr.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = secr.Insert(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "secret.change")

	c.JSON(200, secr)
}

func secretDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)

	secrId, ok := utils.ParseObjectId(c.Param("secr_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := secret.Remove(db, secrId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "secret.change")

	c.JSON(200, nil)
}

func secretGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	secrId, ok := utils.ParseObjectId(c.Param("secr_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	secr, err := secret.Get(db, secrId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if demo.IsDemo() {
		secr.Key = "demo"
		secr.Value = "demo"
	}

	c.JSON(200, secr)
}

func secretsGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	secrs, err := secret.GetAll(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if demo.IsDemo() {
		for _, secr := range secrs {
			secr.Key = "demo"
			secr.Value = "demo"
		}
	}

	c.JSON(200, secrs)
}
