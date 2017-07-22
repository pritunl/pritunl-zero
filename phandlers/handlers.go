package phandlers

import (
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/middlewear"
	"github.com/pritunl/pritunl-zero/static"
)

var (
	store *static.Store
)

func Register(protocol string, engine *gin.Engine) {
	engine.Use(middlewear.Limiter)
	engine.Use(middlewear.Recovery)
	engine.Use(location.New(location.Config{
		Scheme: protocol,
	}))

	dbGroup := engine.Group("")
	dbGroup.Use(middlewear.Database)

	sessGroup := dbGroup.Group("")
	sessGroup.Use(middlewear.SessionProxy)

	authGroup := sessGroup.Group("")
	authGroup.Use(middlewear.AuthProxy)

	engine.GET("/auth/state", authStateGet)
	dbGroup.POST("/auth/session", authSessionPost)
	dbGroup.GET("/auth/request", authRequestGet)
	dbGroup.GET("/auth/callback", authCallbackGet)
	sessGroup.GET("/logout", logoutGet)

	stre, err := static.NewStore(constants.StaticRoot)
	if err != nil {
		panic(err)
	}
	store = stre

	sessGroup.GET("/", staticIndexGet)
}
