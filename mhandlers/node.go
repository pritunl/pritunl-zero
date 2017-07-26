package mhandlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
)

type nodeData struct {
	Id               bson.ObjectId   `json:"id"`
	Name             string          `json:"name"`
	Type             string          `json:"type"`
	Port             int             `json:"port"`
	Protocol         string          `json:"protocol"`
	Certificate      bson.ObjectId   `json:"certificate"`
	ManagementDomain string          `json:"management_domain"`
	Services         []bson.ObjectId `json:"services"`
}

func nodePut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &nodeData{}

	nodeId, ok := utils.ParseObjectId(c.Param("node_id"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	err := c.Bind(data)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	nde, err := node.Get(db, nodeId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	nde.Name = data.Name
	nde.Type = data.Type
	nde.Port = data.Port
	nde.Protocol = data.Protocol
	nde.Certificate = data.Certificate
	nde.ManagementDomain = data.ManagementDomain
	nde.Services = data.Services

	fields := set.NewSet(
		"name",
		"type",
		"port",
		"protocol",
		"certificate",
		"management_domain",
		"services",
	)

	errData, err := nde.Validate(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = nde.CommitFields(db, fields)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	event.PublishDispatch(db, "node.change")

	c.JSON(200, nde)
}

func nodeDelete(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	nodeId, ok := utils.ParseObjectId(c.Param("node_id"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	err := node.Remove(db, nodeId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	event.PublishDispatch(db, "node.change")

	c.JSON(200, nil)
}

func nodeGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	nodeId, ok := utils.ParseObjectId(c.Param("node_id"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	nde, err := node.Get(db, nodeId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, nde)
}

func nodesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	nodes, err := node.GetAll(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, nodes)
}
