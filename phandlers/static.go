package phandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/constants"
	"strings"
)

func staticPath(c *gin.Context, pth string) {
	pth = constants.StaticRoot + pth

	file, ok := store.Files[pth]
	if !ok {
		c.AbortWithStatus(404)
		return
	}

	if constants.StaticCache {
		c.Writer.Header().Add("Cache-Control", "public, max-age=86400")
		c.Writer.Header().Add("ETag", file.Hash)
	} else {
		c.Writer.Header().Add("Cache-Control",
			"no-cache, no-store, must-revalidate")
		c.Writer.Header().Add("Pragma", "no-cache")
		c.Writer.Header().Add("Expires", "0")
	}

	if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		c.Writer.Header().Add("Content-Encoding", "gzip")
		c.Data(200, file.Type, file.GzipData)
	} else {
		c.Data(200, file.Type, file.Data)
	}
}

func staticIndexGet(c *gin.Context) {
	staticPath(c, "/login.html")
}
