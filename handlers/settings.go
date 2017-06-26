package handlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/settings"
)

type settingsData struct {
	ElasticAddress string `json:"elastic_address"`
}

func getSettingsData() *settingsData {
	return &settingsData{
		ElasticAddress: settings.Elastic.Address,
	}
}

func settingsGet(c *gin.Context) {
	data := getSettingsData()
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
	settings.Commit(db, settings.Elastic, set.NewSet(
		"address",
	))

	event.PublishDispatch(db, "settings.change")

	data = getSettingsData()
	c.JSON(200, data)
}
