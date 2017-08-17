package proxy

import (
	"github.com/Sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/auth"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/utils"
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
		utils.WriteStatus(w, 404)
		return true
	}

	valid := auth.CsrfCheck(w, r, host.Domain.Domain)
	if !valid {
		return true
	}

	clientIp := net.ParseIP(node.Self.GetRemoteAddr(r))
	if clientIp != nil {
		for _, network := range host.WhitelistNetworks {
			if network.Contains(clientIp) {
				if wsProxies != nil && r.Header.Get("Upgrade") == "websocket" {
					wsProxies[rand.Intn(wsLen)].ServeHTTP(w, r, nil)
					return true
				}

				wProxies[rand.Intn(wLen)].ServeHTTP(w, r, nil)
				return true
			}
		}
	}

	db := database.GetDatabase()
	defer db.Close()

	cook, sess, err := auth.CookieSessionProxy(db, host.Service, w, r)
	if err != nil {
		WriteError(w, r, 500, err)
		return true
	}

	if sess == nil {
		if cook != nil {
			err = cook.Remove(db)
			if err != nil {
				WriteError(w, r, 500, err)
				return true
			}
		}

		return false
	}

	usr, err := sess.GetUser(db)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			if cook != nil {
				err = cook.Remove(db)
				if err != nil {
					WriteError(w, r, 500, err)
					return true
				}
			}

			return false
		}

		WriteError(w, r, 500, err)
		return true
	}

	errData, err := auth.Validate(db, usr, host.Service, r)
	if err != nil {
		WriteError(w, r, 500, err)
		return true
	}

	if errData != nil {
		err = cook.Remove(db)
		if err != nil {
			WriteError(w, r, 500, err)
			return true
		}

		return false
	}

	if wsProxies != nil && r.Header.Get("Upgrade") == "websocket" {
		wsProxies[rand.Intn(wsLen)].ServeHTTP(w, r, sess)
		return true
	}

	wProxies[rand.Intn(wLen)].ServeHTTP(w, r, sess)
	return true
}

func (p *Proxy) reloadHosts(db *database.Database, services []bson.ObjectId) (
	err error) {

	hosts := map[string]*Host{}

	srvcs, err := service.GetMulti(db, services)
	if err != nil {
		p.Hosts = hosts
		return
	}

	for _, srvc := range srvcs {
		for _, domain := range srvc.Domains {
			whitelistNets := []*net.IPNet{}

			for _, cidr := range srvc.WhitelistNetworks {
				_, network, err := net.ParseCIDR(cidr)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"network": cidr,
						"error":   err,
					}).Error("proxy: Invalid whitelist network")
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
