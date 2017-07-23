package mhandlers

import (
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/middlewear"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/static"
	"net/http"
)

var (
	store      *static.Store
	fileServer http.Handler
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
	sessGroup.Use(middlewear.Session)

	authGroup := sessGroup.Group("")
	authGroup.Use(middlewear.Auth)

	csrfGroup := authGroup.Group("")
	csrfGroup.Use(middlewear.CsrfToken)

	activeCsrfGroup := csrfGroup.Group("")
	activeCsrfGroup.Use(middlewear.UserActive)

	engine.GET("/check", checkGet)

	authGroup.GET("/csrf", csrfGet)

	engine.GET("/auth/state", authStateGet)
	dbGroup.POST("/auth/session", authSessionPost)
	dbGroup.GET("/auth/request", authRequestGet)
	dbGroup.GET("/auth/callback", authCallbackGet)
	sessGroup.GET("/logout", logoutGet)

	csrfGroup.GET("/event", eventGet)

	activeCsrfGroup.GET("/node", nodesGet)
	activeCsrfGroup.PUT("/node/:node_id", nodePut)
	activeCsrfGroup.DELETE("/node/:node_id", nodeDelete)

	activeCsrfGroup.GET("/service", servicesGet)
	activeCsrfGroup.PUT("/service/:service_id", servicePut)
	activeCsrfGroup.POST("/service", servicePost)
	activeCsrfGroup.DELETE("/service/:service_id", serviceDelete)

	activeCsrfGroup.GET("/settings", settingsGet)
	activeCsrfGroup.PUT("/settings", settingsPut)

	activeCsrfGroup.GET("/subscription", subscriptionGet)
	activeCsrfGroup.GET("/subscription/update", subscriptionUpdateGet)
	activeCsrfGroup.POST("/subscription", subscriptionPost)

	activeCsrfGroup.GET("/user", usersGet)
	activeCsrfGroup.GET("/user/:user_id", userGet)
	activeCsrfGroup.PUT("/user/:user_id", userPut)
	activeCsrfGroup.POST("/user", userPost)
	activeCsrfGroup.DELETE("/user", usersDelete)

	if constants.Production {
		sessGroup.GET("/", staticIndexGet)
		engine.GET("/login", staticLoginGet)
		authGroup.GET("/static/*path", staticGet)
	} else {
		fs := gin.Dir(constants.StaticTestingRoot, false)
		fileServer = http.FileServer(fs)

		sessGroup.GET("/", staticTestingGet)
		engine.GET("/login", staticTestingGet)
		engine.GET("/config.js", staticTestingGet)
		engine.GET("/build.js", staticTestingGet)
		engine.GET("/app/*path", staticTestingGet)
		engine.GET("/dist/*path", staticTestingGet)
		engine.GET("/styles/*path", staticTestingGet)
		engine.GET("/node_modules/*path", staticTestingGet)
		engine.GET("/jspm_packages/*path", staticTestingGet)
	}
}

func init() {
	module := requires.New("mhandlers")
	module.After("settings")

	module.Handler = func() (err error) {
		if constants.Production {
			store, err = static.NewStore(constants.StaticRoot)
			if err != nil {
				return
			}
		}

		return
	}
}
