package phandlers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func redirect(c *gin.Context) {
	if strings.ToLower(c.Request.Header.Get("Upgrade")) == "websocket" {
		c.AbortWithStatus(404)
	} else {
		c.Redirect(302, fmt.Sprintf("/?redirect_url=%s",
			url.QueryEscape(c.Request.URL.Path)))
	}
}
