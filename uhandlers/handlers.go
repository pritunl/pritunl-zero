package uhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/constants"
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

	dbGroup := engine.Group("")
	dbGroup.Use(middlewear.Database)

	sessGroup := dbGroup.Group("")
	sessGroup.Use(middlewear.SessionUser)

	authGroup := sessGroup.Group("")
	authGroup.Use(middlewear.AuthUser)

	csrfGroup := authGroup.Group("")
	csrfGroup.Use(middlewear.CsrfToken)

	hsmAuthGroup := dbGroup.Group("")
	hsmAuthGroup.Use(middlewear.AuthHsm)

	engine.NoRoute(middlewear.NotFound)

	engine.GET("/auth/state", authStateGet)
	dbGroup.POST("/auth/session", authSessionPost)
	dbGroup.POST("/auth/secondary", authSecondaryPost)
	dbGroup.GET("/auth/request", authRequestGet)
	dbGroup.GET("/auth/callback", authCallbackGet)
	engine.GET("/auth/u2f/app.json", authU2fAppGet)
	dbGroup.GET("/auth/u2f/register", authU2fRegisterGet)
	dbGroup.POST("/auth/u2f/register", authU2fRegisterPost)
	dbGroup.GET("/auth/u2f/sign", authU2fSignGet)
	dbGroup.POST("/auth/u2f/sign", authU2fSignPost)
	sessGroup.GET("/logout", logoutGet)
	sessGroup.GET("/logout_all", logoutAllGet)

	engine.GET("/check", checkGet)

	authGroup.GET("/csrf", csrfGet)

	csrfGroup.GET("/device", devicesGet)
	csrfGroup.PUT("/device/:device_id", devicePut)
	csrfGroup.DELETE("/device/:device_id", deviceDelete)
	csrfGroup.PUT("/device/:device_id/secondary", deviceU2fSecondaryPut)
	csrfGroup.GET("/device/:device_id/sign", deviceU2fSignGet)
	csrfGroup.POST("/device/:device_id/sign", deviceU2fSignPost)
	csrfGroup.GET("/device/:device_id/register", deviceU2fRegisterGet)
	csrfGroup.POST("/device/:device_id/register", deviceU2fRegisterPost)

	hsmAuthGroup.GET("/hsm", hsmGet)

	sessGroup.GET("/ssh", sshGet)
	csrfGroup.PUT("/ssh/validate/:ssh_token", sshValidatePut)
	csrfGroup.DELETE("/ssh/validate/:ssh_token", sshValidateDelete)
	csrfGroup.PUT("/ssh/secondary", sshSecondaryPut)
	csrfGroup.GET("/ssh/u2f/sign", sshU2fSignGet)
	csrfGroup.POST("/ssh/u2f/sign", sshU2fSignPost)
	dbGroup.POST("/ssh/challenge", sshChallengePost)
	dbGroup.PUT("/ssh/challenge", sshChallengePut)
	dbGroup.POST("/ssh/host", sshHostPost)

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
		authGroup.GET("/config.js", staticTestingGet)
		authGroup.GET("/build.js", staticTestingGet)
		authGroup.GET("/uapp/*path", staticTestingGet)
		authGroup.GET("/dist/*path", staticTestingGet)
		authGroup.GET("/styles/*path", staticTestingGet)
		authGroup.GET("/node_modules/*path", staticTestingGet)
		authGroup.GET("/jspm_packages/*path", staticTestingGet)
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
