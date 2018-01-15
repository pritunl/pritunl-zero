package challenge

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/authority"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/ssh"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"time"
)

type Challenge struct {
	Id            string        `bson:"_id"`
	CertificateId bson.ObjectId `bson:"certificate_id,omitempty"`
	Timestamp     time.Time     `bson:"timestamp"`
	State         string        `bson:"state"`
	PubKey        string        `bson:"pub_key"`
}

func (c *Challenge) Approve(db *database.Database, usr *user.User,
	r *http.Request, secondary bool) (secProvider bson.ObjectId, err error,
	errData *errortypes.ErrorData) {

	allAuthrs, err := authority.GetAll(db)
	if err != nil {
		return
	}

	authrIds := []bson.ObjectId{}
	authrs := []*authority.Authority{}
	for _, authr := range allAuthrs {
		if authr.UserHasAccess(usr) {
			authrIds = append(authrIds, authr.Id)
			authrs = append(authrs, authr)
		}
	}

	policies, err := policy.GetAuthoritiesRoles(db, authrIds, usr.Roles)
	if err != nil {
		return
	}

	for _, polcy := range policies {
		errData, err = polcy.ValidateUser(db, usr, r)
		if err != nil || errData != nil {
			err = c.Deny(db, usr)
			if err != nil {
				return
			}
			return
		}
	}

	keybaseMode := policy.KeybaseMode(policies)
	if keybaseMode == policy.Required {
		err = c.Deny(db, usr)
		if err != nil {
			return
		}

		errData = &errortypes.ErrorData{
			Error:   "keybase_required",
			Message: "Keybase is required for this user",
		}
		return
	}

	for _, polcy := range policies {
		if polcy.AuthoritySecondary != "" {
			secProvider = polcy.AuthoritySecondary
			if !secondary {
				return
			} else {
				break
			}
		}
	}

	if c.State != "" {
		err = errortypes.WriteError{
			errors.New("sshcert: Challenge has already been answered"),
		}
		return
	}

	agnt, err := agent.Parse(db, r)
	if err != nil {
		return
	}

	cert, err := ssh.NewCertificate(db, authrs, usr, agnt, c.PubKey)
	if err != nil {
		return
	}

	if len(cert.Certificates) == 0 {
		c.State = ssh.Unavailable
		c.CertificateId = ""
	} else {
		err = cert.Insert(db)
		if err != nil {
			return
		}

		c.State = ssh.Approved
		c.CertificateId = cert.Id
	}

	coll := db.SshChallenges()

	err = coll.Update(&bson.M{
		"_id":   c.Id,
		"state": "",
	}, c)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (c *Challenge) Deny(db *database.Database, usr *user.User) (err error) {
	if c.State != "" {
		err = errortypes.WriteError{
			errors.New("sshcert: Challenge has already been answered"),
		}
		return
	}

	c.State = ssh.Denied
	c.CertificateId = ""

	coll := db.SshChallenges()

	err = coll.Update(&bson.M{
		"_id":   c.Id,
		"state": "",
	}, c)
	if err != nil {
		err = database.ParseError(err)
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

	pubKey = strings.TrimSpace(pubKey)

	if len(pubKey) > settings.System.SshPubKeyLen {
		err = errortypes.ParseError{
			errors.New("sshcert: Public key too long"),
		}
		return
	}

	token, err := utils.RandStr(48)
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
