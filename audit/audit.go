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
	User      bson.ObjectId `bson:"u" json:"user"`
	Timestamp time.Time     `bson:"t" json:"timestamp"`
	Type      string        `bson:"y" json:"type"`
	Message   string        `bson:"m" json:"message"`
	Agent     *agent.Agent  `bson:"a" json:"agent"`
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
