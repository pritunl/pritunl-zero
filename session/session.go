// Stores sessions in cookies.
package session

import (
	"crypto/hmac"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/rokey"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/useragent"
	"github.com/pritunl/pritunl-zero/utils"
)

type Session struct {
	Id         string             `bson:"_id" json:"id"`
	Type       string             `bson:"type" json:"type"`
	User       primitive.ObjectID `bson:"user" json:"user"`
	Rokey      primitive.ObjectID `bson:"rokey" json:"-"`
	Secret     string             `bson:"secret" json:"-"`
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
	LastActive time.Time          `bson:"last_active" json:"last_active"`
	Removed    bool               `bson:"removed" json:"removed"`
	Agent      *useragent.Agent   `bson:"agent" json:"agent"`
	user       *user.User         `bson:"-" json:"-"`
}

func (s *Session) CheckSignature(db *database.Database, inSig string) (
	valid bool, err error) {

	if s.Rokey.IsZero() || s.Secret == "" {
		return
	}

	rkey, err := rokey.GetId(db, s.Type, s.Rokey)
	if err != nil {
		return
	}

	if rkey == nil {
		return
	}

	if rkey.Secret == "" {
		err = &errortypes.ReadError{
			errors.Wrap(err, "session: Empty secret"),
		}
		return
	}

	hash := hmac.New(sha512.New, []byte(rkey.Secret))
	hash.Write([]byte(s.Secret))
	outSig := base64.RawStdEncoding.EncodeToString(hash.Sum(nil))

	if subtle.ConstantTimeCompare([]byte(inSig), []byte(outSig)) == 1 {
		valid = true
	}

	return
}

func (s *Session) GenerateSignature(db *database.Database) (
	sig string, err error) {

	rkey, err := rokey.Get(db, s.Type)
	if err != nil {
		return
	}

	s.Rokey = rkey.Id

	s.Secret, err = utils.RandStr(64)
	if err != nil {
		return
	}

	if rkey.Secret == "" {
		err = &errortypes.ReadError{
			errors.Wrap(err, "session: Empty secret"),
		}
		return
	}

	hash := hmac.New(sha512.New, []byte(rkey.Secret))
	hash.Write([]byte(s.Secret))
	sig = base64.RawStdEncoding.EncodeToString(hash.Sum(nil))

	return
}

func (s *Session) Active() bool {
	if s.Removed {
		return false
	}

	expire := GetExpire(s.Type)
	maxDuration := GetMaxDuration(s.Type)

	if expire != 0 {
		if time.Since(s.LastActive) > expire {
			return false
		}
	}

	if maxDuration != 0 {
		if time.Since(s.Timestamp) > maxDuration {
			return false
		}
	}

	return true
}

func (s *Session) Update(db *database.Database) (err error) {
	coll := db.Sessions()

	err = coll.FindOneId(s.Id, s)
	if err != nil {
		return
	}

	return
}

func (s *Session) Remove(db *database.Database) (err error) {
	err = Remove(db, s.Id)
	if err != nil {
		return
	}

	return
}

func (s *Session) GetUser(db *database.Database) (usr *user.User, err error) {
	if s.user != nil || db == nil {
		usr = s.user
		return
	}

	usr, err = user.GetUpdate(db, s.User)
	if err != nil {
		return
	}

	s.user = usr

	return
}
