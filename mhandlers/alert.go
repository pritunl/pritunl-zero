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
	"github.com/pritunl/pritunl-zero/alert"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/utils"
)

type alertData struct {
	Id       primitive.ObjectID `json:"id"`
	Name     string             `json:"name"`
	Roles    []string           `json:"roles"`
	Resource string             `json:"resource"`
	Level    int                `json:"level"`
	ValueInt int                `json:"value_int"`
	ValueStr string             `json:"value_str"`
}

type alertsData struct {
	Alerts []*alert.Alert `json:"alerts"`
	Count  int64          `json:"count"`
}

func alertPut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &alertData{}

	alertId, ok := utils.ParseObjectId(c.Param("alert_id"))
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

	alrt, err := alert.Get(db, alertId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	alrt.Name = data.Name
	alrt.Roles = data.Roles
	alrt.Resource = data.Resource
	alrt.Level = data.Level
	alrt.ValueInt = data.ValueInt
	alrt.ValueStr = data.ValueStr

	fields := set.NewSet(
		"name",
		"roles",
		"resource",
		"level",
		"value_int",
		"value_str",
	)

	errData, err := alrt.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = alrt.CommitFields(db, fields)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "alert.change")

	c.JSON(200, alrt)
}

func alertPost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &alertData{
		Name:     "New Alert",
		Resource: alert.SystemHighMemory,
		ValueInt: 90,
		Level:    alert.Medium,
	}

	err := c.Bind(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	alrt := &alert.Alert{
		Name:     data.Name,
		Roles:    data.Roles,
		Resource: data.Resource,
		Level:    data.Level,
		ValueInt: data.ValueInt,
		ValueStr: data.ValueStr,
	}

	errData, err := alrt.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = alrt.Insert(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "alert.change")

	c.JSON(200, alrt)
}

func alertDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)

	alertId, ok := utils.ParseObjectId(c.Param("alert_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := alert.Remove(db, alertId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "alert.change")

	c.JSON(200, nil)
}

func alertsDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	dta := []primitive.ObjectID{}

	err := c.Bind(&dta)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	err = alert.RemoveMulti(db, dta)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "alert.change")

	c.JSON(200, nil)
}

func alertsGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	page, _ := strconv.ParseInt(c.Query("page"), 10, 0)
	pageCount, _ := strconv.ParseInt(c.Query("page_count"), 10, 0)

	query := bson.M{}

	alertId, ok := utils.ParseObjectId(c.Query("id"))
	if ok {
		query["_id"] = alertId
	}

	name := strings.TrimSpace(c.Query("name"))
	if name != "" {
		query["$or"] = []*bson.M{
			&bson.M{
				"name": &bson.M{
					"$regex":   fmt.Sprintf(".*%s.*", regexp.QuoteMeta(name)),
					"$options": "i",
				},
			},
		}
	}

	typ := strings.TrimSpace(c.Query("type"))
	if typ != "" {
		query["type"] = typ
	}

	alerts, count, err := alert.GetAllPaged(
		db, &query, page, pageCount)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	dta := &alertsData{
		Alerts: alerts,
		Count:  count,
	}

	c.JSON(200, dta)
}
