package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/service"
	"gopkg.in/mgo.v2/bson"
)

type serviceData struct {
	Id   bson.ObjectId `json:"id"`
	Name string        `json:"name"`
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
