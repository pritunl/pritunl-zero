package router

import (
	"context"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/mhandlers"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/utils"
	"net/http"
	"sync"
	"time"
)

type Router struct {
	Node             *node.Node
	typ              string
	port             int
	protocol         string
	managementDomain string
	mRouter          *gin.Engine
	pRouters         map[string]*gin.Engine
	waiter           sync.WaitGroup
	lock             sync.Mutex
	redirectServer   *http.Server
	webServer        *http.Server
}

func (r *Router) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	if r.typ == node.Management {
		r.mRouter.ServeHTTP(w, re)
		return
	} else if r.typ == node.ManagementProxy && re.Host == r.managementDomain {
		r.mRouter.ServeHTTP(w, re)
		return
	}

	http.Error(w, "Not found", 404)
}

func (r *Router) initRedirect() (err error) {
	r.redirectServer = &http.Server{
		Addr:         "0.0.0.0:80",
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		Handler: http.HandlerFunc(func(
			w http.ResponseWriter, req *http.Request) {

			req.URL.Host = req.Host
			req.URL.Scheme = "https"

			http.Redirect(w, req, req.URL.String(),
				http.StatusMovedPermanently)
		}),
	}

	return
}

func (r *Router) startRedirect() {
	defer r.waiter.Done()

	err := r.redirectServer.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			err = nil
		} else {
			err = &errortypes.UnknownError{
				errors.Wrap(err, "node: Server listen failed"),
			}
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("router: Redirect server error")
		}
	}
}

func (r *Router) initWeb() (err error) {
	r.typ = r.Node.Type
	r.managementDomain = r.Node.ManagementDomain

	r.port = r.Node.Port
	if r.port == 0 {
		r.port = 443
	}

	r.protocol = r.Node.Protocol
	if r.protocol == "" {
		r.protocol = "https"
	}

	if r.typ == node.Management || r.typ == node.ManagementProxy {
		r.mRouter = gin.New()

		if !constants.Production {
			r.mRouter.Use(gin.Logger())
		}

		mhandlers.Register(r.protocol, r.mRouter)
	}

	r.webServer = &http.Server{
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

	if r.protocol != "http" {
		certExists, e := utils.Exists(constants.CertPath)
		if e != nil {
			err = e
			return
		}

		keyExists, e := utils.Exists(constants.KeyPath)
		if e != nil {
			err = e
			return
		}

		if !certExists || !keyExists {
			err = generateCert(constants.CertPath, constants.KeyPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func (r *Router) startWeb() {
	defer r.waiter.Done()

	if r.protocol == "http" {
		err := r.webServer.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				err = nil
			} else {
				err = &errortypes.UnknownError{
					errors.Wrap(err, "node: Server listen failed"),
				}
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("router: Web server error")
			}
		}
	} else {
		err := r.webServer.ListenAndServeTLS(
			constants.CertPath, constants.KeyPath)
		if err != nil {
			if err == http.ErrServerClosed {
				err = nil
			} else {
				err = &errortypes.UnknownError{
					errors.Wrap(err, "node: Server listen failed"),
				}
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("router: Web server error")
			}
		}
	}

	return
}

func (r *Router) initServers() (err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	err = r.initRedirect()
	if err != nil {
		return
	}

	err = r.initWeb()
	if err != nil {
		return
	}

	return
}

func (r *Router) startServers() {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.redirectServer == nil || r.webServer == nil {
		return
	}

	r.waiter.Add(2)
	go r.startRedirect()
	go r.startWeb()

	time.Sleep(250 * time.Millisecond)

	return
}

func (r *Router) Restart() {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.redirectServer != nil {
		redirectCtx, redirectCancel := context.WithTimeout(
			context.Background(),
			1*time.Second,
		)
		defer redirectCancel()
		r.redirectServer.Shutdown(redirectCtx)
	}
	if r.webServer != nil {
		webCtx, webCancel := context.WithTimeout(
			context.Background(),
			1*time.Second,
		)
		defer webCancel()
		r.webServer.Shutdown(webCtx)
	}

	func() {
		defer func() {
			recover()
		}()
		if r.redirectServer != nil {
			r.redirectServer.Close()
		}
		if r.webServer != nil {
			r.webServer.Close()
		}
	}()

	event.WebSocketsStop()

	r.redirectServer = nil
	r.webServer = nil

	time.Sleep(250 * time.Millisecond)
}

func (r *Router) Run() (err error) {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			r.Restart()
		}
	}()

	for {
		err = r.initServers()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("router: Failed to init web servers")
			time.Sleep(1 * time.Second)
			continue
		}

		r.waiter = sync.WaitGroup{}
		r.startServers()
		r.waiter.Wait()
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
