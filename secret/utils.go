package secret

import (
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
)

func Get(db *database.Database, secrId primitive.ObjectID) (
	secr *Secret, err error) {

	coll := db.Secrets()
	secr = &Secret{}

	err = coll.FindOneId(secrId, secr)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database) (secrs []*Secret, err error) {
	coll := db.Secrets()
	secrs = []*Secret{}

	cursor, err := coll.Find(db, bson.M{})
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		secr := &Secret{}
		err = cursor.Decode(secr)
		if err != nil {
			return
		}

		secrs = append(secrs, secr)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, secrId primitive.ObjectID) (err error) {
	coll := db.Secrets()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": secrId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
