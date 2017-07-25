package certificate

import (
	"crypto/md5"
	"fmt"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"io"
)

type Certificate struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Type        string        `bson:"type" json:"type"`
	Key         string        `bson:"key" json:"key"`
	Certificate string        `bson:"certificate" json:"certificate"`
	AcmeHash    string        `bson:"acme_hash" json:"acme_hash"`
	AcmeAccount string        `bson:"acme_account" json:"acme_account"`
	AcmeDomains []string      `bson:"acme_domains" json:"acme_domains"`
}

func (c *Certificate) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if c.Type != LetsEncrypt {
		c.AcmeAccount = ""
		c.AcmeDomains = []string{}
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

func (c *Certificate) Hash() string {
	hash := md5.New()
	io.WriteString(hash, c.Type)
	io.WriteString(hash, c.Key)
	io.WriteString(hash, c.Certificate)
	io.WriteString(hash, c.AcmeAccount)
	if c.AcmeDomains != nil {
		for _, domain := range c.AcmeDomains {
			io.WriteString(hash, domain)
		}
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (c *Certificate) Write() (err error) {
	err = utils.CreateWrite(constants.KeyPath, c.Key, 0600)
	if err != nil {
		return
	}

	err = utils.CreateWrite(constants.CertPath, c.Certificate, 0666)
	if err != nil {
		return
	}

	return
}
