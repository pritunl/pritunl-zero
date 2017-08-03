package mhandlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/acme"
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
		utils.AbortWithStatus(c, 400)
		return
	}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cert, err := certificate.Get(db, certId)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cert.Name = data.Name
	cert.Type = data.Type
	cert.AcmeAccount = data.AcmeAccount
	cert.AcmeDomains = data.AcmeDomains

	fields := set.NewSet(
		"name",
		"type",
		"acme_account",
		"acme_domains",
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
		err = acme.Update(db, cert)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		err = acme.Renew(db, cert)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}
	}

	event.PublishDispatch(db, "certificate.change")

	c.JSON(200, cert)
}

func certificatePost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &certificateData{
		Name: "New Certificate",
	}

	err := c.Bind(data)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	cert := &certificate.Certificate{
		Name:        data.Name,
		Type:        data.Type,
		AcmeAccount: data.AcmeAccount,
		AcmeDomains: data.AcmeDomains,
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
		err = acme.Update(db, cert)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		err = acme.Renew(db, cert)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}
	}

	event.PublishDispatch(db, "certificate.change")

	c.JSON(200, cert)
}

func certificateDelete(c *gin.Context) {
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

	event.PublishDispatch(db, "certificate.change")

	c.JSON(200, nil)
}

func certificateGet(c *gin.Context) {
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

	c.JSON(200, cert)
}

func certificatesGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	certs, err := certificate.GetAll(db)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	c.JSON(200, certs)
}
