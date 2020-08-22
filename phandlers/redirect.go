package phandlers

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
)

func redirect(c *gin.Context) {
	if c.Request.Header.Get("Upgrade") != "" {
		c.AbortWithStatus(404)
	} else {
		c.Redirect(302, fmt.Sprintf("/?redirect_url=%s",
			url.QueryEscape(c.Request.URL.Path)))
	}
}
