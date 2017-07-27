package policy

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"gopkg.in/mgo.v2/bson"
)

type Rule struct {
	Type  string `bson:"type" json:"type"`
	Value string `bson:"value" json:"value"`
}

type Policy struct {
	Id       bson.ObjectId   `bson:"_id,omitempty" json:"id"`
	Name     string          `bson:"name" json:"name"`
	Services []bson.ObjectId `bson:"services" json:"services"`
	Roles    []string        `bson:"roles" json:"roles"`
	Rules    []*Rule         `bson:"rules" json:"rules"`
}

func (p *Policy) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	return
}

func (p *Policy) Commit(db *database.Database) (err error) {
	coll := db.Policies()

	err = coll.Commit(p.Id, p)
	if err != nil {
		return
	}

	return
}

func (p *Policy) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Policies()

	err = coll.CommitFields(p.Id, p, fields)
	if err != nil {
		return
	}

	return
}

func (p *Policy) Insert(db *database.Database) (err error) {
	coll := db.Policies()

	if p.Id != "" {
		err = &errortypes.DatabaseError{
			errors.New("policy: Policy already exists"),
		}
		return
	}

	err = coll.Insert(p)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
