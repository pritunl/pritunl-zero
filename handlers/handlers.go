package handlers

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/cookie"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/static"
	"github.com/pritunl/pritunl-zero/utils"
	"net/http"
)

var (
	store      *static.Store
	fileServer http.Handler
)

func limiterHand(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1000000)
}

func databaseHand(c *gin.Context) {
	db := database.GetDatabase()
	c.Set("db", db)
	c.Next()
	db.Close()
}

func sessionHand(required bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(*database.Database)

		var sess *session.Session

		cook, err := cookie.Get(c)
		if err == nil {
			sess, err = cook.GetSession(db)
			switch err.(type) {
			case nil:
			case *errortypes.NotFoundError:
				sess = nil
				err = nil
			default:
				c.AbortWithError(500, err)
				return
			}
		}

		if required {
			if sess == nil {
				c.AbortWithStatus(401)
				return
			}

			usr, err := sess.GetUser(db)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}

			if usr.Administrator != "super" {
				sess = nil

				err = cook.Remove(db)
				if err != nil {
					c.AbortWithError(500, err)
					return
				}

				c.AbortWithStatus(401)
				return
			}
		}

		c.Set("session", sess)
	}
}

func activeHand(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	usr, err := sess.GetUser(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	err = usr.SetActive(db)
	if err != nil {
		return
	}

	c.Next()
}

func recoveryHand(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithFields(logrus.Fields{
				"client": utils.GetRemoteAddr(c),
				"error":  errors.New(fmt.Sprintf("%s", r)),
			}).Error("handlers: Handler panic")
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
	}()

	c.Next()
}

func Register(engine *gin.Engine) {
	engine.Use(limiterHand)
	engine.Use(recoveryHand)
	engine.Use(location.New(location.Config{
		Scheme: "https",
	}))

	dbGroup := engine.Group("")
	dbGroup.Use(databaseHand)

	sessGroup := dbGroup.Group("")
	sessGroup.Use(sessionHand(false))

	authGroup := dbGroup.Group("")
	authGroup.Use(sessionHand(true))

	activeAuthGroup := authGroup.Group("")
	activeAuthGroup.Use(activeHand)

	engine.GET("/check", checkGet)

	engine.GET("/auth/state", authStateGet)
	dbGroup.POST("/auth/session", authSessionPost)
	activeAuthGroup.GET("/logout", logoutGet)

	activeAuthGroup.GET("/event", eventGet)

	activeAuthGroup.GET("/settings", settingsGet)
	activeAuthGroup.PUT("/settings", settingsPut)

	activeAuthGroup.GET("/subscription", subscriptionGet)
	activeAuthGroup.GET("/subscription/update", subscriptionUpdateGet)
	activeAuthGroup.POST("/subscription", subscriptionPost)

	activeAuthGroup.GET("/user", usersGet)
	activeAuthGroup.GET("/user/:user_id", userGet)
	activeAuthGroup.PUT("/user/:user_id", userPut)
	activeAuthGroup.POST("/user", userPost)
	activeAuthGroup.DELETE("/user", usersDelete)

	if constants.Production {
		stre, err := static.NewStore(constants.StaticRoot)
		if err != nil {
			panic(err)
		}
		store = stre

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
