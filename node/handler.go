package node

import (
	"fmt"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/service"
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

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	proxy = &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Set("X-Forwarded-For",
				strings.Split(req.RemoteAddr, ":")[0])
			req.Header.Set("X-Forwarded-Proto", h.Node.Protocol)
			req.Header.Set("X-Forwarded-Port", strconv.Itoa(h.Node.Port))

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
