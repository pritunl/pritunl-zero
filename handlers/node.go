package handlers

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
	Id   bson.ObjectId `json:"id"`
	Name string        `json:"name"`
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

	fields := set.NewSet(
		"name",
	)

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

func nodesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	nodes, err := node.GetAll(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, nodes)
}
