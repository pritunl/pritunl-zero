package handlers

import (
	"github.com/gin-gonic/gin"
)

func checkGet(c *gin.Context) {
	c.String(200, "ok")
}
