package authority

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"gopkg.in/mgo.v2/bson"
)

type Authority struct {
	Id    bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name  string        `bson:"name" json:"name"`
	Roles []string      `bson:"roles" json:"roles"`
}

func (a *Authority) Commit(db *database.Database) (err error) {
	coll := db.Certificates()

	err = coll.Commit(a.Id, a)
	if err != nil {
		return
	}

	return
}

func (a *Authority) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Certificates()

	err = coll.CommitFields(a.Id, a, fields)
	if err != nil {
		return
	}

	return
}

func (a *Authority) Insert(db *database.Database) (err error) {
	coll := db.Certificates()

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
