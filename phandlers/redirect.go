package phandlers

import (
	"github.com/gin-gonic/gin"
)

func redirect(c *gin.Context) {
	if c.Request.Header.Get("Upgrade") == "websocket" {
		c.AbortWithStatus(404)
	} else {
		c.Redirect(302, "/")
	}
}
