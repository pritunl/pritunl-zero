package mhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/middlewear"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/static"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
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

	engine.GET("/auth/state", authStateGet)
	dbGroup.POST("/auth/session", authSessionPost)
	dbGroup.GET("/auth/request", authRequestGet)
	dbGroup.GET("/auth/callback", authCallbackGet)
	sessGroup.GET("/logout", logoutGet)

	csrfGroup.GET("/authority", authoritysGet)
	csrfGroup.GET("/authority/:authr_id", authorityGet)
	csrfGroup.PUT("/authority/:authr_id", authorityPut)
	csrfGroup.POST("/authority", authorityPost)
	csrfGroup.DELETE("/authority/:authr_id", authorityDelete)
	dbGroup.GET("/ssh_public_key/:authr_ids", authorityPublicKeyGet)

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

		pushFiles = []string{}
		walk := path.Join(config.StaticTestingRoot, "app")
		err := filepath.Walk(walk, func(
			pth string, _ os.FileInfo, e error) (err error) {

			if e != nil {
				err = e
				return
			}

			if strings.HasSuffix(pth, ".js") ||
				strings.HasSuffix(pth, ".js.map") {

				pth = strings.Replace(pth, walk, "/app", 1)
				pushFiles = append(pushFiles, pth)
			}

			return
		})
		if err != nil {
			panic(err)
		}

		sessGroup.GET("/", staticTestingGet)
		engine.GET("/login", staticTestingGet)
		engine.GET("/logo.png", staticTestingGet)
		authGroup.GET("/config.js", staticTestingGet)
		authGroup.GET("/build.js", staticTestingGet)
		authGroup.GET("/app/*path", staticTestingGet)
		authGroup.GET("/dist/*path", staticTestingGet)
		authGroup.GET("/styles/*path", staticTestingGet)
		authGroup.GET("/node_modules/*path", staticTestingGet)
		authGroup.GET("/jspm_packages/*path", staticTestingGet)
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
