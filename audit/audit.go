package audit

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Event struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	User      bson.ObjectId `bson:"user" json:"user"`
	Timestamp time.Time     `bson:"timestamp" json:"timestamp"`
	Type      string        `bson:"type" json:"type"`
	Agent     *agent.Agent  `bson:"agent" json:"agent"`
	Message   string        `bson:"message" json:"message"`
}

func (e *Event) Insert(db *database.Database) (err error) {
	coll := db.Audits()

	if e.Id != "" {
		err = &errortypes.DatabaseError{
			errors.New("audit: Entry already exists"),
		}
		return
	}

	err = coll.Insert(e)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
