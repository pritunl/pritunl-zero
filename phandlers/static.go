package phandlers

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func staticIndexGet(c *gin.Context) {
	c.Writer.Header().Add("Cache-Control",
		"no-cache, no-store, must-revalidate")
	c.Writer.Header().Add("Pragma", "no-cache")
	c.Writer.Header().Add("Expires", "0")

	if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		c.Writer.Header().Add("Content-Encoding", "gzip")
		c.Data(200, index.Type, index.GzipData)
	} else {
		c.Data(200, index.Type, index.Data)
	}
}

func staticLogoGet(c *gin.Context) {
	c.Writer.Header().Add("Cache-Control",
		"no-cache, no-store, must-revalidate")
	c.Writer.Header().Add("Pragma", "no-cache")
	c.Writer.Header().Add("Expires", "0")

	if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		c.Writer.Header().Add("Content-Encoding", "gzip")
		c.Data(200, logo.Type, logo.GzipData)
	} else {
		c.Data(200, logo.Type, logo.Data)
	}
}
