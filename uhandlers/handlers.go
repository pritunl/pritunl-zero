package uhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/config"
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

	engine.NoRoute(middlewear.NotFound)

	engine.GET("/auth/state", authStateGet)
	dbGroup.POST("/auth/session", authSessionPost)
	dbGroup.GET("/auth/request", authRequestGet)
	dbGroup.GET("/auth/callback", authCallbackGet)
	sessGroup.GET("/logout", logoutGet)
	sessGroup.GET("/logout_all", logoutAllGet)

	engine.GET("/check", checkGet)

	authGroup.GET("/csrf", csrfGet)

	sessGroup.GET("/ssh", sshGet)
	csrfGroup.PUT("/ssh/validate/:ssh_token", sshValidatePut)
	csrfGroup.DELETE("/ssh/validate/:ssh_token", sshValidateDelete)
	dbGroup.PUT("/ssh/challenge", sshChallengePut)
	dbGroup.POST("/ssh/challenge", sshChallengePost)

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
