package keybase

import (
	"fmt"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/agent"
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
	signature string) (certf *ssh.Certificate, err error,
	errData *errortypes.ErrorData) {

	if c.State != "" {
		err = errortypes.WriteError{
			errors.New("keybase: Challenge has already been answered"),
		}
		return
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
		if _, ok := err.(*database.NotFoundError); ok {
			err = nil
			errData = &errortypes.ErrorData{
				Error:   "invalid_keybase",
				Message: "Keybase username is invalid",
			}
		}
		return
	}

	data, err := getInfo(c.Username)
	if err != nil {
		return
	}

	if data.Them.PublicKeys.Primary.UkbId != usr.KeybaseId {
		errData = &errortypes.ErrorData{
			Error: "keybase_id_changed",
			Message: "Keybase identity has changed, " +
				"contact administrator to reset",
		}
		return
	}

	keybaseMode, err := policy.UserKeybaseMode(db, usr)
	if err != nil {
		return
	}

	if keybaseMode == policy.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "keybase_disabled",
			Message: "Keybase cannot be used with this user",
		}
		return
	}

	agnt, err := agent.Parse(db, r)
	if err != nil {
		return
	}

	cert, err := ssh.NewCertificate(db, usr, agnt, c.PubKey)
	if err != nil {
		return
	}

	if len(cert.Certificates) == 0 {
		c.State = ssh.Unavailable
	} else {
		err = cert.Insert(db)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		c.State = ssh.Approved
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

	if len(cert.Certificates) == 0 {
		errData = &errortypes.ErrorData{
			Error: "certificate_unavailable",
			Message: "Cerification was approved but no " +
				"certificates are available",
		}
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
		Type:      CertificateChallenge,
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
