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
	engine.Use(middlewear.Counter)
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

	engine.GET("/auth/state", authStateGet)
	dbGroup.POST("/auth/session", authSessionPost)
	dbGroup.GET("/auth/request", authRequestGet)
	dbGroup.GET("/auth/callback", authCallbackGet)
	sessGroup.GET("/logout", logoutGet)

	csrfGroup.GET("/certificate", certificatesGet)
	csrfGroup.GET("/certificate/:cert_id", certificateGet)
	csrfGroup.PUT("/certificate/:cert_id", certificatePut)
	csrfGroup.POST("/certificate", certificatePost)
	csrfGroup.DELETE("/certificate/:cert_id", certificateDelete)

	engine.GET("/check", checkGet)

	authGroup.GET("/csrf", csrfGet)

	csrfGroup.GET("/event", eventGet)

	csrfGroup.GET("/log", logsGet)
	csrfGroup.GET("/log/:log_id", logGet)

	csrfGroup.GET("/node", nodesGet)
	csrfGroup.GET("/node/:node_id", nodeGet)
	csrfGroup.PUT("/node/:node_id", nodePut)
	csrfGroup.DELETE("/node/:node_id", nodeDelete)

	csrfGroup.GET("/policy", policiesGet)
	csrfGroup.GET("/policy/:policy_id", policyGet)
	csrfGroup.PUT("/policy/:policy_id", policyPut)
	csrfGroup.POST("/policy", policyPost)
	csrfGroup.DELETE("/policy/:policy_id", policyDelete)

	csrfGroup.GET("/service", servicesGet)
	csrfGroup.PUT("/service/:service_id", servicePut)
	csrfGroup.POST("/service", servicePost)
	csrfGroup.DELETE("/service/:service_id", serviceDelete)

	csrfGroup.GET("/session/:user_id", sessionsGet)
	csrfGroup.DELETE("/session/:session_id", sessionDelete)

	csrfGroup.GET("/settings", settingsGet)
	csrfGroup.PUT("/settings", settingsPut)

	csrfGroup.GET("/subscription", subscriptionGet)
	csrfGroup.GET("/subscription/update", subscriptionUpdateGet)
	csrfGroup.POST("/subscription", subscriptionPost)

	csrfGroup.PUT("/theme", themePut)

	csrfGroup.GET("/user", usersGet)
	csrfGroup.GET("/user/:user_id", userGet)
	csrfGroup.PUT("/user/:user_id", userPut)
	csrfGroup.POST("/user", userPost)
	csrfGroup.DELETE("/user", usersDelete)

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
