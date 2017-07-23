package phandlers

import (
	"github.com/gin-gonic/gin"
)

func redirect(c *gin.Context) {
	c.Redirect(302, "/")
}
