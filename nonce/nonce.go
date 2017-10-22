package nonce

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"time"
)

type nonce struct {
	Id        string    `bson:"_id"`
	Timestamp time.Time `bson:"timestamp"`
}

func Validate(db *database.Database, nce string) (err error) {
	doc := &nonce{
		Id:        nce,
		Timestamp: time.Now(),
	}

	coll := db.Nonces()

	err = coll.Insert(doc)
	if err != nil {
		switch err.(type) {
		case *database.DuplicateKeyError:
			err = &errortypes.AuthenticationError{
				errors.New("nonce: Duplicate authentication nonce"),
			}
		}
		return
	}

	return
}
