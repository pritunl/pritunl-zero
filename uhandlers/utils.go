package uhandlers

import (
	"github.com/gin-gonic/gin"
)

type redirectData struct {
	Redirect string `json:"redirect"`
}

func redirectQuery(c *gin.Context, query string) {
	if query != "" {
		c.Redirect(302, "/?"+query)
	} else {
		c.Redirect(302, "/"+query)
	}
}

func redirectQueryJson(c *gin.Context, query string) {
	data := redirectData{
		Redirect: "/",
	}

	if query != "" {
		data.Redirect += "?" + query
	}

	c.JSON(202, data)
}
