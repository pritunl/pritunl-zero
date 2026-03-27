package phandlers

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func getRedirectPath(redirectUrl string) string {
	if redirectUrl != "" {
		parsed, e := url.Parse(redirectUrl)
		if e == nil {
			return parsed.Path
		}
	}
	return ""
}

type redirectData struct {
	Redirect string `json:"redirect"`
}

func redirectQuery(c *gin.Context, query string) {
	redirect := ""

	vals, err := url.ParseQuery(query)
	if err == nil {
		redirect = getRedirectPath(vals.Get("redirect_url"))
	}

	if redirect != "" {
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
