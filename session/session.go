// Stores sessions in cookies.
package session

import (
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Session struct {
	Id         string        `bson:"_id" json:"id"`
	User       bson.ObjectId `bson:"user" json:"user"`
	Timestamp  time.Time     `bson:"timestamp" json:"timestamp"`
	LastActive time.Time     `bson:"last_active" json:"last_active"`
	Removed    bool          `bson:"removed" json:"removed"`
	Agent      *agent.Agent  `bson:"agent" json:"agent"`
	user       *user.User    `bson:"-" json:"-"`
}

func (s *Session) Active() bool {
	if s.Removed {
		return false
	}

	if settings.Auth.Expire != 0 {
		if time.Since(s.LastActive) > time.Duration(
			settings.Auth.Expire)*time.Hour {

			return false
		}
	}

	if settings.Auth.MaxDuration != 0 {
		if time.Since(s.Timestamp) > time.Duration(
			settings.Auth.MaxDuration)*time.Hour {

			return false
		}
	}

	return true
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

	usr, err = user.GetUpdate(db, s.User)
	if err != nil {
		return
	}

	s.user = usr

	return
}
