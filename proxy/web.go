package proxy

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/logger"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/search"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

type web struct {
	reqHost     string
	serverHost  string
	serverProto string
	proxyProto  string
	proxyPort   int
	Transport   http.RoundTripper
	ErrorLog    *log.Logger
}

func (w *web) ServeHTTP(rw http.ResponseWriter, r *http.Request,
	authr *authorizer.Authorizer) {

	prxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Set("X-Forwarded-For",
				strings.Split(req.RemoteAddr, ":")[0])
			req.Header.Set("X-Forwarded-Proto", w.proxyProto)
			req.Header.Set("X-Forwarded-Port", strconv.Itoa(w.proxyPort))

			if w.reqHost != "" {
				req.Host = w.reqHost
			}

			req.URL.Scheme = w.serverProto
			req.URL.Host = w.serverHost

			stripCookie(req)

			if settings.Elastic.ProxyRequests {
				index := search.Request{
					Address:   node.Self.GetRemoteAddr(req),
					Timestamp: time.Now(),
					Scheme:    req.URL.Scheme,
					Host:      req.URL.Host,
					Path:      req.URL.Path,
					Query:     req.URL.Query(),
					Header:    req.Header,
				}

				if authr.IsValid() {
					usr, _ := authr.GetUser(nil)

					if usr != nil {
						index.User = usr.Id.Hex()
						index.Session = authr.SessionId()
					}
				}

				contentType := strings.ToLower(req.Header.Get("Content-Type"))
				if search.RequestTypes.Contains(contentType) &&
					req.ContentLength != 0 &&
					req.Body != nil {

					bodyCopy := &bytes.Buffer{}
					tee := io.TeeReader(req.Body, bodyCopy)
					body, _ := ioutil.ReadAll(tee)
					req.Body.Close()
					req.Body = utils.NopCloser{bodyCopy}
					index.Body = string(body)
				}

				index.Index()
			}
		},
		Transport: w.Transport,
		ErrorLog:  w.ErrorLog,
	}

	prxy.ServeHTTP(rw, r)
}

func newWeb(proxyProto string, proxyPort int, host *Host,
	server *service.Server) (w *web) {

	dialTimeout := time.Duration(
		settings.Router.DialTimeout) * time.Second
	dialKeepAlive := time.Duration(
		settings.Router.DialKeepAlive) * time.Second
	maxIdleConns := settings.Router.MaxIdleConns
	maxIdleConnsPerHost := settings.Router.MaxIdleConnsPerHost
	idleConnTimeout := time.Duration(
		settings.Router.IdleConnTimeout) * time.Second
	handshakeTimeout := time.Duration(
		settings.Router.HandshakeTimeout) * time.Second
	continueTimeout := time.Duration(
		settings.Router.ContinueTimeout) * time.Second

	var tlsConfig *tls.Config
	if settings.Router.SkipVerify || net.ParseIP(server.Hostname) != nil {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	writer := &logger.ErrorWriter{
		Message: "node: Proxy server error",
		Fields: logrus.Fields{
			"service": host.Service.Name,
			"domain":  host.Domain.Domain,
			"server": fmt.Sprintf(
				"%s://%s:%d",
				server.Protocol,
				server.Hostname,
				server.Port,
			),
		},
		Filters: []string{
			"context canceled",
		},
	}

	w = &web{
		reqHost:     host.Domain.Host,
		serverProto: server.Protocol,
		serverHost:  fmt.Sprintf("%s:%d", server.Hostname, server.Port),
		proxyProto:  proxyProto,
		proxyPort:   proxyPort,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   dialTimeout,
				KeepAlive: dialKeepAlive,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          maxIdleConns,
			MaxIdleConnsPerHost:   maxIdleConnsPerHost,
			IdleConnTimeout:       idleConnTimeout,
			TLSHandshakeTimeout:   handshakeTimeout,
			ExpectContinueTimeout: continueTimeout,
			TLSClientConfig:       tlsConfig,
		},
		ErrorLog: log.New(writer, "", 0),
	}

	return
}
