package keybase

import (
	"fmt"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/authority"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/sshcert"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"time"
)

type Challenge struct {
	Id        string    `bson:"_id"`
	Type      string    `bson:"type"`
	Username  string    `bson:"username"`
	Timestamp time.Time `bson:"timestamp"`
	State     string    `bson:"state"`
	PubKey    string    `bson:"pub_key"`
}

func (c *Challenge) Message() string {
	return fmt.Sprintf(
		"%s&%s&%s&%s",
		c.Id,
		c.Type,
		c.Username,
		c.PubKey,
	)
}

func (c *Challenge) Validate(db *database.Database, r *http.Request,
	signature string) (certf *sshcert.Certificate, err error,
	errData *errortypes.ErrorData) {

	if c.State != "" {
		err = errortypes.WriteError{
			errors.New("keybase: Challenge has already been answered"),
		}
	}

	valid, err := VerifySig(c.Message(), signature, c.Username)
	if err != nil {
		return
	}

	if !valid {
		errData = &errortypes.ErrorData{
			Error:   "invalid_signature",
			Message: "Keybase signature is invalid",
		}
		return
	}

	usr, err := user.GetKeybase(db, c.Username)
	if err != nil {
		return
	}

	cert := &sshcert.Certificate{
		Id:               bson.NewObjectId(),
		UserId:           usr.Id,
		AuthorityIds:     []bson.ObjectId{},
		Timestamp:        time.Now(),
		PubKey:           c.PubKey,
		Certificates:     []string{},
		CertificatesInfo: []*sshcert.Info{},
	}

	agnt, err := agent.Parse(db, r)
	if err != nil {
		return
	}
	cert.Agent = agnt

	authrs, err := authority.GetAll(db)
	if err != nil {
		return
	}

	for _, authr := range authrs {
		if !authr.UserHasAccess(usr) {
			continue
		}

		crt, certStr, e := authr.CreateCertificate(usr, c.PubKey)
		if e != nil {
			err = e
			return
		}

		info := &sshcert.Info{
			Expires:    time.Unix(int64(crt.ValidBefore), 0),
			Serial:     fmt.Sprintf("%d", crt.Serial),
			Principals: crt.ValidPrincipals,
			Extensions: []string{},
		}

		for permission := range crt.Permissions.Extensions {
			info.Extensions = append(info.Extensions, permission)
		}

		cert.AuthorityIds = append(cert.AuthorityIds, authr.Id)
		cert.Certificates = append(cert.Certificates, certStr)
		cert.CertificatesInfo = append(cert.CertificatesInfo, info)
	}

	if len(cert.Certificates) == 0 {
		c.State = Unavailable
	} else {
		err = cert.Insert(db)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		c.State = Approved
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

	certf = cert

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

func NewChallenge(db *database.Database, username, pubKey string) (
	chal *Challenge, err error) {

	pubKey = strings.TrimSpace(pubKey)

	if len(pubKey) > settings.System.SshPubKeyLen {
		err = errortypes.ParseError{
			errors.New("keybase: Public key too long"),
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
		Username:  username,
		PubKey:    pubKey,
	}

	err = chal.Insert(db)
	if err != nil {
		err = database.ParseError(err)
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
