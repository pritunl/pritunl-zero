package secret

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
)

type Secret struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Comment string             `bson:"comment" json:"comment"`
	Type    string             `bson:"type" json:"type"`
	Key     string             `bson:"key" json:"key"`
	Region  string             `bson:"region" json:"region"`
	Value   string             `bson:"value" json:"value"`
}

func (c *Secret) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	switch c.Type {
	case AWS, "":
		c.Type = AWS

		if c.Region == "" {
			c.Region = "us-east-1"
		}

		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "invalid_secret_type",
			Message: "Secret type invalid",
		}
		return
	}

	return
}

func (c *Secret) Commit(db *database.Database) (err error) {
	coll := db.Secrets()

	err = coll.Commit(c.Id, c)
	if err != nil {
		return
	}

	return
}

func (c *Secret) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Secrets()

	err = coll.CommitFields(c.Id, c, fields)
	if err != nil {
		return
	}

	return
}

func (c *Secret) Insert(db *database.Database) (err error) {
	coll := db.Secrets()

	if !c.Id.IsZero() {
		err = &errortypes.DatabaseError{
			errors.New("secret: Secret already exists"),
		}
		return
	}

	_, err = coll.InsertOne(db, c)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
