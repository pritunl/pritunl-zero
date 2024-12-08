package router

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/acme"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/endpoint"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/mhandlers"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/phandlers"
	"github.com/pritunl/pritunl-zero/proxy"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/uhandlers"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Router struct {
	nodeHash         []byte
	singleType       bool
	managementType   bool
	userType         bool
	proxyType        bool
	port             int
	noRedirectServer bool
	protocol         string
	certificates     *Certificates
	managementDomain string
	userDomain       string
	endpointDomain   string
	stateLock        sync.Mutex
	mRouter          *gin.Engine
	uRouter          *gin.Engine
	pRouter          *gin.Engine
	waiter           sync.WaitGroup
	lock             sync.Mutex
	redirectServer   *http.Server
	webServer        *http.Server
	proxy            *proxy.Proxy
	stop             bool
}

func (r *Router) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	if node.Self.ForwardedProtoHeader != "" &&
		strings.ToLower(re.Header.Get(
			node.Self.ForwardedProtoHeader)) == "http" {

		re.URL.Host = utils.StripPort(re.Host)
		re.URL.Scheme = "https"

		http.Redirect(w, re, re.URL.String(),
			http.StatusMovedPermanently)
		return
	}

	if r.singleType {
		if r.managementType {
			r.mRouter.ServeHTTP(w, re)
		} else if r.userType {
			r.uRouter.ServeHTTP(w, re)
		} else if r.proxyType {
			if !r.proxy.ServeHTTP(w, re) {
				r.pRouter.ServeHTTP(w, re)
			}
		} else {
			utils.WriteStatus(w, 520)
		}
		return
	} else {
		hst := utils.StripPort(re.Host)
		if r.managementType && hst == r.managementDomain {
			r.mRouter.ServeHTTP(w, re)
			return
		} else if r.userType && (hst == r.userDomain ||
			(r.endpointDomain != "" && hst == r.endpointDomain)) {

			r.uRouter.ServeHTTP(w, re)
			return
		} else if r.proxyType {
			if !r.proxy.ServeHTTP(w, re) {
				r.pRouter.ServeHTTP(w, re)
			}
			return
		}
	}

	if re.URL.Path == "/check" {
		utils.WriteText(w, 200, "ok")
		return
	}

	utils.WriteStatus(w, 404)
}

func (r *Router) initRedirect() (err error) {
	r.redirectServer = &http.Server{
		Addr:           ":80",
		ReadTimeout:    1 * time.Minute,
		WriteTimeout:   1 * time.Minute,
		IdleTimeout:    1 * time.Minute,
		MaxHeaderBytes: 8192,
		Handler: http.HandlerFunc(func(
			w http.ResponseWriter, req *http.Request) {

			if strings.HasPrefix(req.URL.Path, acme.AcmePath) {
				token := acme.ParsePath(req.URL.Path)
				token = utils.FilterStr(token, 96)
				if token != "" {
					chal, err := acme.GetChallenge(token)
					if err != nil {
						utils.WriteStatus(w, 400)
					} else {
						logrus.WithFields(logrus.Fields{
							"token": token,
						}).Info("router: Acme challenge requested")
						utils.WriteText(w, 200, chal.Resource)
					}
					return
				}
			} else if req.URL.Path == "/check" {
				utils.WriteText(w, 200, "ok")
				return
			}

			newHost := utils.StripPort(req.Host)
			if r.port != 443 {
				newHost += fmt.Sprintf(":%d", r.port)
			}

			req.URL.Host = newHost
			req.URL.Scheme = "https"

			http.Redirect(w, req, req.URL.String(),
				http.StatusMovedPermanently)
		}),
	}

	return
}

func (r *Router) startRedirect() {
	defer r.waiter.Done()

	if r.port == 80 || r.noRedirectServer {
		return
	}

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
	r.managementType = node.Self.IsManagement()
	r.userType = node.Self.IsUser()
	r.proxyType = node.Self.IsProxy()
	r.managementDomain = node.Self.ManagementDomain
	r.userDomain = node.Self.UserDomain
	r.endpointDomain = node.Self.EndpointDomain
	r.noRedirectServer = node.Self.NoRedirectServer

	if r.managementType && !r.userType && !r.proxyType {
		r.singleType = true
	} else if r.userType && !r.proxyType && !r.managementType {
		r.singleType = true
	} else if r.proxyType && !r.managementType && !r.userType {
		r.singleType = true
	} else {
		r.singleType = false
	}

	r.port = node.Self.Port
	if r.port == 0 {
		r.port = 443
	}

	r.protocol = node.Self.Protocol
	if r.protocol == "" {
		r.protocol = "https"
	}

	if r.managementType {
		r.mRouter = gin.New()

		if constants.DebugWeb {
			r.mRouter.Use(gin.Logger())
		}

		mhandlers.Register(r.mRouter)
	}

	if r.userType {
		r.uRouter = gin.New()

		if constants.DebugWeb {
			r.uRouter.Use(gin.Logger())
		}

		uhandlers.Register(r.uRouter)
	}

	if r.proxyType {
		r.pRouter = gin.New()

		if constants.DebugWeb {
			r.pRouter.Use(gin.Logger())
		}

		phandlers.Register(r.proxy, r.pRouter)
	}

	readTimeout := time.Duration(settings.Router.ReadTimeout) * time.Second
	headerTimeout := time.Duration(
		settings.Router.HeaderTimeout) * time.Second
	writeTimeout := time.Duration(settings.Router.WriteTimeout) * time.Second
	idleTimeout := time.Duration(settings.Router.IdleTimeout) * time.Second

	r.webServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", r.port),
		Handler:           r,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: headerTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    settings.Router.MaxHeaderBytes,
	}

	if settings.Router.H2cSupport {
		h2s := &http2.Server{
			IdleTimeout:     idleTimeout,
			ReadIdleTimeout: readTimeout,
		}

		r.webServer.Handler = h2c.NewHandler(r, h2s)
	}

	return
}

func (r *Router) startWeb() {
	defer r.waiter.Done()

	logrus.WithFields(logrus.Fields{
		"production":          constants.Production,
		"protocol":            r.protocol,
		"port":                r.port,
		"read_timeout":        settings.Router.ReadTimeout,
		"write_timeout":       settings.Router.WriteTimeout,
		"idle_timeout":        settings.Router.IdleTimeout,
		"read_header_timeout": settings.Router.HeaderTimeout,
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
				return
			}
		}
	} else {
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
			NextProtos: []string{"h2"},
			CipherSuites: []uint16{
				tls.TLS_AES_128_GCM_SHA256,                        // 0x1301
				tls.TLS_AES_256_GCM_SHA384,                        // 0x1302
				tls.TLS_CHACHA20_POLY1305_SHA256,                  // 0x1303
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,       // 0xc02b
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,         // 0xc02f
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,       // 0xc02c
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,         // 0xc030
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256, // 0xcca9
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,   // 0xcca8
			},
			GetCertificate: r.certificates.GetCertificate,
		}
		tlsConfig.Certificates = []tls.Certificate{}

		listener, err := tls.Listen("tcp", r.webServer.Addr, tlsConfig)
		if err != nil {
			err = &errortypes.UnknownError{
				errors.Wrap(err, "router: TLS listen failed"),
			}
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("router: Web server TLS error")
			return
		}

		err = r.webServer.Serve(listener)
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
				return
			}
		}
	}

	return
}

func (r *Router) initServers() (err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	err = r.certificates.Init()
	if err != nil {
		return
	}

	err = r.updateState()
	if err != nil {
		return
	}

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
		_ = r.redirectServer.Shutdown(redirectCtx)
	}
	if r.webServer != nil {
		webCtx, webCancel := context.WithTimeout(
			context.Background(),
			1*time.Second,
		)
		defer webCancel()
		_ = r.webServer.Shutdown(webCtx)
	}

	func() {
		defer func() {
			recover()
		}()
		if r.redirectServer != nil {
			_ = r.redirectServer.Close()
		}
		if r.webServer != nil {
			_ = r.webServer.Close()
		}
	}()

	event.WebSocketsStop()
	proxy.WebSocketsStop()
	endpoint.WebSocketsStop()

	r.redirectServer = nil
	r.webServer = nil

	time.Sleep(250 * time.Millisecond)
}

func (r *Router) Shutdown() {
	r.stop = true
	r.Restart()
	time.Sleep(1 * time.Second)
	r.Restart()
	time.Sleep(1 * time.Second)
	r.Restart()
}

func (r *Router) hashNode() []byte {
	hash := md5.New()
	_, _ = io.WriteString(hash, node.Self.Type)
	_, _ = io.WriteString(hash, node.Self.ManagementDomain)
	_, _ = io.WriteString(hash, node.Self.UserDomain)
	_, _ = io.WriteString(hash, node.Self.EndpointDomain)
	_, _ = io.WriteString(hash, strconv.Itoa(node.Self.Port))
	_, _ = io.WriteString(hash, fmt.Sprintf("%t", node.Self.NoRedirectServer))
	_, _ = io.WriteString(hash, node.Self.Protocol)

	_, _ = io.WriteString(hash, strconv.Itoa(settings.Router.ReadTimeout))
	_, _ = io.WriteString(hash, strconv.Itoa(settings.Router.HeaderTimeout))
	_, _ = io.WriteString(hash, strconv.Itoa(settings.Router.WriteTimeout))
	_, _ = io.WriteString(hash, strconv.Itoa(settings.Router.IdleTimeout))
	_, _ = io.WriteString(hash, strconv.Itoa(settings.Router.MaxHeaderBytes))

	certs := node.Self.CertificateObjs
	if certs != nil {
		for _, cert := range certs {
			_, _ = io.WriteString(hash, cert.Hash())
		}
	}

	return hash.Sum(nil)
}

func (r *Router) watchNode() {
	for {
		time.Sleep(1 * time.Second)

		hash := r.hashNode()
		if bytes.Compare(r.nodeHash, hash) != 0 {
			r.nodeHash = hash
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			r.Restart()
			time.Sleep(2 * time.Second)
		}
	}
}

func (r *Router) updateState() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	r.stateLock.Lock()
	defer r.stateLock.Unlock()

	err = r.certificates.Update(db)
	if err != nil {
		return
	}

	return
}

func (r *Router) watchState() {
	for {
		time.Sleep(4 * time.Second)

		err := r.updateState()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("proxy: Failed to load proxy state")
		}
	}
}

func (r *Router) Run() (err error) {
	r.nodeHash = r.hashNode()
	go r.watchNode()
	go r.watchState()

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

		if r.stop {
			break
		}
	}

	return
}

func (r *Router) Init() {
	if constants.DebugWeb {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r.certificates = &Certificates{}
	r.proxy = &proxy.Proxy{}
	r.proxy.Init()
}
