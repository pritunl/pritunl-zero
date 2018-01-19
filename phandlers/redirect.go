package phandlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
)

func redirect(c *gin.Context) {
	if c.Request.Header.Get("Upgrade") == "websocket" {
		c.AbortWithStatus(404)
	} else {
		c.Redirect(302, fmt.Sprintf("/?redirect_url=%s",
			url.QueryEscape(c.Request.URL.Path)))
	}
}
