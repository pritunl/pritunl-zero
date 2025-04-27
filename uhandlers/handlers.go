package uhandlers

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
)

func Register(engine *gin.Engine) {
	engine.Use(middlewear.Limiter)
	engine.Use(middlewear.Counter)
	engine.Use(middlewear.Recovery)
	engine.Use(middlewear.Headers)

	dbGroup := engine.Group("")
	dbGroup.Use(middlewear.Database)

	authrGroup := dbGroup.Group("")
	authrGroup.Use(middlewear.HasAuthority)

	sessGroup := dbGroup.Group("")
	sessGroup.Use(middlewear.SessionUser)

	authGroup := sessGroup.Group("")
	authGroup.Use(middlewear.AuthUser)

	csrfGroup := authGroup.Group("")
	csrfGroup.Use(middlewear.CsrfToken)

	// hsmAuthGroup := dbGroup.Group("")
	// hsmAuthGroup.Use(middlewear.AuthHsm)

	engine.NoRoute(middlewear.NotFound)

	engine.GET("/auth/state", authStateGet)
	dbGroup.POST("/auth/session", authSessionPost)
	dbGroup.POST("/auth/secondary", authSecondaryPost)
	dbGroup.GET("/auth/request", authRequestGet)
	dbGroup.GET("/auth/callback", authCallbackGet)
	engine.GET("/auth/u2f/app.json", authU2fAppGet)
	dbGroup.GET("/auth/webauthn/request", authWanRequestGet)
	dbGroup.POST("/auth/webauthn/respond", authWanRespondPost)
	dbGroup.GET("/auth/webauthn/register", authWanRegisterGet)
	dbGroup.POST("/auth/webauthn/register", authWanRegisterPost)
	sessGroup.GET("/logout", logoutGet)
	sessGroup.GET("/logout_all", logoutAllGet)

	engine.GET("/check", checkGet)

	authGroup.GET("/csrf", csrfGet)

	csrfGroup.GET("/device", devicesGet)
	csrfGroup.PUT("/device/:device_id", devicePut)
	csrfGroup.DELETE("/device/:device_id", deviceDelete)
	csrfGroup.PUT("/device/:device_id/secondary", deviceSecondaryPut)
	csrfGroup.GET("/device/:device_id/request", deviceWanRequestGet)
	csrfGroup.POST("/device/:device_id/respond", deviceWanRespondPost)
	csrfGroup.GET("/device/:device_id/register", deviceWanRegisterGet)
	csrfGroup.POST("/device/:device_id/register", deviceWanRegisterPost)

	dbGroup.PUT("/endpoint/:endpoint_id/register",
		handlers.EndpointRegisterPut)
	dbGroup.GET("/endpoint/:endpoint_id/comm",
		handlers.EndpointCommGet)

	// hsmAuthGroup.GET("/hsm", hsmGet)

	sessGroup.GET("/ssh", sshGet)
	csrfGroup.PUT("/ssh/validate/:ssh_token", sshValidatePut)
	csrfGroup.DELETE("/ssh/validate/:ssh_token", sshValidateDelete)
	csrfGroup.PUT("/ssh/secondary", sshSecondaryPut)
	csrfGroup.GET("/ssh/webauthn/request", sshWanRequestGet)
	csrfGroup.POST("/ssh/webauthn/respond", sshWanRespondPost)
	authrGroup.PUT("/ssh/challenge", sshChallengePut)
	authrGroup.POST("/ssh/challenge", sshChallengePost)
	authrGroup.POST("/ssh/host", sshHostPost)

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
	module := requires.New("uhandlers")
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
