package handlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/settings"
)

type settingsData struct {
	ElasticAddress string `json:"elastic_address"`
}

func settingsGet(c *gin.Context) {
	data := &settingsData{
		ElasticAddress: settings.Elastic.Address,
	}

	c.JSON(200, data)
}

func settingsPut(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &settingsData{}

	err := c.Bind(data)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	settings.Elastic.Address = data.ElasticAddress
	fields := set.NewSet(
		"elastic_address",
	)
	settings.Commit(db, "elastic", fields)

	c.JSON(200, data)
}
