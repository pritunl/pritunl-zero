package node

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

var Self *Node

type Node struct {
	Id               bson.ObjectId                       `bson:"_id" json:"id"`
	Name             string                              `bson:"name" json:"name"`
	Type             string                              `bson:"type" json:"type"`
	Timestamp        time.Time                           `bson:"timestamp" json:"timestamp"`
	Port             int                                 `bson:"port" json:"port"`
	Protocol         string                              `bson:"protocol" json:"protocol"`
	ManagementDomain string                              `bson:"management_domain" json:"management_domain"`
	Memory           float64                             `bson:"memory" json:"memory"`
	Load1            float64                             `bson:"load1" json:"load1"`
	Load5            float64                             `bson:"load5" json:"load5"`
	Load15           float64                             `bson:"load15" json:"load15"`
	Services         []bson.ObjectId                     `bson:"services" json:"services"`
	DomainServices   map[string]*service.Service         `bson:"-" json:"-"`
	DomainProxies    map[string][]*httputil.ReverseProxy `bson:"-" json:"-"`
}

func (n *Node) loadDomainServices(db *database.Database) (err error) {
	serviceDomains := map[string]*service.Service{}

	services, err := service.GetMulti(db, n.Services)
	if err != nil {
		n.DomainServices = serviceDomains
		return
	}

	for _, srvc := range services {
		for _, domain := range srvc.Domains {
			serviceDomains[domain] = srvc
		}
	}

	n.DomainServices = serviceDomains

	return
}

func (n *Node) initProxy(srvc *service.Service, server *service.Server) (
	proxy *httputil.ReverseProxy) {

	proxy = &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Set("X-Forwarded-For",
				strings.Split(req.RemoteAddr, ":")[0])
			req.Header.Set("X-Forwarded-Proto", n.Protocol)
			req.Header.Set("X-Forwarded-Port", strconv.Itoa(n.Port))

			req.URL.Scheme = server.Protocol
			req.URL.Host = fmt.Sprintf(
				"%s:%d", server.Hostname, server.Port)
		},
	}

	return
}

func (n *Node) initProxies() {
	proxies := map[string][]*httputil.ReverseProxy{}

	for domain, srvc := range n.DomainServices {
		domainProxies := []*httputil.ReverseProxy{}
		for _, server := range srvc.Servers {
			domainProxies = append(domainProxies, n.initProxy(srvc, server))
		}
		proxies[domain] = domainProxies
	}

	n.DomainProxies = proxies

	return
}

func (n *Node) Load(db *database.Database) (err error) {
	err = n.loadDomainServices(db)
	if err != nil {
		return
	}

	n.initProxies()

	return
}

func (n *Node) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if n.Services == nil {
		n.Services = []bson.ObjectId{}
	}

	if n.Protocol != "http" && n.Protocol != "https" {
		errData = &errortypes.ErrorData{
			Error:   "node_protocol_invalid",
			Message: "Invalid node server protocol",
		}
		return
	}

	if n.Port < 1 || n.Port > 65535 {
		errData = &errortypes.ErrorData{
			Error:   "node_port_invalid",
			Message: "Invalid node server port",
		}
		return
	}

	if n.Type != ManagementProxy {
		n.ManagementDomain = ""
	}

	n.Format()

	return
}

func (n *Node) Format() {
	utils.SortObjectIds(n.Services)
}

func (n *Node) Commit(db *database.Database) (err error) {
	coll := db.Nodes()

	err = coll.Commit(n.Id, n)
	if err != nil {
		return
	}

	return
}

func (n *Node) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Nodes()

	err = coll.CommitFields(n.Id, n, fields)
	if err != nil {
		return
	}

	return
}

func (n *Node) update(db *database.Database) (err error) {
	coll := db.Nodes()

	change := mgo.Change{
		Update: &bson.M{
			"$set": &bson.M{
				"timestamp": n.Timestamp,
				"memory":    n.Memory,
				"load1":     n.Load1,
				"load5":     n.Load5,
				"load15":    n.Load15,
			},
		},
		Upsert:    false,
		ReturnNew: true,
	}

	_, err = coll.Find(&bson.M{
		"_id": n.Id,
	}).Apply(change, n)
	if err != nil {
		return
	}

	return
}

func (n *Node) keepalive() {
	db := database.GetDatabase()
	defer db.Close()

	for {
		n.Timestamp = time.Now()

		mem, err := utils.MemoryUsed()
		if err != nil {
			n.Memory = 0

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("node: Failed to get memory")
		} else {
			n.Memory = mem
		}

		load, err := utils.LoadAverage()
		if err != nil {
			n.Load1 = 0
			n.Load5 = 0
			n.Load15 = 0

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("node: Failed to get load")
		} else {
			n.Load1 = load.Load1
			n.Load5 = load.Load5
			n.Load15 = load.Load15
		}

		err = n.update(db)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("node: Failed to update node")
		}

		err = n.Load(db)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("node: Failed to load service domains")
		}

		time.Sleep(1 * time.Second)
	}
}

func (n *Node) Init() (err error) {
	_ = service.Server{}

	db := database.GetDatabase()
	defer db.Close()

	coll := db.Nodes()

	err = coll.FindOneId(n.Id, n)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		default:
			return
		}
	}

	if n.Name == "" {
		n.Name = utils.RandName()
	}

	if n.Type == "" {
		n.Type = Management
	}

	_, err = coll.UpsertId(n.Id, &bson.M{
		"$set": &bson.M{
			"_id":       n.Id,
			"name":      n.Name,
			"type":      n.Type,
			"timestamp": time.Now(),
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	err = n.Load(db)
	if err != nil {
		return
	}

	event.PublishDispatch(db, "node.change")

	Self = n

	go n.keepalive()

	return
}
