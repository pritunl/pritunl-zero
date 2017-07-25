package certificate

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"gopkg.in/mgo.v2/bson"
)

type Certificate struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Type        string        `bson:"type" json:"type"`
	Key         string        `bson:"key" json:"key"`
	Certificate string        `bson:"certificate" json:"certificate"`
	AcmeHost    bson.ObjectId `bson:"acme_host" json:"acme_host"`
	AcmeAccount string        `bson:"acme_account" json:"acme_account"`
	AcmeDomains []string      `bson:"acme_domains" json:"acme_domains"`
}

func (c *Certificate) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if c.Type != LetsEncrypt {
		c.AcmeDomains = []string{}
		c.AcmeAccount = ""
	}

	if c.AcmeDomains == nil {
		c.AcmeDomains = []string{}
	}

	return
}

func (c *Certificate) Commit(db *database.Database) (err error) {
	coll := db.Certificates()

	err = coll.Commit(c.Id, c)
	if err != nil {
		return
	}

	return
}

func (c *Certificate) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Certificates()

	err = coll.CommitFields(c.Id, c, fields)
	if err != nil {
		return
	}

	return
}

func (c *Certificate) Insert(db *database.Database) (err error) {
	coll := db.Certificates()

	if c.Id != "" {
		err = &errortypes.DatabaseError{
			errors.New("certificate: Certificate already exists"),
		}
		return
	}

	err = coll.Insert(c)
	if err != nil {
		return
	}

	return
}
