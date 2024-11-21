package proxy

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/logger"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/searches"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
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
				node.Self.GetRemoteAddr(req))
			req.Header.Set("X-Forwarded-Host", req.Host)
			req.Header.Set("X-Forwarded-Proto", w.proxyProto)
			req.Header.Set("X-Forwarded-Port", strconv.Itoa(w.proxyPort))

			if authr != nil {
				usr, _ := authr.GetUser(nil)
				if usr != nil {
					req.Header.Set("X-Forwarded-User", usr.Username)
				}
			}

			if w.reqHost != "" {
				req.Host = w.reqHost
			}

			req.URL.Scheme = w.serverProto
			req.URL.Host = w.serverHost

			stripCookieHeaders(req)

			if settings.Elastic.ProxyRequests {
				index := searches.Request{
					Address:   node.Self.GetRemoteAddr(req),
					Timestamp: time.Now(),
					Scheme:    req.URL.Scheme,
					Host:      req.URL.Host,
					Path:      req.URL.Path,
					Query:     req.URL.Query(),
					Header:    req.Header.Clone(),
				}

				if authr.IsValid() {
					usr, _ := authr.GetUser(nil)

					if usr != nil {
						index.User = usr.Id.Hex()
						index.Username = usr.Username
						index.Session = authr.SessionId()
					}
				}

				contentType := strings.ToLower(req.Header.Get("Content-Type"))
				if searches.RequestTypes.Contains(contentType) &&
					req.ContentLength != 0 &&
					req.Body != nil {

					bodyCopy := &bytes.Buffer{}
					tee := io.TeeReader(req.Body, bodyCopy)
					body, _ := ioutil.ReadAll(tee)
					_ = req.Body.Close()
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
	headerTimeout := time.Duration(
		settings.Router.HeaderTimeout) * time.Second

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
	}

	if host.Service.Http2 {
		tlsConfig.NextProtos = []string{"h2"}
	}

	if settings.Router.SkipVerify || net.ParseIP(server.Hostname) != nil {
		tlsConfig.InsecureSkipVerify = true
	}

	if host.ClientCertificate != nil {
		tlsConfig.Certificates = []tls.Certificate{
			*host.ClientCertificate,
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

	var dialer2 *net.Dialer
	if host.Service.Http2 {
		dialer2 = &net.Dialer{
			Timeout:   dialTimeout,
			KeepAlive: dialKeepAlive,
			DualStack: true,
		}
	}

	transportFix := &TransportFix{
		transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   dialTimeout,
				KeepAlive: dialKeepAlive,
				DualStack: true,
			}).DialContext,
			MaxResponseHeaderBytes: int64(
				settings.Router.MaxResponseHeaderBytes),
			MaxIdleConns:          maxIdleConns,
			MaxIdleConnsPerHost:   maxIdleConnsPerHost,
			ResponseHeaderTimeout: headerTimeout,
			IdleConnTimeout:       idleConnTimeout,
			TLSHandshakeTimeout:   handshakeTimeout,
			ExpectContinueTimeout: continueTimeout,
			TLSClientConfig:       tlsConfig,
		},
	}

	if host.Service.Http2 {
		transportFix.transport2 = &http2.Transport{
			TLSClientConfig: tlsConfig,
			DialTLS: func(network, addr string, cfg *tls.Config) (
				conn net.Conn, err error) {

				dialConn, err := dialer2.Dial(network, addr)
				if err != nil {
					err = &errortypes.RequestError{
						errors.Wrap(err, "proxy: Transport dialer error"),
					}
					return
				}

				if server.Protocol == "http" {
					conn = dialConn
				} else {
					conn = tls.Client(dialConn, cfg)
				}
				return
			},
			AllowHTTP: true,
		}
	}

	w = &web{
		reqHost:     host.Domain.Host,
		serverProto: server.Protocol,
		serverHost:  utils.FormatHostPort(server.Hostname, server.Port),
		proxyProto:  proxyProto,
		proxyPort:   proxyPort,
		Transport:   transportFix,
		ErrorLog:    log.New(writer, "", 0),
	}

	return
}
