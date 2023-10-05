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
		var urlBuf strings.Builder

		path := c.Request.URL.EscapedPath()
		if path == "" || path[0] != '/' {
			urlBuf.WriteByte('/')
		}

		urlBuf.WriteString(path)

		if c.Request.URL.RawQuery != "" {
			urlBuf.WriteByte('?')
			urlBuf.WriteString(c.Request.URL.RawQuery)
		}

		c.Redirect(302, fmt.Sprintf("/?redirect_url=%s",
			url.QueryEscape(urlBuf.String())))
	}
}
