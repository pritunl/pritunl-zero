package node

import (
	"container/list"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/certificate"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	Self *Node
)

type Node struct {
	Id                   bson.ObjectId              `bson:"_id" json:"id"`
	Name                 string                     `bson:"name" json:"name"`
	Type                 string                     `bson:"type" json:"type"`
	Timestamp            time.Time                  `bson:"timestamp" json:"timestamp"`
	Port                 int                        `bson:"port" json:"port"`
	Protocol             string                     `bson:"protocol" json:"protocol"`
	Certificate          bson.ObjectId              `bson:"certificate" json:"certificate"`
	Certificates         []bson.ObjectId            `bson:"certificates" json:"certificates"`
	ManagementDomain     string                     `bson:"management_domain" json:"management_domain"`
	UserDomain           string                     `bson:"user_domain" json:"user_domain"`
	Services             []bson.ObjectId            `bson:"services" json:"services"`
	Authorities          []bson.ObjectId            `bson:"authorities" json:"authorities"`
	RequestsMin          int64                      `bson:"requests_min" json:"requests_min"`
	ForwardedForHeader   string                     `bson:"forwarded_for_header" json:"forwarded_for_header"`
	ForwardedProtoHeader string                     `bson:"forwarded_proto_header" json:"forwarded_proto_header"`
	Memory               float64                    `bson:"memory" json:"memory"`
	Load1                float64                    `bson:"load1" json:"load1"`
	Load5                float64                    `bson:"load5" json:"load5"`
	Load15               float64                    `bson:"load15" json:"load15"`
	SoftwareVersion      string                     `bson:"software_version" json:"software_version"`
	Version              int                        `bson:"version" json:"-"`
	CertificateObjs      []*certificate.Certificate `bson:"-" json:"-"`
	reqLock              sync.Mutex                 `bson:"-" json:"-"`
	reqCount             *list.List                 `bson:"-" json:"-"`
}

func (n *Node) AddRequest() {
	n.reqLock.Lock()
	back := n.reqCount.Back()
	back.Value = back.Value.(int) + 1
	n.reqLock.Unlock()
}

func (n *Node) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if n.Services == nil {
		n.Services = []bson.ObjectId{}
	}

	if n.Authorities == nil {
		n.Authorities = []bson.ObjectId{}
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

	if n.Certificates == nil || n.Protocol != "https" {
		n.Certificates = []bson.ObjectId{}
	}

	if n.Type == "" {
		n.Type = Management
	}

	if n.Type == Management {
		n.ManagementDomain = ""
		n.UserDomain = ""
	} else {
		if !strings.Contains(n.Type, Management) {
			n.ManagementDomain = ""
		}
		if !strings.Contains(n.Type, User) {
			n.UserDomain = ""
		}
	}

	if !strings.Contains(n.Type, Proxy) {
		n.Services = []bson.ObjectId{}
	}

	if !strings.Contains(n.Type, Bastion) {
		n.Authorities = []bson.ObjectId{}
	}

	n.Format()

	return
}

func (n *Node) Format() {
	utils.SortObjectIds(n.Services)
	utils.SortObjectIds(n.Certificates)
}

func (n *Node) SetActive() {
	if time.Since(n.Timestamp) > 30*time.Second {
		n.RequestsMin = 0
		n.Memory = 0
		n.Load1 = 0
		n.Load5 = 0
		n.Load15 = 0
	}
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

func (n *Node) GetRemoteAddr(r *http.Request) (addr string) {
	if n.ForwardedForHeader != "" {
		addr = strings.TrimSpace(
			strings.SplitN(r.Header.Get(n.ForwardedForHeader), ",", 1)[0])
		if addr != "" {
			return
		}
	}

	addr = utils.StripPort(r.RemoteAddr)
	return
}

func (n *Node) update(db *database.Database) (err error) {
	coll := db.Nodes()

	change := mgo.Change{
		Update: &bson.M{
			"$set": &bson.M{
				"timestamp":    n.Timestamp,
				"requests_min": n.RequestsMin,
				"memory":       n.Memory,
				"load1":        n.Load1,
				"load5":        n.Load5,
				"load15":       n.Load15,
			},
		},
		Upsert:    false,
		ReturnNew: true,
	}

	nde := &Node{}

	_, err = coll.Find(&bson.M{
		"_id": n.Id,
	}).Apply(change, nde)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	n.Id = nde.Id
	n.Name = nde.Name
	n.Type = nde.Type
	n.Port = nde.Port
	n.Protocol = nde.Protocol
	n.Certificates = nde.Certificates
	n.ManagementDomain = nde.ManagementDomain
	n.UserDomain = nde.UserDomain
	n.Services = nde.Services
	n.Authorities = nde.Authorities
	n.ForwardedForHeader = nde.ForwardedForHeader
	n.ForwardedProtoHeader = nde.ForwardedProtoHeader

	return
}

func (n *Node) loadCerts(db *database.Database) (err error) {
	certObjs := []*certificate.Certificate{}

	if n.Certificates == nil || len(n.Certificates) == 0 {
		n.CertificateObjs = certObjs
		return
	}

	for _, certId := range n.Certificates {
		cert, e := certificate.Get(db, certId)
		if e != nil {
			switch e.(type) {
			case *database.NotFoundError:
				e = nil
				break
			default:
				err = e
				return
			}
		} else {
			certObjs = append(certObjs, cert)
		}
	}

	n.CertificateObjs = certObjs

	return
}

func (n *Node) sync() {
	db := database.GetDatabase()
	defer db.Close()

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

	err = n.loadCerts(db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("node: Failed to load node certificate")
	}
}

func (n *Node) keepalive() {
	for {
		n.sync()
		time.Sleep(1 * time.Second)
	}
}

func (n *Node) reqInit() {
	n.reqLock.Lock()
	n.reqCount = list.New()
	for i := 0; i < 60; i++ {
		n.reqCount.PushBack(0)
	}
	n.reqLock.Unlock()
}

func (n *Node) reqSync() {
	for {
		time.Sleep(1 * time.Second)

		n.reqLock.Lock()

		var count int64
		for elm := n.reqCount.Front(); elm != nil; elm = elm.Next() {
			count += int64(elm.Value.(int))
		}
		n.RequestsMin = count

		n.reqCount.Remove(n.reqCount.Front())
		n.reqCount.PushBack(0)

		n.reqLock.Unlock()
	}
}

func (n *Node) Init() (err error) {
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

	n.SoftwareVersion = constants.Version

	if n.Name == "" {
		n.Name = utils.RandName()
	}

	if n.Type == "" {
		n.Type = Management
	}

	if n.Protocol == "" {
		n.Protocol = "https"
	}

	if n.Port == 0 {
		n.Port = 443
	}

	if n.Services == nil {
		n.Services = []bson.ObjectId{}
	}

	if n.Authorities == nil {
		n.Authorities = []bson.ObjectId{}
	}

	_, err = coll.UpsertId(n.Id, &bson.M{
		"$set": &bson.M{
			"_id":              n.Id,
			"name":             n.Name,
			"type":             n.Type,
			"timestamp":        time.Now(),
			"protocol":         n.Protocol,
			"port":             n.Port,
			"services":         n.Services,
			"authorities":      n.Authorities,
			"software_version": n.SoftwareVersion,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	n.reqInit()

	err = n.loadCerts(db)
	if err != nil {
		return
	}

	event.PublishDispatch(db, "node.change")

	Self = n

	go n.keepalive()
	go n.reqSync()

	return
}

func init() {
	module := requires.New("node")
	module.After("settings")

	module.Handler = func() (err error) {
		db := database.GetDatabase()
		defer db.Close()

		nodes, err := GetAll(db)
		if err != nil {
			return
		}

		for _, node := range nodes {
			if node.Version < 1 {
				changed := set.NewSet("version")
				node.Version = 1

				if node.Certificate != "" &&
					(node.Certificates == nil ||
						len(node.Certificates) == 0) {

					node.Certificates = []bson.ObjectId{
						node.Certificate,
					}
					changed.Add("certificates")
				}

				err = node.CommitFields(
					db,
					changed,
				)
				if err != nil {
					return
				}
			}
		}

		return
	}
}
