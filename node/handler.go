package node

import (
	"crypto/tls"
	"fmt"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/settings"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

type Host struct {
	Service *service.Service
	Domain  *service.Domain
}

type Handler struct {
	Node    *Node
	Hosts   map[string]*Host
	Proxies map[string][]*httputil.ReverseProxy
}

func (h *Handler) loadServices(db *database.Database) (err error) {
	hosts := map[string]*Host{}

	services, err := service.GetMulti(db, h.Node.Services)
	if err != nil {
		h.Hosts = hosts
		return
	}

	for _, srvc := range services {
		for _, domain := range srvc.Domains {
			srvcDomain := &Host{
				Service: srvc,
				Domain:  domain,
			}

			hosts[domain.Domain] = srvcDomain
		}
	}

	h.Hosts = hosts

	return
}

func (h *Handler) initProxy(host *Host, server *service.Server) (
	proxy *httputil.ReverseProxy) {

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

	transport := &http.Transport{
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
	}

	proxy = &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Set("X-Forwarded-For",
				strings.Split(req.RemoteAddr, ":")[0])
			req.Header.Set("X-Forwarded-Proto", h.Node.Protocol)
			req.Header.Set("X-Forwarded-Port", strconv.Itoa(h.Node.Port))

			if host.Domain.Host != "" {
				req.Host = host.Domain.Host
			}

			req.URL.Scheme = server.Protocol
			req.URL.Host = fmt.Sprintf(
				"%s:%d", server.Hostname, server.Port)

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
		},
		Transport: transport,
	}

	return
}

func (h *Handler) initProxies() {
	proxies := map[string][]*httputil.ReverseProxy{}

	for domain, host := range h.Hosts {
		domainProxies := []*httputil.ReverseProxy{}
		for _, server := range host.Service.Servers {
			domainProxies = append(
				domainProxies,
				h.initProxy(host, server),
			)
		}
		proxies[domain] = domainProxies
	}

	h.Proxies = proxies

	return
}

func (h *Handler) Load(db *database.Database) (err error) {
	err = h.loadServices(db)
	if err != nil {
		return
	}

	h.initProxies()

	return
}
