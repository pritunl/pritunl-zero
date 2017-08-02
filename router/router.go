package router

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/acme"
	"github.com/pritunl/pritunl-zero/auth"
	"github.com/pritunl/pritunl-zero/certificate"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/mhandlers"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/phandlers"
	"github.com/pritunl/pritunl-zero/utils"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Router struct {
	nodeHash         []byte
	typ              string
	port             int
	protocol         string
	certificate      *certificate.Certificate
	managementDomain string
	mRouter          *gin.Engine
	pRouter          *gin.Engine
	waiter           sync.WaitGroup
	lock             sync.Mutex
	redirectServer   *http.Server
	webServer        *http.Server
}

func (r *Router) proxy(w http.ResponseWriter, re *http.Request, hst string) {
	host := node.Self.Handler.Hosts[hst]
	proxies, ok := node.Self.Handler.Proxies[hst]
	n := len(proxies)

	if host != nil && ok && n != 0 {
		ndeScheme := ""
		ndePort := ""
		if node.Self.Protocol == "http" {
			ndeScheme = "http"
			if node.Self.Port != 80 {
				ndePort += fmt.Sprintf(":%d", node.Self.Port)
			}
		} else {
			ndeScheme = "https"
			if node.Self.Port != 443 {
				ndePort += fmt.Sprintf(":%d", node.Self.Port)
			}
		}
		domain := fmt.Sprintf(
			"%s://%s%s",
			ndeScheme,
			host.Domain.Domain,
			ndePort,
		)

		origin := re.Header.Get("Origin")
		if origin != "" {
			u, err := url.Parse(origin)
			if err != nil {
				http.Error(w, "Server error", 500)
				return
			}
			origin = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
		}

		referer := re.Header.Get("Referer")
		if referer != "" {
			u, err := url.Parse(referer)
			if err != nil {
				http.Error(w, "Server error", 500)
				return
			}
			referer = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
		}

		if origin != "" && origin != domain {
			http.Error(w, "CSRF origin error", 401)
			return
		}
		if referer != "" && referer != domain {
			http.Error(w, "CSRF referer error", 401)
			return
		}
		if origin == "" && referer == "" {
			http.Error(w, "CSRF origin referer error", 401)
			return
		}

		db := database.GetDatabase()
		defer db.Close()

		cook, sess, err := auth.CookieSessionProxy(db, host.Service, w, re)
		if err != nil {
			http.Error(w, "Server error", 500)
			return
		}

		if sess == nil {
			err = cook.Remove(db)
			if err != nil {
				http.Error(w, "Server error", 500)
				return
			}

			r.pRouter.ServeHTTP(w, re)
			return
		}

		usr, err := sess.GetUser(db)
		if err != nil {
			http.Error(w, "Server error", 500)
			return
		}

		errData, err := auth.Validate(db, usr, host.Service, re)
		if err != nil {
			http.Error(w, "Server error", 500)
			return
		}

		if errData != nil {
			err = cook.Remove(db)
			if err != nil {
				http.Error(w, "Server error", 500)
				return
			}

			r.pRouter.ServeHTTP(w, re)
			return
		}

		proxies[rand.Intn(n)].ServeHTTP(w, re)
		return
	}

	http.Error(w, "Not found", 404)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	hst := utils.StripPort(re.Host)
	if r.typ == node.Management {
		r.mRouter.ServeHTTP(w, re)
		return
	} else if r.typ == node.ManagementProxy && hst == r.managementDomain {
		r.mRouter.ServeHTTP(w, re)
		return
	} else {
		r.proxy(w, re, hst)
		return
	}

	http.Error(w, "Not found", 404)
}

func (r *Router) initRedirect() (err error) {
	r.redirectServer = &http.Server{
		Addr:         ":80",
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		Handler: http.HandlerFunc(func(
			w http.ResponseWriter, req *http.Request) {

			if strings.HasPrefix(req.URL.Path, acme.AcmePath) {
				token := acme.ParsePath(req.URL.Path)
				if token != "" {
					chal, err := acme.GetChallenge(token)
					if err == nil {
						logrus.WithFields(logrus.Fields{
							"token": token,
						}).Info("router: Acme challenge requested")
						io.WriteString(w, chal.Resource)
						return
					}
				}
			}

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

	logrus.WithFields(logrus.Fields{
		"production": constants.Production,
		"protocol":   "http",
		"port":       80,
	}).Info("router: Starting redirect server")

	err := r.redirectServer.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			err = nil
		} else {
			err = &errortypes.UnknownError{
				errors.Wrap(err, "router: Server listen failed"),
			}
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("router: Redirect server error")
		}
	}
}

func (r *Router) initWeb() (err error) {
	r.typ = node.Self.Type
	r.managementDomain = node.Self.ManagementDomain
	r.certificate = node.Self.CertificateObj

	r.port = node.Self.Port
	if r.port == 0 {
		r.port = 443
	}

	r.protocol = node.Self.Protocol
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

	if r.typ == node.Proxy || r.typ == node.ManagementProxy {
		r.pRouter = gin.New()

		if !constants.Production {
			r.pRouter.Use(gin.Logger())
		}

		phandlers.Register(r.protocol, r.pRouter)
	}

	r.webServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", r.port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 4096,
	}

	if r.protocol != "http" {
		if r.certificate != nil {
			err = r.certificate.Write()
			if err != nil {
				return
			}
		} else {
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

	logrus.WithFields(logrus.Fields{
		"production": constants.Production,
		"protocol":   r.protocol,
		"port":       r.port,
	}).Info("router: Starting web server")

	if r.protocol == "http" {
		err := r.webServer.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				err = nil
			} else {
				err = &errortypes.UnknownError{
					errors.Wrap(err, "router: Server listen failed"),
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
					errors.Wrap(err, "router: Server listen failed"),
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

func (r *Router) hashNode() []byte {
	hash := md5.New()
	io.WriteString(hash, node.Self.Type)
	io.WriteString(hash, node.Self.ManagementDomain)
	io.WriteString(hash, strconv.Itoa(node.Self.Port))
	io.WriteString(hash, node.Self.Protocol)

	cert := node.Self.CertificateObj
	if cert != nil {
		io.WriteString(hash, cert.Hash())
	}

	return hash.Sum(nil)
}

func (r *Router) watchNode() {
	for {
		time.Sleep(1 * time.Second)

		hash := r.hashNode()
		if bytes.Compare(r.nodeHash, hash) != 0 {
			r.nodeHash = hash
			r.Restart()
			time.Sleep(2 * time.Second)
		}
	}

	return
}

func (r *Router) Run() (err error) {
	r.nodeHash = r.hashNode()
	go r.watchNode()

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
