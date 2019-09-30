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
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/logger"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/search"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
)

type webIsolated struct {
	reqHost     string
	serverHost  string
	serverProto string
	proxyProto  string
	proxyPort   int
	Client      *http.Client
	ErrorLog    *log.Logger
}

func (w *webIsolated) ServeHTTP(rw http.ResponseWriter, r *http.Request,
	authr *authorizer.Authorizer) {

	reqUrl := utils.ProxyUrl(r.URL, w.serverProto, w.serverHost)

	srcBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errortypes.ReadError{
			errors.Wrap(err, "request: Read request failed"),
		}
		WriteError(rw, r, 500, err)
		return
	}

	reqBody := bytes.NewBuffer(srcBody)
	req, err := http.NewRequest(r.Method, reqUrl.String(), reqBody)
	if err != nil {
		err = errortypes.RequestError{
			errors.Wrap(err, "request: Create request failed"),
		}
		WriteError(rw, r, 500, err)
		return
	}

	utils.CopyHeaders(req.Header, r.Header)
	req.Header.Set("X-Forwarded-For",
		node.Self.GetRemoteAddr(req))
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

	stripCookieHeaders(req)

	if settings.Elastic.ProxyRequests {
		index := search.Request{
			Address:   node.Self.GetRemoteAddr(r),
			Timestamp: time.Now(),
			Scheme:    reqUrl.Scheme,
			Host:      reqUrl.Host,
			Path:      reqUrl.Path,
			Query:     reqUrl.Query(),
			Header:    r.Header,
		}

		if authr.IsValid() {
			usr, _ := authr.GetUser(nil)

			if usr != nil {
				index.User = usr.Id.Hex()
				index.Username = usr.Username
				index.Session = authr.SessionId()
			}
		}

		contentType := strings.ToLower(r.Header.Get("Content-Type"))
		if search.RequestTypes.Contains(contentType) &&
			req.ContentLength != 0 && srcBody != nil {

			index.Body = string(srcBody)
		}

		index.Index()
	}

	resp, err := w.Client.Do(req)
	if err != nil {
		err = errortypes.RequestError{
			errors.Wrap(err, "request: Request failed"),
		}
		WriteError(rw, r, 500, err)
		return
	}
	defer resp.Body.Close()

	utils.CopyHeaders(rw.Header(), resp.Header)
	rw.WriteHeader(resp.StatusCode)
	io.Copy(rw, resp.Body)
}

func newWebIsolated(proxyProto string, proxyPort int, host *Host,
	server *service.Server) (w *webIsolated) {

	requestTimeout := time.Duration(
		settings.Router.RequestTimeout) * time.Second
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

	w = &webIsolated{
		reqHost:     host.Domain.Host,
		serverProto: server.Protocol,
		serverHost:  utils.FormatHostPort(server.Hostname, server.Port),
		proxyProto:  proxyProto,
		proxyPort:   proxyPort,
		Client: &http.Client{
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
			CheckRedirect: func(r *http.Request, v []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Timeout: requestTimeout,
		},
		ErrorLog: log.New(writer, "", 0),
	}

	return
}
