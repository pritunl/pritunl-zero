package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/static"
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
	sess := c.MustGet("session").(*session.Session)
	if sess == nil {
		c.Redirect(302, "/login")
		return
	}

	staticPath(c, "/index.html")
}

func staticLoginGet(c *gin.Context) {
	staticPath(c, "/login.html")
}

func staticGet(c *gin.Context) {
	staticPath(c, "/static"+c.Params.ByName("path"))
}

func staticTestingGet(c *gin.Context) {
	pth := c.Params.ByName("path")
	if pth == "" {
		if c.Request.URL.Path == "/config.js" {
			pth = "config.js"
		} else if c.Request.URL.Path == "/build.js" {
			pth = "build.js"
		} else if c.Request.URL.Path == "/login" {
			c.Request.URL.Path = "/login.html"
			pth = "login.html"
		} else {
			sess := c.MustGet("session").(*session.Session)
			if sess == nil {
				c.Redirect(302, "/login")
				return
			}

			pth = "index.html"
		}
	}

	if strings.HasPrefix(c.Request.URL.Path, "/node_modules/") ||
		strings.HasPrefix(c.Request.URL.Path, "/jspm_packages/") {

		c.Writer.Header().Add("Cache-Control", "public, max-age=86400")
	} else {
		c.Writer.Header().Add("Cache-Control",
			"no-cache, no-store, must-revalidate")
		c.Writer.Header().Add("Pragma", "no-cache")
		c.Writer.Header().Add("Expires", "0")
	}

	c.Writer.Header().Add("Content-Type", static.GetMimeType(pth))
	fileServer.ServeHTTP(c.Writer, c.Request)
}
