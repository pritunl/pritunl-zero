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
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/utils"
)

type serviceData struct {
	Id                primitive.ObjectID       `json:"id"`
	Name              string                   `json:"name"`
	Type              string                   `json:"type"`
	ShareSession      bool                     `json:"share_session"`
	LogoutPath        string                   `json:"logout_path"`
	WebSockets        bool                     `json:"websockets"`
	DisableCsrfCheck  bool                     `json:"disable_csrf_check"`
	ClientAuthority   primitive.ObjectID       `json:"client_authority"`
	Domains           []*service.Domain        `json:"domains"`
	Roles             []string                 `json:"roles"`
	Servers           []*service.Server        `json:"servers"`
	WhitelistNetworks []string                 `json:"whitelist_networks"`
	WhitelistPaths    []*service.WhitelistPath `json:"whitelist_paths"`
	WhitelistOptions  bool                     `json:"whitelist_options"`
}

type servicesData struct {
	Services []*service.Service `json:"services"`
	Count    int64              `json:"count"`
}

func servicePut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &serviceData{}

	serviceId, ok := utils.ParseObjectId(c.Param("service_id"))
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

	srvce, err := service.Get(db, serviceId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	srvce.Name = data.Name
	srvce.Type = data.Type
	srvce.ShareSession = data.ShareSession
	srvce.LogoutPath = data.LogoutPath
	srvce.WebSockets = data.WebSockets
	srvce.DisableCsrfCheck = data.DisableCsrfCheck
	srvce.ClientAuthority = data.ClientAuthority
	srvce.Domains = data.Domains
	srvce.Roles = data.Roles
	srvce.Servers = data.Servers
	srvce.WhitelistNetworks = data.WhitelistNetworks
	srvce.WhitelistPaths = data.WhitelistPaths
	srvce.WhitelistOptions = data.WhitelistOptions

	fields := set.NewSet(
		"name",
		"type",
		"share_session",
		"logout_path",
		"websockets",
		"disable_csrf_check",
		"client_authority",
		"domains",
		"roles",
		"servers",
		"whitelist_networks",
		"whitelist_paths",
		"whitelist_options",
	)

	errData, err := srvce.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = srvce.CommitFields(db, fields)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "service.change")

	c.JSON(200, srvce)
}

func servicePost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &serviceData{
		Name:         "New Service",
		ShareSession: true,
	}

	err := c.Bind(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	srvce := &service.Service{
		Name:              data.Name,
		Type:              data.Type,
		ShareSession:      data.ShareSession,
		LogoutPath:        data.LogoutPath,
		WebSockets:        data.WebSockets,
		DisableCsrfCheck:  data.DisableCsrfCheck,
		ClientAuthority:   data.ClientAuthority,
		Roles:             data.Roles,
		Domains:           data.Domains,
		Servers:           data.Servers,
		WhitelistNetworks: data.WhitelistNetworks,
		WhitelistPaths:    data.WhitelistPaths,
		WhitelistOptions:  data.WhitelistOptions,
	}

	errData, err := srvce.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = srvce.Insert(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "service.change")

	c.JSON(200, srvce)
}

func serviceDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)

	serviceId, ok := utils.ParseObjectId(c.Param("service_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := service.Remove(db, serviceId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "service.change")

	c.JSON(200, nil)
}

func servicesDelete(c *gin.Context) {
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

	err = service.RemoveMulti(db, dta)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "service.change")

	c.JSON(200, nil)
}

func servicesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	serviceNames := c.Query("service_names")
	if serviceNames == "true" {
		insts, err := service.GetAllName(db)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(200, insts)
	} else {
		page, _ := strconv.ParseInt(c.Query("page"), 10, 0)
		pageCount, _ := strconv.ParseInt(c.Query("page_count"), 10, 0)

		query := bson.M{}

		serviceId, ok := utils.ParseObjectId(c.Query("id"))
		if ok {
			query["_id"] = serviceId
		}

		name := strings.TrimSpace(c.Query("name"))
		if name != "" {
			query["$or"] = []*bson.M{
				&bson.M{
					"name": &bson.M{
						"$regex": fmt.Sprintf(".*%s.*",
							regexp.QuoteMeta(name)),
						"$options": "i",
					},
				},
			}
		}

		typ := strings.TrimSpace(c.Query("type"))
		if typ != "" {
			query["type"] = typ
		}

		organization, ok := utils.ParseObjectId(c.Query("organization"))
		if ok {
			query["organization"] = organization
		}

		services, count, err := service.GetAllPaged(
			db, &query, page, pageCount)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		dta := &servicesData{
			Services: services,
			Count:    count,
		}

		c.JSON(200, dta)
	}
}
