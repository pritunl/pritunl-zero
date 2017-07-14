package handlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
)

type serviceData struct {
	Id    bson.ObjectId `json:"id"`
	Name  string        `json:"name"`
	Roles []string      `json:"roles"`
}

func servicePut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &serviceData{}

	serviceId, ok := utils.ParseObjectId(c.Param("service_id"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	err := c.Bind(data)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	srvce, err := service.Get(db, serviceId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	srvce.Name = data.Name
	srvce.Roles = data.Roles

	fields := set.NewSet(
		"name",
		"roles",
	)

	errData, err := srvce.Validate(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = srvce.CommitFields(db, fields)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	event.PublishDispatch(db, "service.change")

	c.JSON(200, srvce)
}

func servicePost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &serviceData{}

	err := c.Bind(data)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	srvce := &service.Service{}

	if data.Name == "" {
		srvce.Name = "New Service"
	} else {
		srvce.Name = data.Name
	}

	srvce.Roles = data.Roles

	err = srvce.Insert(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	event.PublishDispatch(db, "service.change")

	c.JSON(200, srvce)
}

func serviceDelete(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	serviceId, ok := utils.ParseObjectId(c.Param("service_id"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	err := service.Remove(db, serviceId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	event.PublishDispatch(db, "service.change")

	c.JSON(200, nil)
}

func servicesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	services, err := service.GetAll(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, services)
}
