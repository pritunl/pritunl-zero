package phandlers

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
	store *static.Store
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

func sessionHand(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	var sess *session.Session

	cook, err := cookie.Get(c.Writer, c.Request)
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

	c.Set("session", sess)
	c.Set("cookie", cook)
}

func authHand(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	if sess == nil {
		c.AbortWithStatus(401)
		return
	}

	_, err := sess.GetUser(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
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

func Register(protocol string, engine *gin.Engine) {
	engine.Use(limiterHand)
	engine.Use(recoveryHand)
	engine.Use(location.New(location.Config{
		Scheme: protocol,
	}))

	dbGroup := engine.Group("")
	dbGroup.Use(databaseHand)

	sessGroup := dbGroup.Group("")
	sessGroup.Use(sessionHand)

	authGroup := sessGroup.Group("")
	authGroup.Use(authHand)

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
