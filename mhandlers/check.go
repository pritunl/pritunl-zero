package mhandlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/check"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/endpoints"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/utils"
)

type checkData struct {
	Id         bson.ObjectID   `json:"id"`
	Name       string          `json:"name"`
	Roles      []string        `json:"roles"`
	Frequency  int             `json:"frequency"`
	Type       string          `json:"type"`
	Targets    []string        `json:"targets"`
	Timeout    int             `json:"timeout"`
	Method     string          `json:"method"`
	StatusCode int             `json:"status_code"`
	Headers    []*check.Header `json:"headers"`
}

type checksData struct {
	Checks []*check.Check `json:"checks"`
	Count  int64          `json:"count"`
}

func checkPut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &checkData{}

	checkId, ok := utils.ParseObjectId(c.Param("check_id"))
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

	chck, err := check.Get(db, checkId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	chck.Name = data.Name
	chck.Roles = data.Roles
	chck.Frequency = data.Frequency
	chck.Type = data.Type
	chck.Targets = data.Targets
	chck.Timeout = data.Timeout
	chck.Method = data.Method
	chck.StatusCode = data.StatusCode
	chck.Headers = data.Headers

	fields := set.NewSet(
		"name",
		"roles",
		"frequency",
		"type",
		"targets",
		"timeout",
		"method",
		"status_code",
		"headers",
	)

	errData, err := chck.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = chck.CommitFields(db, fields)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "check.change")

	c.JSON(200, chck)
}

func checkPost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &checkData{
		Name: "New Check",
		Type: "http",
	}

	err := c.Bind(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	chck := &check.Check{
		Name:       data.Name,
		Roles:      data.Roles,
		Frequency:  data.Frequency,
		Type:       data.Type,
		Targets:    data.Targets,
		Timeout:    data.Timeout,
		Method:     data.Method,
		StatusCode: data.StatusCode,
		Headers:    data.Headers,
	}

	errData, err := chck.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = chck.Insert(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "check.change")

	c.JSON(200, chck)
}

func checkDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)

	checkId, ok := utils.ParseObjectId(c.Param("check_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := check.Remove(db, checkId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "check.change")

	c.JSON(200, nil)
}

func checksDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	dta := []bson.ObjectID{}

	err := c.Bind(&dta)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	err = check.RemoveMulti(db, dta)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "check.change")

	c.JSON(200, nil)
}

func checksGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	page, _ := strconv.ParseInt(c.Query("page"), 10, 0)
	pageCount, _ := strconv.ParseInt(c.Query("page_count"), 10, 0)

	query := bson.M{}

	checkId, ok := utils.ParseObjectId(c.Query("id"))
	if ok {
		query["_id"] = checkId
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

	role := strings.TrimSpace(c.Query("role"))
	if role != "" {
		if strings.HasPrefix(role, "~") {
			role := role[1:]
			if strings.HasPrefix(role, "!") {
				query["roles"] = &bson.M{
					"$not": &bson.M{
						"$regex": fmt.Sprintf(".*%s.*",
							regexp.QuoteMeta(role[1:])),
						"$options": "i",
					},
				}
			} else {
				query["$or"] = []*bson.M{
					&bson.M{
						"roles": &bson.M{
							"$regex": fmt.Sprintf(".*%s.*",
								regexp.QuoteMeta(role)),
							"$options": "i",
						},
					},
				}
			}
		} else {
			if strings.HasPrefix(role, "!") {
				role = strings.TrimLeft(role, "!")
				query["roles"] = &bson.M{
					"$ne": role,
				}
			} else {
				query["roles"] = role
			}
		}
	}

	checks, count, err := check.GetAllPaged(
		db, &query, page, pageCount)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	dta := &checksData{
		Checks: checks,
		Count:  count,
	}

	c.JSON(200, dta)
}

func checkChartGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	checkId, ok := utils.ParseObjectId(c.Param("check_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	resource := c.Query("resource")

	period, _ := strconv.ParseInt(c.Query("period"), 10, 0)
	if period == 0 {
		period = 1440
	}

	interval, _ := strconv.ParseInt(c.Query("interval"), 10, 0)
	if interval == 0 {
		interval = 24
	}

	chck, err := check.Get(db, checkId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	startTime := time.Now().UTC().Add(time.Duration(-period) * time.Minute)
	endTime := time.Now().UTC()

	data, err := endpoints.GetChart(c, db, chck.Id, resource,
		startTime, endTime, time.Duration(interval)*time.Minute)
	if err != nil {
		return
	}

	chartData := &endpointChartData{
		HasData: len(data) > 0,
		Data:    data,
	}

	c.JSON(200, chartData)
}

func checkLogGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	checkId, ok := utils.ParseObjectId(c.Param("check_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	resource := c.Query("resource")

	chck, err := check.Get(db, checkId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	data, err := endpoints.GetLog(c, db, chck.Id, resource)
	if err != nil {
		return
	}

	c.JSON(200, data)
}
