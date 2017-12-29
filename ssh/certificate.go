package ssh

import (
	"fmt"
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/authority"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Info struct {
	Serial     string    `bson:"serial" json:"serial"`
	Expires    time.Time `bson:"expires" json:"expires"`
	Principals []string  `bson:"principals" json:"principals"`
	Extensions []string  `bson:"extensions" json:"extensions"`
}

type Host struct {
	Domain             string `bson:"domain" json:"domain"`
	ProxyHost          string `bson:"proxy_host" json:"proxy_host"`
	StrictHostChecking bool   `bson:"strict_host_checking" json:"strict_host_checking"`
}

type Certificate struct {
	Id                     bson.ObjectId   `bson:"_id,omitempty" json:"id"`
	UserId                 bson.ObjectId   `bson:"user_id,omitempty" json:"user_id"`
	AuthorityIds           []bson.ObjectId `bson:"authority_ids" json:"authority_ids"`
	Timestamp              time.Time       `bson:"timestamp" json:"timestamp"`
	PubKey                 string          `bson:"pub_key"`
	Hosts                  []*Host         `bson:"hosts" json:"hosts"`
	CertificateAuthorities []string        `bson:"certificate_authorities" json:"-"`
	Certificates           []string        `bson:"certificates" json:"-"`
	CertificatesInfo       []*Info         `bson:"certificates_info" json:"certificates_info"`
	Agent                  *agent.Agent    `bson:"agent" json:"agent"`
}

func (c *Certificate) Commit(db *database.Database) (err error) {
	coll := db.SshCertificates()

	err = coll.Commit(c.Id, c)
	if err != nil {
		return
	}

	return
}

func (c *Certificate) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.SshCertificates()

	err = coll.CommitFields(c.Id, c, fields)
	if err != nil {
		return
	}

	return
}

func (c *Certificate) Insert(db *database.Database) (err error) {
	coll := db.SshCertificates()

	err = coll.Insert(c)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetCertificate(db *database.Database, certId bson.ObjectId) (
	cert *Certificate, err error) {

	coll := db.SshCertificates()
	cert = &Certificate{}

	err = coll.FindOneId(certId, cert)
	if err != nil {
		return
	}

	return
}

func GetCertificates(db *database.Database, userId bson.ObjectId,
	page, pageCount int) (certs []*Certificate, count int, err error) {

	coll := db.SshCertificates()
	certs = []*Certificate{}

	qury := coll.Find(&bson.M{
		"user_id": userId,
	})

	count, err = qury.Count()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	skip := utils.Min(page*pageCount, utils.Max(0, count-pageCount))

	cursor := qury.Sort("-timestamp").Skip(skip).Limit(pageCount).Iter()

	cert := &Certificate{}
	for cursor.Next(cert) {
		certs = append(certs, cert)
		cert = &Certificate{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func NewCertificate(db *database.Database, authrs []*authority.Authority,
	usr *user.User, agnt *agent.Agent, pubKey string) (cert *Certificate,
	err error) {

	cert = &Certificate{
		Id:           bson.NewObjectId(),
		UserId:       usr.Id,
		AuthorityIds: []bson.ObjectId{},
		Timestamp:    time.Now(),
		PubKey:       pubKey,
		Hosts:        []*Host{},
		CertificateAuthorities: []string{},
		Certificates:           []string{},
		CertificatesInfo:       []*Info{},
		Agent:                  agnt,
	}

	for _, authr := range authrs {
		if !authr.UserHasAccess(usr) {
			continue
		}

		crt, certStr, e := authr.CreateCertificate(usr, pubKey)
		if e != nil {
			err = e
			return
		}

		info := &Info{
			Expires:    time.Unix(int64(crt.ValidBefore), 0),
			Serial:     fmt.Sprintf("%d", crt.Serial),
			Principals: crt.ValidPrincipals,
			Extensions: []string{},
		}

		for permission := range crt.Permissions.Extensions {
			info.Extensions = append(info.Extensions, permission)
		}

		certAuthr := authr.GetCertAuthority()
		if certAuthr != "" {
			cert.CertificateAuthorities = append(
				cert.CertificateAuthorities,
				certAuthr,
			)
		}

		if authr.StrictHostChecking || authr.HostProxy != "" {
			hst := &Host{
				Domain:             authr.GetHostDomain(),
				ProxyHost:          authr.HostProxy,
				StrictHostChecking: authr.StrictHostChecking,
			}
			cert.Hosts = append(cert.Hosts, hst)
		}

		cert.AuthorityIds = append(cert.AuthorityIds, authr.Id)
		cert.Certificates = append(cert.Certificates, certStr)
		cert.CertificatesInfo = append(cert.CertificatesInfo, info)
	}

	return
}
