package authority

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type Authority struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Type       string        `bson:"type" json:"type"`
	Roles      []string      `bson:"roles" json:"roles"`
	PrivateKey string        `bson:"private_key" json:"private_key"`
}

func (a *Authority) GenerateRsaPrivateKey() (err error) {
	keyBytes, err := GenerateRsaKey()
	if err != nil {
		return
	}

	a.PrivateKey = strings.TrimSpace(string(keyBytes))

	return
}

func (a *Authority) GenerateEcPrivateKey() (err error) {
	keyBytes, err := GenerateEcKey()
	if err != nil {
		return
	}

	a.PrivateKey = strings.TrimSpace(string(keyBytes))

	return
}

func (a *Authority) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if a.Type == "" {
		a.Type = Local
	}

	return
}

func (a *Authority) Commit(db *database.Database) (err error) {
	coll := db.Authorities()

	err = coll.Commit(a.Id, a)
	if err != nil {
		return
	}

	return
}

func (a *Authority) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Authorities()

	err = coll.CommitFields(a.Id, a, fields)
	if err != nil {
		return
	}

	return
}

func (a *Authority) Insert(db *database.Database) (err error) {
	coll := db.Authorities()

	if a.Id != "" {
		err = &errortypes.DatabaseError{
			errors.New("authority: Authority already exists"),
		}
		return
	}

	err = coll.Insert(a)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
