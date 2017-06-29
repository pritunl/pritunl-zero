package handlers

import (
	"github.com/gin-gonic/gin"
)

func subscriptionGet(c *gin.Context) {
	c.JSON(200, struct {
		Active   bool   `json:"active"`
		Status   string `json:"status"`
		Plan     string `json:"plan"`
		Quantity int    `json:"quantity"`
	}{
		Active:   false,
		Status:   "active",
		Plan:     "zero0",
		Quantity: 1,
	})
}
