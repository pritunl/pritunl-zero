package sshcert

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/authority"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Challenge struct {
	Id            string        `bson:"_id"`
	CertificateId bson.ObjectId `bson:"certificate_id,omitempty"`
	Timestamp     time.Time     `bson:"timestamp"`
	State         string        `bson:"state"`
	PubKey        string        `bson:"pub_key"`
}

func (c *Challenge) Approve(db *database.Database, usr *user.User) (
	err error) {

	if c.State != "" {
		err = errortypes.WriteError{
			errors.New("sshcert: Challenge has already been answered"),
		}
	}

	cert := &Certificate{
		Id:           bson.NewObjectId(),
		UserId:       usr.Id,
		AuthorityIds: []bson.ObjectId{},
		Timestamp:    time.Now(),
		Certificates: []string{},
	}

	authrs, err := authority.GetAll(db)
	if err != nil {
		return
	}

	for _, authr := range authrs {
		if !authr.UserHasAccess(usr) {
			continue
		}

		certStr, e := authr.CreateCertificate(usr, c.PubKey)
		if e != nil {
			err = e
			return
		}

		cert.AuthorityIds = append(cert.AuthorityIds, authr.Id)
		cert.Certificates = append(cert.Certificates, certStr)
	}

	if len(cert.Certificates) == 0 {
		c.State = Unavailable
		c.CertificateId = ""
	} else {
		err = cert.Insert(db)
		if err != nil {
			return
		}

		c.State = Approved
		c.CertificateId = cert.Id
	}

	coll := db.SshChallenges()

	err = coll.Update(&bson.M{
		"_id":   c.Id,
		"state": "",
	}, c)
	if err != nil {
		return
	}

	return
}

func (c *Challenge) Deny(db *database.Database, usr *user.User) (err error) {
	if c.State != "" {
		err = errortypes.WriteError{
			errors.New("sshcert: Challenge has already been answered"),
		}
	}

	c.State = Denied
	c.CertificateId = ""

	coll := db.SshChallenges()

	err = coll.Update(&bson.M{
		"_id":   c.Id,
		"state": "",
	}, c)
	if err != nil {
		return
	}

	return
}

func (c *Challenge) Commit(db *database.Database) (err error) {
	coll := db.SshChallenges()

	err = coll.Commit(c.Id, c)
	if err != nil {
		return
	}

	return
}

func (c *Challenge) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.SshChallenges()

	err = coll.CommitFields(c.Id, c, fields)
	if err != nil {
		return
	}

	return
}

func (c *Challenge) Insert(db *database.Database) (err error) {
	coll := db.SshChallenges()

	err = coll.Insert(c)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func NewChallenge(db *database.Database, pubKey string) (
	chal *Challenge, err error) {

	if len(pubKey) > settings.System.SshPubKeyLen {
		err = errortypes.ParseError{
			errors.New("sshcert: Public key too long"),
		}
		return
	}

	token, err := utils.RandStr(32)
	if err != nil {
		return
	}

	chal = &Challenge{
		Id:        token,
		Timestamp: time.Now(),
		PubKey:    pubKey,
	}

	err = chal.Insert(db)
	if err != nil {
		return
	}

	return
}

func GetChallenge(db *database.Database, chalId string) (
	chal *Challenge, err error) {

	coll := db.SshChallenges()
	chal = &Challenge{}

	err = coll.FindOneId(chalId, chal)
	if err != nil {
		return
	}

	return
}
