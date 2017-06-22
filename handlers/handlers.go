package handlers

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/static"
	"github.com/pritunl/pritunl-zero/utils"
	"net/http"
)

var (
	store      *static.Store
	fileServer http.Handler
)

// Limit size of request body
func Limiter(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1000000)
}

// Get database from session
func Database(c *gin.Context) {
	db := database.GetDatabase()
	c.Set("db", db)
	c.Next()
	db.Close()
}

// Recover panics
func Recovery(c *gin.Context) {
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

// Register all endpoint handlers
func Register(engine *gin.Engine) {
	engine.Use(Limiter)
	engine.Use(Recovery)

	dbGroup := engine.Group("")
	dbGroup.Use(Database)

	engine.GET("/check", checkGet)

	if constants.Production {
		stre, err := static.NewStore(constants.StaticRoot)
		if err != nil {
			panic(err)
		}
		store = stre

		engine.GET("/", staticIndexGet)
		engine.GET("/static/*path", staticGet)
	} else {
		fs := gin.Dir(constants.StaticTestingRoot, false)
		fileServer = http.FileServer(fs)

		engine.GET("/", staticTestingGet)
		engine.GET("/config.js", staticTestingGet)
		engine.GET("/build.js", staticTestingGet)
		engine.GET("/app/*path", staticTestingGet)
		engine.GET("/dist/*path", staticTestingGet)
		engine.GET("/styles/*path", staticTestingGet)
		engine.GET("/node_modules/*path", staticTestingGet)
		engine.GET("/jspm_packages/*path", staticTestingGet)
	}
}
