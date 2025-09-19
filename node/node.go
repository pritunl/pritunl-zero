package node

import (
	"container/list"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/pritunl/pritunl-zero/certificate"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/sirupsen/logrus"
)

var (
	Self *Node
)

type Node struct {
	Id                   bson.ObjectID              `bson:"_id" json:"id"`
	Name                 string                     `bson:"name" json:"name"`
	Type                 string                     `bson:"type" json:"type"`
	Timestamp            time.Time                  `bson:"timestamp" json:"timestamp"`
	Port                 int                        `bson:"port" json:"port"`
	NoRedirectServer     bool                       `bson:"no_redirect_server" json:"no_redirect_server"`
	Protocol             string                     `bson:"protocol" json:"protocol"`
	Certificate          bson.ObjectID              `bson:"certificate" json:"certificate"`
	Certificates         []bson.ObjectID            `bson:"certificates" json:"certificates"`
	SelfCertificate      string                     `bson:"self_certificate_key" json:"-"`
	SelfCertificateKey   string                     `bson:"self_certificate" json:"-"`
	ManagementDomain     string                     `bson:"management_domain" json:"management_domain"`
	UserDomain           string                     `bson:"user_domain" json:"user_domain"`
	WebauthnDomain       string                     `bson:"webauthn_domain" json:"webauthn_domain"`
	EndpointDomain       string                     `bson:"endpoint_domain" json:"endpoint_domain"`
	Services             []bson.ObjectID            `bson:"services" json:"services"`
	Authorities          []bson.ObjectID            `bson:"authorities" json:"authorities"`
	RequestsMin          int64                      `bson:"requests_min" json:"requests_min"`
	ForwardedForHeader   string                     `bson:"forwarded_for_header" json:"forwarded_for_header"`
	ForwardedProtoHeader string                     `bson:"forwarded_proto_header" json:"forwarded_proto_header"`
	Memory               float64                    `bson:"memory" json:"memory"`
	Load1                float64                    `bson:"load1" json:"load1"`
	Load5                float64                    `bson:"load5" json:"load5"`
	Load15               float64                    `bson:"load15" json:"load15"`
	SoftwareVersion      string                     `bson:"software_version" json:"software_version"`
	Hostname             string                     `bson:"hostname" json:"hostname"`
	Version              int                        `bson:"version" json:"-"`
	CertificateObjs      []*certificate.Certificate `bson:"-" json:"-"`
	reqCount             *list.List                 `bson:"-" json:"-"`
	lock                 sync.Mutex                 `bson:"-" json:"-"`
}

func (n *Node) Copy() *Node {
	n.lock.Lock()
	defer n.lock.Unlock()

	nde := &Node{
		Id:                   n.Id,
		Name:                 n.Name,
		Type:                 n.Type,
		Timestamp:            n.Timestamp,
		Port:                 n.Port,
		NoRedirectServer:     n.NoRedirectServer,
		Protocol:             n.Protocol,
		Certificate:          n.Certificate,
		Certificates:         n.Certificates,
		SelfCertificate:      n.SelfCertificate,
		SelfCertificateKey:   n.SelfCertificateKey,
		ManagementDomain:     n.ManagementDomain,
		UserDomain:           n.UserDomain,
		WebauthnDomain:       n.WebauthnDomain,
		EndpointDomain:       n.EndpointDomain,
		Services:             n.Services,
		Authorities:          n.Authorities,
		RequestsMin:          n.RequestsMin,
		ForwardedForHeader:   n.ForwardedForHeader,
		ForwardedProtoHeader: n.ForwardedProtoHeader,
		Memory:               n.Memory,
		Load1:                n.Load1,
		Load5:                n.Load5,
		Load15:               n.Load15,
		SoftwareVersion:      n.SoftwareVersion,
		Hostname:             n.Hostname,
		Version:              n.Version,
		CertificateObjs:      n.CertificateObjs,
	}

	return nde
}

func (n *Node) IsManagement() bool {
	return strings.Contains(n.Type, Management)
}

func (n *Node) IsUser() bool {
	return strings.Contains(n.Type, User)
}

func (n *Node) IsProxy() bool {
	return strings.Contains(n.Type, Proxy)
}

func (n *Node) HasType(nodeType string) bool {
	return strings.Contains(n.Type, nodeType)
}

func (n *Node) AddType(nodeType string) bool {
	if n.Type == "" {
		n.Type = nodeType
		return true
	}

	types := strings.Split(n.Type, "_")
	if slices.Contains(types, nodeType) {
		return false
	}

	types = append(types, nodeType)
	n.Type = strings.Join(types, "_")

	return true
}

func (n *Node) RemoveType(nodeType string) bool {
	if n.Type == "" {
		return false
	}

	types := strings.Split(n.Type, "_")
	found := false
	newTypes := []string{}

	for _, typ := range types {
		if typ == nodeType {
			found = true
		} else {
			newTypes = append(newTypes, typ)
		}
	}

	if found {
		n.Type = strings.Join(newTypes, "_")
		return true
	}

	return false
}

func (n *Node) IsOnline() bool {
	return time.Since(n.Timestamp) <= time.Duration(
		settings.System.NodeTimestampTtl)*time.Second
}

func (n *Node) AddRequest() {
	n.lock.Lock()
	back := n.reqCount.Back()
	back.Value = back.Value.(int) + 1
	n.lock.Unlock()
}

func (n *Node) GetWebauthn(origin string, strict bool) (
	web *webauthn.WebAuthn, err error) {

	webauthnDomain := n.WebauthnDomain
	if webauthnDomain == "" {
		if strict {
			err = &errortypes.ReadError{
				errors.New("node: Webauthn domain not configured"),
			}
			return
		} else {
			userN := strings.Count(n.UserDomain, ".")
			adminN := strings.Count(n.ManagementDomain, ".")

			if userN <= adminN {
				webauthnDomain = n.UserDomain
			} else {
				webauthnDomain = n.ManagementDomain
			}
		}
	}

	web, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "Pritunl Zero",
		RPID:          webauthnDomain,
		RPOrigins:     []string{origin},
	})
	if err != nil {
		err = utils.ParseWebauthnError(err)
		return
	}

	return
}

func (n *Node) HasService(srvcId bson.ObjectID) bool {
	return slices.Contains(n.Services, srvcId)
}

func (n *Node) AddService(srvcId bson.ObjectID) bool {
	if n.HasService(srvcId) {
		return false
	}

	n.Services = append(n.Services, srvcId)
	return true
}

func (n *Node) RemoveService(srvcId bson.ObjectID) bool {
	for i, serviceId := range n.Services {
		if serviceId == srvcId {
			n.Services = append(n.Services[:i], n.Services[i+1:]...)
			return true
		}
	}

	return false
}

func (n *Node) HasCertificate(certId bson.ObjectID) bool {
	if slices.Contains(n.Certificates, certId) {
		return true
	}

	return false
}

func (n *Node) AddCertificate(certId bson.ObjectID) bool {
	if n.HasCertificate(certId) {
		return false
	}

	n.Certificates = append(n.Certificates, certId)
	return true
}

func (n *Node) RemoveCertificate(certId bson.ObjectID) bool {
	for i, certObjId := range n.Certificates {
		if certObjId == certId {
			n.Certificates = append(
				n.Certificates[:i], n.Certificates[i+1:]...)
			return true
		}
	}

	return false
}

func (n *Node) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	n.Name = utils.FilterName(n.Name)

	if n.Services == nil {
		n.Services = []bson.ObjectID{}
	}

	if n.Authorities == nil {
		n.Authorities = []bson.ObjectID{}
	}

	if n.Protocol != Http && n.Protocol != Https {
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
		n.Certificates = []bson.ObjectID{}
	}

	if n.Type == "" {
		n.Type = Management
	}

	typs := strings.Split(n.Type, "_")
	for _, typ := range typs {
		switch typ {
		case Management, User, Proxy, Bastion:
			break
		default:
			errData = &errortypes.ErrorData{
				Error:   "type_invalid",
				Message: "Invalid node type",
			}
			return
		}
	}

	if n.Type == Management {
		n.ManagementDomain = ""
		n.UserDomain = ""
		n.EndpointDomain = ""
	} else {
		if !strings.Contains(n.Type, Management) {
			n.ManagementDomain = ""
		}
		if !strings.Contains(n.Type, User) {
			n.UserDomain = ""
			n.EndpointDomain = ""
		}
	}

	if !strings.Contains(n.Type, Proxy) {
		n.Services = []bson.ObjectID{}
	}

	if !strings.Contains(n.Type, Bastion) {
		n.Authorities = []bson.ObjectID{}
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
			strings.SplitN(r.Header.Get(n.ForwardedForHeader), ",", 2)[0])
		if addr != "" {
			return
		}

		logrus.WithFields(logrus.Fields{
			"forwarded_header": n.ForwardedForHeader,
		}).Warn("node: Unsafe node forwarded header")
	}

	addr = utils.StripPort(r.RemoteAddr)
	return
}

func (n *Node) SafeGetRemoteAddr(r *http.Request) (addr string,
	header bool, valid bool) {

	if n.ForwardedForHeader != "" {
		addr = strings.TrimSpace(
			strings.SplitN(r.Header.Get(n.ForwardedForHeader), ",", 2)[0])
		if addr != "" {
			header = true
			valid = true
		}
		return
	}

	addr = utils.StripPort(r.RemoteAddr)
	valid = true
	return
}

func (n *Node) update(db *database.Database) (err error) {
	coll := db.Nodes()

	nde := &Node{}
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	err = coll.FindOneAndUpdate(
		db,
		&bson.M{
			"_id": n.Id,
		},
		&bson.M{
			"$set": &bson.M{
				"timestamp":    n.Timestamp,
				"requests_min": n.RequestsMin,
				"memory":       n.Memory,
				"load1":        n.Load1,
				"load5":        n.Load5,
				"load15":       n.Load15,
				"hostname":     n.Hostname,
			},
		},
		opts,
	).Decode(nde)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	n.Id = nde.Id
	n.Name = nde.Name
	n.Type = nde.Type
	n.Port = nde.Port
	n.NoRedirectServer = nde.NoRedirectServer
	n.Protocol = nde.Protocol
	n.Certificates = nde.Certificates
	n.SelfCertificate = nde.SelfCertificate
	n.SelfCertificateKey = nde.SelfCertificateKey
	n.ManagementDomain = nde.ManagementDomain
	n.UserDomain = nde.UserDomain
	n.EndpointDomain = nde.EndpointDomain
	n.WebauthnDomain = nde.WebauthnDomain
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

func (n *Node) sync() (nde *Node) {
	db := database.GetDatabase()
	defer db.Close()

	nde = n.Copy()
	n = nde

	n.Timestamp = time.Now()

	mem, err := utils.GetMemInfo()
	if err != nil {
		n.Memory = 0

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("node: Failed to get memory")
	} else {
		n.Memory = mem.UsedPercent
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

	hostname, err := os.Hostname()
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "node: Failed to get hostname"),
		}
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("node: Failed to get hostname")
	}
	n.Hostname = hostname

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

	var reqCount *list.List
	if Self != nil {
		Self.lock.Lock()
		reqCount = utils.CopyList(Self.reqCount)
		Self.lock.Unlock()
	} else {
		reqCount = list.New()
		for i := 0; i < 60; i++ {
			reqCount.PushBack(0)
		}
	}

	var count int64
	for elm := reqCount.Front(); elm != nil; elm = elm.Next() {
		count += int64(elm.Value.(int))
	}
	n.RequestsMin = count

	reqCount.Remove(reqCount.Front())
	reqCount.PushBack(0)

	n.reqCount = reqCount

	Self = n

	return
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
		n.Services = []bson.ObjectID{}
	}

	if n.Authorities == nil {
		n.Authorities = []bson.ObjectID{}
	}

	_, err = coll.UpdateOne(
		db,
		&bson.M{
			"_id": n.Id,
		},
		&bson.M{
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
		},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	n.sync()

	_ = event.PublishDispatch(db, "node.change")

	go func() {
		nde := n
		for {
			if constants.Interrupt {
				return
			}

			nde = nde.sync()
			time.Sleep(1 * time.Second)
		}
	}()

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

				if !node.Certificate.IsZero() &&
					(node.Certificates == nil ||
						len(node.Certificates) == 0) {

					node.Certificates = []bson.ObjectID{
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
