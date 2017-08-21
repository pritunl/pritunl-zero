package phandlers

import (
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/middlewear"
	"github.com/pritunl/pritunl-zero/proxy"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/static"
	"path/filepath"
)

var (
	index *static.File
	logo  *static.File
)

func Register(prxy *proxy.Proxy, protocol string, engine *gin.Engine) {
	engine.Use(middlewear.Limiter)
	engine.Use(middlewear.Counter)
	engine.Use(middlewear.Recovery)
	engine.Use(location.New(location.Config{
		Scheme: protocol,
	}))

	engine.Use(func(c *gin.Context) {
		var srvc *service.Service
		host := prxy.Hosts[c.Request.Host]
		if host != nil {
			srvc = host.Service
		}
		c.Set("service", srvc)
	})

	engine.NoRoute(redirect)

	dbGroup := engine.Group("")
	dbGroup.Use(middlewear.Database)

	sessGroup := dbGroup.Group("")
	sessGroup.Use(middlewear.SessionProxy)

	engine.GET("/auth/state", authStateGet)
	dbGroup.POST("/auth/session", authSessionPost)
	dbGroup.GET("/auth/request", authRequestGet)
	dbGroup.GET("/auth/callback", authCallbackGet)
	sessGroup.GET("/logout", logoutGet)

	engine.GET("/", staticIndexGet)
	engine.GET("/logo.png", staticLogoGet)
}

func init() {
	module := requires.New("phandlers")
	module.After("settings")

	module.Handler = func() (err error) {
		root := ""
		if constants.Production {
			root = config.StaticRoot
		} else {
			root = config.StaticTestingRoot
		}

		index, err = static.NewFile(filepath.Join(root, "login.html"))
		if err != nil {
			return
		}

		logo, err = static.NewFile(filepath.Join(root, "logo.png"))
		if err != nil {
			return
		}

		return
	}
}
