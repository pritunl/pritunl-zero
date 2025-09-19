package acme

import (
	"time"

	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
)

type Challenge struct {
	Id        string    `bson:"_id"`
	Resource  string    `bson:"resource"`
	Timestamp time.Time `bson:"timestamp"`
}

func (c *Challenge) Insert(db *database.Database) (err error) {
	coll := db.AcmeChallenges()

	_, err = coll.InsertOne(db, c)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (c *Challenge) Remove(db *database.Database) (err error) {
	coll := db.AcmeChallenges()

	_, err = coll.DeleteOne(db, &bson.M{
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

type ChallengeMsg struct {
	Token    string `bson:"token" json:"token"`
	Response string `bson:"response" json:"response"`
}

func (c *ChallengeMsg) Publish(db *database.Database) (err error) {
	err = event.Publish(db, "acme", c)
	if err != nil {
		return
	}

	return
}
