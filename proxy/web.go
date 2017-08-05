package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/logger"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/settings"
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

func (w *web) Director(req *http.Request) {
	req.Header.Set("X-Forwarded-For",
		strings.Split(req.RemoteAddr, ":")[0])
	req.Header.Set("X-Forwarded-Proto", w.proxyProto)
	req.Header.Set("X-Forwarded-Port", strconv.Itoa(w.proxyPort))

	if w.reqHost != "" {
		req.Host = w.reqHost
	}

	req.URL.Scheme = w.serverProto
	req.URL.Host = w.serverHost

	cookie := req.Header.Get("Cookie")
	start := strings.Index(cookie, "pritunl-zero=")
	if start != -1 {
		str := cookie[start:]
		end := strings.Index(str, ";")
		if end != -1 {
			if len(str) > end+1 && string(str[end+1]) == " " {
				end += 1
			}
			cookie = cookie[:start] + cookie[start+end+1:]
		} else {
			cookie = cookie[:start]
		}
	}

	cookie = strings.TrimSpace(cookie)

	if len(cookie) > 0 {
		req.Header.Set("Cookie", cookie)
	} else {
		req.Header.Del("Cookie")
	}
}

func (w *web) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	prxy := &httputil.ReverseProxy{
		Director:  w.Director,
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
