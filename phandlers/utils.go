package phandlers

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/service"
)

func getRedirectPath(c *gin.Context, srvc *service.Service,
	redirectUrl string) string {

	host := ""
	if c.Request.Host != "" {
		host = c.Request.Host
	} else if c.Request.URL.Host != "" {
		host = c.Request.URL.Host
	}

	redirectHost := ""
	if host != "" {
		for _, domain := range srvc.Domains {
			if domain.Domain == host {
				redirectHost = domain.Domain
			}
		}
	}

	if redirectHost != "" && redirectUrl != "" {
		parsed, e := url.Parse(redirectUrl)
		if e == nil {
			parsed.Scheme = "https"
			parsed.Host = redirectHost
			parsed.User = nil
			parsed.Opaque = ""
			return parsed.String()
		}
	}

	return ""
}

type redirectData struct {
	Redirect string `json:"redirect"`
}

func redirectQuery(c *gin.Context, srvc *service.Service, query string) {
	redirect := ""

	vals, err := url.ParseQuery(query)
	if err == nil {
		redirect = getRedirectPath(c, srvc, vals.Get("redirect_url"))
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
