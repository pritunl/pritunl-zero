package sshcert

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/database"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Certificate struct {
	Id           bson.ObjectId   `bson:"_id,omitempty" json:"id"`
	UserId       bson.ObjectId   `bson:"user_id" json:"user_id"`
	AuthorityIds []bson.ObjectId `bson:"authority_ids" json:"authority_ids"`
	Timestamp    time.Time       `bson:"timestamp" json:"timestamp"`
	Certificates []string        `bson:"certificates" json:"certificates"`
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
