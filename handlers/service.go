package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
)

type serviceData struct {
	Id   bson.ObjectId `json:"id"`
	Name string        `json:"name"`
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
		srvce.Name = utils.RandName()
	} else {
		srvce.Name = data.Name
	}

	err = srvce.Insert(db)
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
