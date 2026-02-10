package mhandlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/acme"
	"github.com/pritunl/pritunl-zero/certificate"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/utils"
)

type certificateData struct {
	Id          bson.ObjectID `json:"id"`
	Name        string        `json:"name"`
	Comment     string        `json:"comment"`
	Type        string        `json:"type"`
	Key         string        `json:"key"`
	Certificate string        `json:"certificate"`
	AcmeDomains []string      `json:"acme_domains"`
	AcmeType    string        `json:"acme_type"`
	AcmeAuth    string        `json:"acme_auth"`
	AcmeSecret  bson.ObjectID `json:"acme_secret"`
	Refresh     bool          `json:"refresh"`
}

type certificatesData struct {
	Certificates []*certificate.Certificate `json:"certificates"`
	Count        int64                      `json:"count"`
}

func certificatePut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &certificateData{}

	certId, ok := utils.ParseObjectId(c.Param("cert_id"))
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

	cert, err := certificate.Get(db, certId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cert.Name = data.Name
	cert.Comment = data.Comment
	cert.Type = data.Type
	cert.AcmeDomains = data.AcmeDomains
	cert.AcmeType = data.AcmeType
	cert.AcmeAuth = data.AcmeAuth
	cert.AcmeSecret = data.AcmeSecret

	fields := set.NewSet(
		"name",
		"comment",
		"type",
		"acme_domains",
		"acme_type",
		"acme_auth",
		"acme_secret",
		"info",
	)

	if cert.Type != certificate.LetsEncrypt {
		cert.Key = data.Key
		fields.Add("key")
		cert.Certificate = data.Certificate
		fields.Add("certificate")
	}

	errData, err := cert.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = cert.CommitFields(db, fields)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if cert.Type == certificate.LetsEncrypt {
		acme.RenewBackground(cert, data.Refresh)
	}

	_ = event.PublishDispatch(db, "certificate.change")

	c.JSON(200, cert)
}

func certificatePost(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &certificateData{
		Name: "New Certificate",
	}

	err := c.Bind(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	cert := &certificate.Certificate{
		Name:        data.Name,
		Comment:     data.Comment,
		Type:        data.Type,
		AcmeDomains: data.AcmeDomains,
		AcmeType:    data.AcmeType,
		AcmeAuth:    data.AcmeAuth,
		AcmeSecret:  data.AcmeSecret,
	}

	if cert.Type != certificate.LetsEncrypt {
		cert.Key = data.Key
		cert.Certificate = data.Certificate
	}

	errData, err := cert.Validate(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = cert.Insert(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if cert.Type == certificate.LetsEncrypt {
		acme.RenewBackground(cert, false)
	}

	_ = event.PublishDispatch(db, "certificate.change")

	c.JSON(200, cert)
}

func certificateDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)

	certId, ok := utils.ParseObjectId(c.Param("cert_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	err := certificate.Remove(db, certId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "certificate.change")
	_ = event.PublishDispatch(db, "node.change")

	c.JSON(200, nil)
}

func certificatesDelete(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := []bson.ObjectID{}

	err := c.Bind(&data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	err = certificate.RemoveMulti(db, data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	event.PublishDispatch(db, "certificate.change")

	c.JSON(200, nil)
}

func certificateGet(c *gin.Context) {
	if demo.IsDemo() {
		cert := demo.Certificates[0]
		c.JSON(200, cert)
		return
	}

	db := c.MustGet("db").(*database.Database)

	certId, ok := utils.ParseObjectId(c.Param("cert_id"))
	if !ok {
		utils.AbortWithStatus(c, 400)
		return
	}

	cert, err := certificate.Get(db, certId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	if demo.IsDemo() {
		cert.Key = "demo"
		cert.AcmeAccount = "demo"
	}

	c.JSON(200, cert)
}

func certificatesGet(c *gin.Context) {
	if demo.IsDemo() {
		data := &certificatesData{
			Certificates: demo.Certificates,
			Count:        int64(len(demo.Certificates)),
		}

		c.JSON(200, data)
		return
	}

	db := c.MustGet("db").(*database.Database)

	if c.Query("names") == "true" {
		certs, err := certificate.GetAllNames(db, &bson.M{})
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.JSON(200, certs)
		return
	}

	page, _ := strconv.ParseInt(c.Query("page"), 10, 0)
	pageCount, _ := strconv.ParseInt(c.Query("page_count"), 10, 0)

	query := bson.M{}

	certificateId, ok := utils.ParseObjectId(c.Query("id"))
	if ok {
		query["_id"] = certificateId
	}

	name := strings.TrimSpace(c.Query("name"))
	if name != "" {
		query["name"] = &bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", regexp.QuoteMeta(name)),
			"$options": "i",
		}
	}

	organization, ok := utils.ParseObjectId(c.Query("organization"))
	if ok {
		query["organization"] = organization
	}

	comment := strings.TrimSpace(c.Query("comment"))
	if comment != "" {
		query["comment"] = &bson.M{
			"$regex":   fmt.Sprintf(".*%s.*", comment),
			"$options": "i",
		}
	}

	certs, count, err := certificate.GetAllPaged(db, &query, page, pageCount)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	data := &certificatesData{
		Certificates: certs,
		Count:        count,
	}

	c.JSON(200, data)
}
