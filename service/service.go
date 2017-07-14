package service

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"gopkg.in/mgo.v2/bson"
	"sort"
)

type Service struct {
	Id    bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name  string        `bson:"name" json:"name"`
	Roles []string      `bson:"roles" json:"roles"`
}

func (s *Service) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	s.Format()

	return
}

func (s *Service) Format() {
	sort.Strings(s.Roles)
}

func (s *Service) Commit(db *database.Database) (err error) {
	coll := db.Services()

	err = coll.Commit(s.Id, s)
	if err != nil {
		return
	}

	return
}

func (s *Service) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Services()

	err = coll.CommitFields(s.Id, s, fields)
	if err != nil {
		return
	}

	return
}

func (s *Service) Insert(db *database.Database) (err error) {
	coll := db.Services()

	if s.Id != "" {
		err = &errortypes.DatabaseError{
			errors.New("service: Service already exists"),
		}
		return
	}

	err = coll.Insert(s)
	if err != nil {
		return
	}

	return
}
