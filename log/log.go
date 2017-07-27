package log

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Entry struct {
	Id        bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	Level     string                 `bson:"level" json:"level"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
	Message   string                 `bson:"message" json:"message"`
	Stack     string                 `bson:"stack" json:"stack"`
	Fields    map[string]interface{} `bson:"fields" json:"fields"`
}

func (e *Entry) Insert(db *database.Database) (err error) {
	coll := db.Logs()

	if e.Id != "" {
		err = &errortypes.DatabaseError{
			errors.New("log: Entry already exists"),
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
