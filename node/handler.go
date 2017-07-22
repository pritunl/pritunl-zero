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

var (
	Transport http.RoundTripper = &http.Transport{
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
)

type Handler struct {
	Node     *Node
	Services map[string]*service.Service
	Proxies  map[string][]*httputil.ReverseProxy
}

func (h *Handler) loadServices(db *database.Database) (err error) {
	serviceDomains := map[string]*service.Service{}

	services, err := service.GetMulti(db, h.Node.Services)
	if err != nil {
		h.Services = serviceDomains
		return
	}

	for _, srvc := range services {
		for _, domain := range srvc.Domains {
			serviceDomains[domain] = srvc
		}
	}

	h.Services = serviceDomains

	return
}

func (h *Handler) initProxy(srvc *service.Service, server *service.Server) (
	proxy *httputil.ReverseProxy) {

	proxy = &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Set("X-Forwarded-For",
				strings.Split(req.RemoteAddr, ":")[0])
			req.Header.Set("X-Forwarded-Proto", h.Node.Protocol)
			req.Header.Set("X-Forwarded-Port", strconv.Itoa(h.Node.Port))

			req.URL.Scheme = server.Protocol
			req.URL.Host = fmt.Sprintf(
				"%s:%d", server.Hostname, server.Port)
		},
		Transport: Transport,
	}

	return
}

func (h *Handler) initProxies() {
	proxies := map[string][]*httputil.ReverseProxy{}

	for domain, srvc := range h.Services {
		domainProxies := []*httputil.ReverseProxy{}
		for _, server := range srvc.Servers {
			domainProxies = append(domainProxies, h.initProxy(srvc, server))
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
