package audit

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Fields map[string]interface{}

type Audit struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	User      bson.ObjectId `bson:"u" json:"user"`
	Timestamp time.Time     `bson:"t" json:"timestamp"`
	Type      string        `bson:"y" json:"type"`
	Fields    Fields        `bson:"f" json:"fields"`
	Agent     *agent.Agent  `bson:"a" json:"agent"`
}

func (a *Audit) Insert(db *database.Database) (err error) {
	coll := db.Audits()

	if a.Id != "" {
		err = &errortypes.DatabaseError{
			errors.New("audit: Entry already exists"),
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
