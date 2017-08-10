package audit

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type Event struct {
	Id        bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	User      bson.ObjectId          `bson:"u" json:"user"`
	Timestamp time.Time              `bson:"t" json:"timestamp"`
	Type      string                 `bson:"y" json:"type"`
	Fields    map[string]interface{} `bson:"f" json:"fields"`
	Agent     *agent.Agent           `bson:"a" json:"agent"`
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

func NewEvent(db *database.Database, r *http.Request,
	userId bson.ObjectId, typ, msg string) (err error) {

	agnt, err := agent.Parse(db, r)
	if err != nil {
		return
	}

	evt := &Event{
		User:      userId,
		Timestamp: time.Now(),
		Type:      typ,
		Message:   msg,
		Agent:     agnt,
	}

	err = evt.Insert(db)
	if err != nil {
		return
	}

	return
}
