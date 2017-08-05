package phandlers

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func staticIndexGet(c *gin.Context) {
	file := index

	c.Writer.Header().Add("Cache-Control",
		"no-cache, no-store, must-revalidate")
	c.Writer.Header().Add("Pragma", "no-cache")
	c.Writer.Header().Add("Expires", "0")

	if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		c.Writer.Header().Add("Content-Encoding", "gzip")
		c.Data(200, file.Type, file.GzipData)
	} else {
		c.Data(200, file.Type, file.Data)
	}
}
