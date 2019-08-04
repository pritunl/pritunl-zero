// Stores sessions in cookies.
package session

import (
	"time"

	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/user"
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
	Agent      *agent.Agent       `bson:"agent" json:"agent"`
	user       *user.User         `bson:"-" json:"-"`
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
