package audit

import (
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/useragent"
)

type Fields map[string]any

type Audit struct {
	Id        bson.ObjectID    `bson:"_id,omitempty" json:"id"`
	User      bson.ObjectID    `bson:"u" json:"user"`
	Timestamp time.Time        `bson:"t" json:"timestamp"`
	Type      string           `bson:"y" json:"type"`
	Fields    Fields           `bson:"f" json:"fields"`
	Agent     *useragent.Agent `bson:"a" json:"agent"`
}

func (a *Audit) Insert(db *database.Database) (err error) {
	coll := db.Audits()

	if !a.Id.IsZero() {
		err = &errortypes.DatabaseError{
			errors.New("audit: Entry already exists"),
		}
		return
	}

	_, err = coll.InsertOne(db, a)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
