package middlewear

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/auth"
	"github.com/pritunl/pritunl-zero/cookie"
	"github.com/pritunl/pritunl-zero/csrf"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/utils"
	"net/http"
)

func Limiter(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1000000)
}

func Database(c *gin.Context) {
	db := database.GetDatabase()
	c.Set("db", db)
	c.Next()
	db.Close()
}

func Service(c *gin.Context) {
	srvc := node.Self.Handler.Hosts[c.Request.Host].Service
	c.Set("service", srvc)
}

func Session(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	cook, sess, err := auth.CookieSession(db, c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Set("session", sess)
	c.Set("cookie", cook)
}

func SessionProxy(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	cook, sess, err := auth.CookieSessionProxy(db, c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Set("session", sess)
	c.Set("cookie", cook)
}

func Auth(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)
	cook := c.MustGet("cookie").(*cookie.Cookie)

	if sess == nil {
		c.AbortWithStatus(401)
		return
	}

	usr, err := sess.GetUser(db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if usr.Disabled || usr.Administrator != "super" {
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

func AuthProxy(c *gin.Context) {
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

func CsrfToken(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	sess := c.MustGet("session").(*session.Session)

	if sess == nil {
		c.AbortWithStatus(401)
		return
	}

	token := ""
	if c.Request.Header.Get("Upgrade") == "websocket" {
		token = c.Query("csrf_token")
	} else {
		token = c.Request.Header.Get("Csrf-Token")
	}

	valid, err := csrf.ValidateToken(db, sess.Id, token)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if !valid {
		c.AbortWithStatus(401)
		return
	}
}

func UserActive(c *gin.Context) {
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

func Recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithFields(logrus.Fields{
				"client": utils.GetRemoteAddr(c),
				"error":  errors.New(fmt.Sprintf("%s", r)),
			}).Error("middlewear: Handler panic")
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
	}()

	c.Next()
}
