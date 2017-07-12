package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"gopkg.in/mgo.v2/bson"
)

type serviceData struct {
	Id   bson.ObjectId `json:"id"`
	Name string        `json:"name"`
}

var temp = []*serviceData{
	&serviceData{
		Id:   bson.NewObjectId(),
		Name: "test1",
	},
	&serviceData{
		Id:   bson.NewObjectId(),
		Name: "test2",
	},
}

func servicesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	_ = db

	c.JSON(200, temp)
}
