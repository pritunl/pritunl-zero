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
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/utils"
)

type nodeData struct {
	Id                   bson.ObjectID   `json:"id"`
	Name                 string          `json:"name"`
	Type                 string          `json:"type"`
	Port                 int             `json:"port"`
	NoRedirectServer     bool            `json:"no_redirect_server"`
	Protocol             string          `json:"protocol"`
	Certificates         []bson.ObjectID `json:"certificates"`
	ManagementDomain     string          `json:"management_domain"`
	UserDomain           string          `json:"user_domain"`
	EndpointDomain       string          `json:"endpoint_domain"`
	WebauthnDomain       string          `json:"webauthn_domain"`
	Services             []bson.ObjectID `json:"services"`
	Authorities          []bson.ObjectID `json:"authorities"`
	ForwardedForHeader   string          `json:"forwarded_for_header"`
	ForwardedProtoHeader string          `json:"forwarded_proto_header"`
}

type nodesData struct {
	Nodes []*node.Node `json:"nodes"`
	Count int64        `json:"count"`
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
	nde.EndpointDomain = data.EndpointDomain
	nde.WebauthnDomain = data.WebauthnDomain
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
		"endpoint_domain",
		"webauthn_domain",
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

	_ = event.PublishDispatch(db, "node.change")

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

	_ = event.PublishDispatch(db, "node.change")

	c.JSON(200, nil)
}

func nodeGet(c *gin.Context) {
	if demo.IsDemo() {
		nde := demo.Nodes[0]
		nde.Timestamp = time.Now()
		c.JSON(200, nde)
		return
	}

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

	c.JSON(200, nde)
}

func nodesGet(c *gin.Context) {
	if demo.IsDemo() {
		for _, nde := range demo.Nodes {
			nde.Timestamp = time.Now()
		}

		data := &nodesData{
			Nodes: demo.Nodes,
			Count: int64(len(demo.Nodes)),
		}

		c.JSON(200, data)
		return
	}

	db := c.MustGet("db").(*database.Database)
	page, _ := strconv.ParseInt(c.Query("page"), 10, 0)
	pageCount, _ := strconv.ParseInt(c.Query("page_count"), 10, 0)

	query := bson.M{}

	nodeId, ok := utils.ParseObjectId(c.Query("id"))
	if ok {
		query["_id"] = nodeId
	}

	name := strings.TrimSpace(c.Query("name"))
	if name != "" {
		query["name"] = &bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", regexp.QuoteMeta(name)),
			"$options": "i",
		}
	}

	nodes, count, err := node.GetAllPaged(db, &query, page, pageCount)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	data := &nodesData{
		Nodes: nodes,
		Count: count,
	}

	c.JSON(200, data)
}
