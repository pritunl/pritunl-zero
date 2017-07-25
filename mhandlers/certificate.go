package mhandlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/certificate"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
)

type certificateData struct {
	Id          bson.ObjectId `json:"id"`
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Key         string        `json:"key"`
	Certificate string        `json:"certificate"`
	AcmeAccount string        `json:"acme_account"`
	AcmeDomains []string      `json:"acme_domains"`
}

func certificatePut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &certificateData{}

	certId, ok := utils.ParseObjectId(c.Param("cert_id"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	err := c.Bind(data)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	cert, err := certificate.Get(db, certId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	cert.Name = data.Name
	cert.Type = data.Type
	cert.Key = data.Key
	cert.Certificate = data.Certificate
	cert.AcmeAccount = data.AcmeAccount
	cert.AcmeDomains = data.AcmeDomains

	fields := set.NewSet(
		"name",
		"type",
		"key",
		"certificate",
		"acme_account",
		"acme_domains",
	)

	errData, err := cert.Validate(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = cert.CommitFields(db, fields)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	event.PublishDispatch(db, "certificate.change")

	c.JSON(200, cert)
}

func certificatePost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &certificateData{}

	err := c.Bind(data)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	cert := &certificate.Certificate{
		Name:        data.Name,
		Type:        data.Type,
		Key:         data.Key,
		Certificate: data.Certificate,
		AcmeAccount: data.AcmeAccount,
		AcmeDomains: data.AcmeDomains,
	}

	errData, err := cert.Validate(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if errData != nil {
		c.JSON(400, errData)
		return
	}

	err = cert.Insert(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	event.PublishDispatch(db, "certificate.change")

	c.JSON(200, cert)
}

func certificateDelete(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	certId, ok := utils.ParseObjectId(c.Param("cert_id"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	err := certificate.Remove(db, certId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	event.PublishDispatch(db, "certificate.change")

	c.JSON(200, nil)
}

func certificateGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	certId, ok := utils.ParseObjectId(c.Param("cert_id"))
	if !ok {
		c.AbortWithStatus(400)
		return
	}

	cert, err := certificate.Get(db, certId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, cert)
}

func certificatesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	certs, err := certificate.GetAll(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, certs)
}
