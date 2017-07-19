package mhandlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/subscription"
	"strings"
)

type subscriptionPostData struct {
	License string `json:"license"`
}

func subscriptionGet(c *gin.Context) {
	c.JSON(200, subscription.Subscription)
}

func subscriptionUpdateGet(c *gin.Context) {
	errData, err := subscription.Update()
	if err != nil {
		if errData != nil {
			c.JSON(400, errData)
		} else {
			c.AbortWithError(500, err)
		}
		return
	}

	c.JSON(200, subscription.Subscription)
}

func subscriptionPost(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	data := &subscriptionPostData{}

	err := c.Bind(data)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	license := strings.TrimSpace(data.License)
	license = strings.Replace(license, "BEGIN LICENSE", "", 1)
	license = strings.Replace(license, "END LICENSE", "", 1)
	license = strings.Replace(license, "-", "", -1)
	license = strings.Replace(license, " ", "", -1)
	license = strings.Replace(license, "\n", "", -1)

	settings.System.License = license

	errData, err := subscription.Update()
	if err != nil {
		settings.System.License = ""
		if errData != nil {
			c.JSON(400, errData)
		} else {
			c.AbortWithError(500, err)
		}
		return
	}

	err = settings.Commit(db, settings.System, set.NewSet(
		"license",
	))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	event.PublishDispatch(db, "subscription.change")
	event.PublishDispatch(db, "settings.change")

	c.JSON(200, subscription.Subscription)
}
