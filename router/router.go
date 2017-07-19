package router

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/mhandlers"
	"github.com/pritunl/pritunl-zero/node"
	"net/http"
	"time"
)

type Router struct {
	Node     *node.Node
	typ      string
	port     int
	protocol string
	mRouter  *gin.Engine
	pRouters map[string]*gin.Engine
}

func (r *Router) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	if r.typ == node.Management {
		r.mRouter.ServeHTTP(w, re)
		return
	}

	http.Error(w, "Not found", 404)
}

func (r *Router) Run() (err error) {
	r.typ = r.Node.Type

	r.port = r.Node.Port
	if r.port == 0 {
		r.port = 443
	}

	r.protocol = r.Node.Protocol
	if r.protocol == "" {
		r.protocol = "https"
	}

	if r.typ == node.Management {
		r.mRouter = gin.New()

		if !constants.Production {
			r.mRouter.Use(gin.Logger())
		}

		mhandlers.Register(r.mRouter)
	}

	server := &http.Server{
		Addr:           fmt.Sprintf("0.0.0.0:%d", r.port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 4096,
	}

	logrus.WithFields(logrus.Fields{
		"type":       r.typ,
		"production": constants.Production,
	}).Info("node: Starting node")

	if r.protocol == "http" {
		err = server.ListenAndServe()
		if err != nil {
			err = &errortypes.UnknownError{
				errors.Wrap(err, "node: Server listen failed"),
			}
			return
		}
	}

	return
}

func init() {
	if constants.Production {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}
