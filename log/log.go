package log

import (
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/requires"
)

var published = false

type Entry struct {
	Id        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Level     string                 `bson:"level" json:"level"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
	Message   string                 `bson:"message" json:"message"`
	Stack     string                 `bson:"stack" json:"stack"`
	Fields    map[string]interface{} `bson:"fields" json:"fields"`
}

func (e *Entry) Insert(db *database.Database) (err error) {
	coll := db.Logs()

	if !e.Id.IsZero() {
		err = &errortypes.DatabaseError{
			errors.New("log: Entry already exists"),
		}
		return
	}

	_, err = coll.InsertOne(db, e)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	published = true

	return
}

func publish() {
	db := database.GetDatabase()
	defer db.Close()

	event.PublishDispatch(db, "log.change")
}

func initSender() {
	for {
		time.Sleep(1500 * time.Millisecond)
		if published {
			published = false
			publish()
		}
	}
}

func init() {
	module := requires.New("log")
	module.After("logger")

	module.Handler = func() (err error) {
		go initSender()
		return
	}
}
