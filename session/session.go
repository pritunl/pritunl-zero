// Stores sessions in cookies.
package session

import (
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/user"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Session struct {
	Id        string        `bson:"_id" json:"-"`
	UserId    bson.ObjectId `bson:"user_id" json:"-"`
	Timestamp time.Time     `bson:"timestamp" json:"-"`
	user      *user.User    `bson:"-" json:"-"`
}

func (s *Session) Remove(db *database.Database) (err error) {
	err = Remove(db, s.Id)
	if err != nil {
		return
	}

	return
}

func (s *Session) GetUser(db *database.Database) (usr *user.User, err error) {
	if s.user != nil {
		usr = s.user
		return
	}

	usr, err = user.Get(db, s.UserId)
	if err != nil {
		return
	}

	s.user = usr

	return
}
