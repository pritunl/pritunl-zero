package uhandlers

import (
	"github.com/gin-gonic/gin"
	"net/url"
)

type redirectData struct {
	Redirect string `json:"redirect"`
}

func redirectQuery(c *gin.Context, query string) {
	if query != "" {
		vals, _ := url.ParseQuery(query)
		if vals.Get("redirect") == "ssh-validate" {
			vals.Del("redirect")
			c.Redirect(302, "/ssh/validate?"+vals.Encode())
			return
		}
	}
	c.Redirect(302, "/")
}

func redirectQueryJson(c *gin.Context, query string) {
	data := redirectData{
		Redirect: "/",
	}

	if query != "" {
		vals, _ := url.ParseQuery(query)
		if vals.Get("redirect") == "ssh-validate" {
			vals.Del("redirect")
			data.Redirect = "/ssh/validate?" + vals.Encode()
		}
	}

	c.JSON(202, data)
}
