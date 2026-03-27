package phandlers

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

type redirectData struct {
	Redirect string `json:"redirect"`
}

func redirectQuery(c *gin.Context, query string) {
	redirect := ""

	vals, err := url.ParseQuery(query)
	if err == nil {
		redirect = vals.Get("redirect_url")
	}

	if redirect != "" {
		// Prevent open redirect by ensuring the URL is a relative path
		parsed, err := url.Parse(redirect)
		if err != nil || parsed.Host != "" || parsed.Scheme != "" {
			redirect = "/"
		}
		c.Redirect(302, redirect)
	} else {
		c.Redirect(302, "/")
	}
}

func redirectJson(c *gin.Context, redirect string) {
	if redirect == "" {
		redirect = "/"
	}

	data := redirectData{
		Redirect: redirect,
	}

	c.JSON(202, data)
}
