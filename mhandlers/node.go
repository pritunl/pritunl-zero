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
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/utils"
)

type nodeData struct {
	Id                   primitive.ObjectID   `json:"id"`
	Name                 string               `json:"name"`
	Type                 string               `json:"type"`
	Port                 int                  `json:"port"`
	NoRedirectServer     bool                 `json:"no_redirect_server"`
	Protocol             string               `json:"protocol"`
	Certificates         []primitive.ObjectID `json:"certificates"`
	ManagementDomain     string               `json:"management_domain"`
	UserDomain           string               `json:"user_domain"`
	Services             []primitive.ObjectID `json:"services"`
	Authorities          []primitive.ObjectID `json:"authorities"`
	ForwardedForHeader   string               `json:"forwarded_for_header"`
	ForwardedProtoHeader string               `json:"forwarded_proto_header"`
}

func nodePut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &nodeData{}

	nodeId, ok := utils.ParseObjectId(c.Param("node_id"))
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

	nde, err := node.Get(db, nodeId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	nde.Name = data.Name
	nde.Type = data.Type
	nde.Port = data.Port
	nde.NoRedirectServer = data.NoRedirectServer
	nde.Protocol = data.Protocol
	nde.Certificates = data.Certificates
	nde.ManagementDomain = data.ManagementDomain
	nde.UserDomain = data.UserDomain
	nde.Services = data.Services
	nde.Authorities = data.Authorities
	nde.ForwardedForHeader = data.ForwardedForHeader
	nde.ForwardedProtoHeader = data.ForwardedProtoHeader

	fields := set.NewSet(
		"name",
		"type",
		"port",
		"no_redirect_server",
		"protocol",
		"certificates",
		"management_domain",
		"user_domain",
		"services",
		"authorities",
		"forwarded_for_header",
		"forwarded_proto_header",
	)

	errData, err := nde.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = nde.CommitFields(db, fields)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "node.change")

	c.JSON(200, nde)
}

func nodeDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)

	nodeId, ok := utils.ParseObjectId(c.Param("node_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := node.Remove(db, nodeId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "node.change")

	c.JSON(200, nil)
}

func nodeGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	nodeId, ok := utils.ParseObjectId(c.Param("node_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	nde, err := node.Get(db, nodeId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if demo.IsDemo() {
		nde.RequestsMin = 32
		nde.Memory = 25.0
		nde.Load1 = 10.0
		nde.Load5 = 15.0
		nde.Load15 = 20.0
	}

	c.JSON(200, nde)
}

func nodesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	nodes, err := node.GetAll(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if demo.IsDemo() {
		for _, nde := range nodes {
			nde.RequestsMin = 32
			nde.Memory = 25.0
			nde.Load1 = 10.0
			nde.Load5 = 15.0
			nde.Load15 = 20.0
		}
	}

	c.JSON(200, nodes)
}
