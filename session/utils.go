package session

import (
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func Get(db *database.Database, id string) (
	sess *Session, err error) {

	coll := db.Sessions()
	sess = &Session{}

	err = coll.FindOneId(id, sess)
	if err != nil {
		return
	}

	return
}

func New(db *database.Database, userId bson.ObjectId) (
	sess *Session, err error) {

	id, err := utils.RandStr(32)
	if err != nil {
		return
	}

	coll := db.Sessions()
	sess = &Session{
		Id:        id,
		UserId:    userId,
		Timestamp: time.Now(),
	}

	err = coll.Insert(sess)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, id string) (err error) {
	coll := db.Sessions()

	err = coll.RemoveId(id)
	if err != nil {
		err = database.ParseError(err)

		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		default:
			return
		}
	}

	return
}
