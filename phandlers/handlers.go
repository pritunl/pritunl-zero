package phandlers

import (
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/middlewear"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/static"
	"path/filepath"
)

var (
	index *static.File
)

func Register(protocol string, engine *gin.Engine) {
	engine.Use(middlewear.Limiter)
	engine.Use(middlewear.Counter)
	engine.Use(middlewear.Recovery)
	engine.Use(location.New(location.Config{
		Scheme: protocol,
	}))
	engine.Use(middlewear.Service)

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

	sessGroup.GET("/", staticIndexGet)
}

func init() {
	module := requires.New("phandlers")
	module.After("settings")

	module.Handler = func() (err error) {
		root := ""
		if constants.Production {
			root = constants.StaticRoot
		} else {
			root = constants.StaticTestingRoot
		}

		index, err = static.NewFile(filepath.Join(root, "login.html"))
		if err != nil {
			return
		}

		return
	}
}
