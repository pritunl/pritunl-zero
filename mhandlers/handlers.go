package mhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/handlers"
	"github.com/pritunl/pritunl-zero/middlewear"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/static"
)

var (
	store      *static.Store
	fileServer http.Handler
	pushFiles  []string
)

func Register(engine *gin.Engine) {
	engine.Use(middlewear.Limiter)
	engine.Use(middlewear.Counter)
	engine.Use(middlewear.Recovery)
	engine.Use(middlewear.Headers)

	dbGroup := engine.Group("")
	dbGroup.Use(middlewear.Database)

	sessGroup := dbGroup.Group("")
	sessGroup.Use(middlewear.SessionAdmin)

	authGroup := sessGroup.Group("")
	authGroup.Use(middlewear.AuthAdmin)

	csrfGroup := authGroup.Group("")
	csrfGroup.Use(middlewear.CsrfToken)

	engine.NoRoute(middlewear.NotFound)

	csrfGroup.GET("/audit/:user_id", auditsGet)

	csrfGroup.GET("/alert", alertsGet)
	csrfGroup.PUT("/alert/:alert_id", alertPut)
	csrfGroup.POST("/alert", alertPost)
	csrfGroup.DELETE("/alert", alertsDelete)
	csrfGroup.DELETE("/alert/:alert_id", alertDelete)

	engine.GET("/auth/state", authStateGet)
	dbGroup.POST("/auth/session", authSessionPost)
	dbGroup.POST("/auth/secondary", authSecondaryPost)
	dbGroup.GET("/auth/request", authRequestGet)
	dbGroup.GET("/auth/callback", authCallbackGet)
	dbGroup.GET("/auth/webauthn/request", authWanRequestGet)
	dbGroup.POST("/auth/webauthn/respond", authWanRespondPost)
	dbGroup.GET("/auth/webauthn/register", authWanRegisterGet)
	dbGroup.POST("/auth/webauthn/register", authWanRegisterPost)
	sessGroup.GET("/logout", logoutGet)

	csrfGroup.GET("/authority", authoritysGet)
	csrfGroup.GET("/authority/:authr_id", authorityGet)
	csrfGroup.PUT("/authority/:authr_id", authorityPut)
	csrfGroup.POST("/authority", authorityPost)
	csrfGroup.DELETE("/authority/:authr_id", authorityDelete)
	csrfGroup.POST("/authority/:authr_id/token", authorityTokenPost)
	csrfGroup.DELETE("/authority/:authr_id/token/:token",
		authorityTokenDelete)
	dbGroup.GET("/ssh_public_key/:authr_ids", authorityPublicKeyGet)

	csrfGroup.GET("/certificate", certificatesGet)
	csrfGroup.GET("/certificate/:cert_id", certificateGet)
	csrfGroup.PUT("/certificate/:cert_id", certificatePut)
	csrfGroup.POST("/certificate", certificatePost)
	csrfGroup.DELETE("/certificate/:cert_id", certificateDelete)

	engine.GET("/check", checkGet)

	authGroup.GET("/csrf", csrfGet)

	csrfGroup.GET("/device/:user_id", devicesGet)
	csrfGroup.PUT("/device/:device_id", devicePut)
	csrfGroup.POST("/device", devicePost)
	csrfGroup.DELETE("/device/:device_id", deviceDelete)
	csrfGroup.POST("/device/:resource_id/:method", deviceMethodPost)
	csrfGroup.GET("/device/:user_id/webauthn/register", deviceWanRegisterGet)
	csrfGroup.POST("/device/:resource_id/webauthn/register",
		deviceWanRegisterPost)

	csrfGroup.GET("/endpoint", endpointsGet)
	csrfGroup.PUT("/endpoint/:endpoint_id", endpointPut)
	csrfGroup.POST("/endpoint", endpointPost)
	csrfGroup.DELETE("/endpoint", endpointsDelete)
	csrfGroup.DELETE("/endpoint/:endpoint_id", endpointDelete)
	csrfGroup.GET("/endpoint/:endpoint_id/chart", endpointChartGet)
	csrfGroup.GET("/endpoint/:endpoint_id/log", endpointLogGet)

	dbGroup.PUT("/endpoint/:endpoint_id/register",
		handlers.EndpointRegisterPut)
	dbGroup.GET("/endpoint/:endpoint_id/comm",
		handlers.EndpointCommGet)

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
	csrfGroup.DELETE("/service", servicesDelete)
	csrfGroup.DELETE("/service/:service_id", serviceDelete)

	csrfGroup.GET("/session/:user_id", sessionsGet)
	csrfGroup.DELETE("/session/:session_id", sessionDelete)

	csrfGroup.GET("/settings", settingsGet)
	csrfGroup.PUT("/settings", settingsPut)

	csrfGroup.GET("/sshcertificate/:user_id", sshcertsGet)

	csrfGroup.GET("/subscription", subscriptionGet)
	csrfGroup.GET("/subscription/update", subscriptionUpdateGet)
	csrfGroup.POST("/subscription", subscriptionPost)

	csrfGroup.PUT("/theme", themePut)

	csrfGroup.GET("/user", usersGet)
	csrfGroup.GET("/user/:user_id", userGet)
	csrfGroup.PUT("/user/:user_id", userPut)
	csrfGroup.POST("/user", userPost)
	csrfGroup.DELETE("/user", usersDelete)

	engine.GET("/robots.txt", middlewear.RobotsGet)

	if constants.Production {
		sessGroup.GET("/", staticIndexGet)
		engine.GET("/login", staticLoginGet)
		engine.GET("/logo.png", staticLogoGet)
		authGroup.GET("/static/*path", staticGet)
	} else {
		fs := gin.Dir(config.StaticTestingRoot, false)
		fileServer = http.FileServer(fs)

		sessGroup.GET("/", staticTestingGet)
		engine.GET("/login", staticTestingGet)
		engine.GET("/logo.png", staticTestingGet)
		authGroup.GET("/static/*path", staticTestingGet)
	}
}

func init() {
	module := requires.New("mhandlers")
	module.After("settings")

	module.Handler = func() (err error) {
		if constants.Production {
			store, err = static.NewStore(config.StaticRoot)
			if err != nil {
				return
			}
		}

		return
	}
}
