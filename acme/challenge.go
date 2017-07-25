package acme

import (
	"github.com/pritunl/pritunl-zero/database"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Challenge struct {
	Id        string    `bson:"_id"`
	Resource  string    `bson:"resource"`
	Timestamp time.Time `bson:"timestamp"`
}

func (c *Challenge) Insert(db *database.Database) (err error) {
	coll := db.AcmeChallenges()

	err = coll.Insert(c)
	if err != nil {
		database.ParseError(err)
		return
	}

	return
}

func (c *Challenge) Remove(db *database.Database) (err error) {
	coll := db.AcmeChallenges()

	err = coll.Remove(&bson.M{
		"_id": c.Id,
	})
	if err != nil {
		err = database.ParseError(err)
		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		default:
			return
		}
	}

	return
}
