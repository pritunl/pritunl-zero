package proxy

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/auth"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/pritunl/pritunl-zero/validator"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"net"
	"net/http"
	"time"
)

type Host struct {
	Service           *service.Service
	Domain            *service.Domain
	WhitelistNetworks []*net.IPNet
}

type Proxy struct {
	Hosts     map[string]*Host
	nodeHash  []byte
	wProxies  map[string][]*web
	wsProxies map[string][]*webSocket
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) bool {
	hst := utils.StripPort(r.Host)

	host := p.Hosts[hst]
	wProxies := p.wProxies[hst]
	wsProxies := p.wsProxies[hst]

	wLen := 0
	if wProxies != nil {
		wLen = len(wProxies)
	}

	wsLen := 0
	if wsProxies != nil {
		wsLen = len(wsProxies)
	}

	if host == nil || wLen == 0 {
		if r.URL.Path == "/check" {
			utils.WriteText(w, 200, "ok")
			return true
		}

		utils.WriteStatus(w, 404)
		return true
	}

	if !host.Service.DisableCsrfCheck {
		valid := auth.CsrfCheck(w, r, host.Domain.Domain)
		if !valid {
			return true
		}
	}

	db := database.GetDatabase()
	defer db.Close()

	clientIp := net.ParseIP(node.Self.GetRemoteAddr(r))
	if clientIp != nil {
		for _, network := range host.WhitelistNetworks {
			if network.Contains(clientIp) {
				if wsProxies != nil &&
					r.Header.Get("Upgrade") == "websocket" {

					wsProxies[rand.Intn(wsLen)].ServeHTTP(
						w, r, db, authorizer.NewProxy())
					return true
				}

				wProxies[rand.Intn(wLen)].ServeHTTP(
					w, r, authorizer.NewProxy())
				return true
			}
		}
	}

	authr, err := authorizer.AuthorizeProxy(db, host.Service, w, r)
	if err != nil {
		WriteError(w, r, 500, err)
		return true
	}

	if !authr.IsValid() {
		err = authr.Clear(db, w, r)
		if err != nil {
			WriteError(w, r, 500, err)
			return true
		}

		return false
	}

	usr, err := authr.GetUser(db)
	if err != nil {
		WriteError(w, r, 500, err)
		return true
	}

	if usr == nil {
		err = authr.Clear(db, w, r)
		if err != nil {
			WriteError(w, r, 500, err)
			return true
		}

		return false
	}

	active, err := auth.SyncUser(db, usr)
	if err != nil {
		WriteError(w, r, 500, err)
		return true
	}

	if !active {
		err = session.RemoveAll(db, usr.Id)
		if err != nil {
			WriteError(w, r, 500, err)
			return true
		}

		err = authr.Clear(db, w, r)
		if err != nil {
			WriteError(w, r, 500, err)
			return true
		}

		return false
	}

	_, _, errData, err := validator.ValidateProxy(
		db, usr, authr.IsApi(), host.Service, r)
	if err != nil {
		WriteError(w, r, 500, err)
		return true
	}

	if errData != nil {
		err = authr.Clear(db, w, r)
		if err != nil {
			WriteError(w, r, 500, err)
			return true
		}

		return false
	}

	if wsProxies != nil && r.Header.Get("Upgrade") == "websocket" {
		wsProxies[rand.Intn(wsLen)].ServeHTTP(w, r, db, authr)
		return true
	}

	if host.Service.LogoutPath != "" && r.URL.Path == host.Service.LogoutPath {
		err = authr.Clear(db, w, r)
		if err != nil {
			WriteError(w, r, 500, err)
			return true
		}

		return false
	}

	wProxies[rand.Intn(wLen)].ServeHTTP(w, r, authr)
	return true
}

func (p *Proxy) reloadHosts(db *database.Database, services []bson.ObjectId) (
	err error) {

	hosts := map[string]*Host{}
	appId := ""
	facets := []string{}

	if node.Self.UserDomain != "" {
		appId = fmt.Sprintf("https://%s/auth/u2f/app.json",
			node.Self.UserDomain)
	}

	nodeServices := set.NewSet()
	for _, srvc := range services {
		nodeServices.Add(srvc)
	}

	nodes, err := node.GetAll(db)
	if err != nil {
		return
	}

	for _, nde := range nodes {
		if appId == "" {
			appId = fmt.Sprintf("https://%s/auth/u2f/app.json",
				nde.UserDomain)
		}
		if nde.UserDomain != "" {
			facets = append(facets,
				fmt.Sprintf("https://%s", nde.UserDomain))
		}
		if nde.ManagementDomain != "" {
			facets = append(facets,
				fmt.Sprintf("https://%s", nde.ManagementDomain))
		}
	}

	srvcs, err := service.GetAll(db)
	if err != nil {
		p.Hosts = hosts
		return
	}

	for _, srvc := range srvcs {
		nodeService := nodeServices.Contains(srvc.Id)

		for _, domain := range srvc.Domains {
			facets = append(facets, fmt.Sprintf("https://%s", domain.Domain))

			if !nodeService {
				continue
			}
			whitelistNets := []*net.IPNet{}

			for _, cidr := range srvc.WhitelistNetworks {
				_, network, err := net.ParseCIDR(cidr)
				if err != nil {
					err = &errortypes.ParseError{
						errors.Wrap(err, "proxy: Failed to parse network"),
					}

					logrus.WithFields(logrus.Fields{
						"network": cidr,
						"error":   err,
					}).Error("proxy: Invalid whitelist network")
					err = nil
					continue
				}

				whitelistNets = append(whitelistNets, network)
			}

			srvcDomain := &Host{
				Service:           srvc,
				Domain:            domain,
				WhitelistNetworks: whitelistNets,
			}

			hosts[domain.Domain] = srvcDomain
		}
	}

	settings.Local.AppId = appId
	settings.Local.Facets = facets

	p.Hosts = hosts

	return
}

func (p *Proxy) reloadProxies(db *database.Database, proto string, port int) (
	err error) {

	wProxies := map[string][]*web{}
	wsProxies := map[string][]*webSocket{}

	for domain, host := range p.Hosts {
		domainProxies := []*web{}
		for _, server := range host.Service.Servers {
			prxy := newWeb(proto, port, host, server)
			domainProxies = append(domainProxies, prxy)
		}
		wProxies[domain] = domainProxies

		if host.Service.WebSockets {
			domainWsProxies := []*webSocket{}
			for _, server := range host.Service.Servers {
				prxy := newWebSocket(proto, port, host, server)
				domainWsProxies = append(domainWsProxies, prxy)
			}
			wsProxies[domain] = domainWsProxies
		}
	}

	p.wProxies = wProxies
	p.wsProxies = wsProxies

	return
}

func (p *Proxy) update() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	proto := node.Self.Protocol
	port := node.Self.Port
	services := node.Self.Services

	err = p.reloadHosts(db, services)
	if err != nil {
		return
	}

	err = p.reloadProxies(db, proto, port)
	if err != nil {
		return
	}

	return
}

func (p *Proxy) watchNode() {
	for {
		err := p.update()
		if err != nil {
			p.nodeHash = []byte{}
			p.Hosts = map[string]*Host{}
			p.wProxies = map[string][]*web{}
			p.wsProxies = map[string][]*webSocket{}

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("proxy: Failed to load proxy state")

			return
		}

		time.Sleep(3 * time.Second)
	}

	return
}

func (p *Proxy) Init() {
	p.Hosts = map[string]*Host{}
	p.wProxies = map[string][]*web{}
	p.wsProxies = map[string][]*webSocket{}
	go p.watchNode()
}
