package authority

import (
	"github.com/pritunl/pritunl-zero/database"
	"gopkg.in/mgo.v2/bson"
)

func Get(db *database.Database, authrId bson.ObjectId) (
	authr *Authority, err error) {

	coll := db.Authorities()
	authr = &Authority{}

	err = coll.FindOneId(authrId, authr)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database) (authrs []*Authority, err error) {
	coll := db.Authorities()
	authrs = []*Authority{}

	cursor := coll.Find(bson.M{}).Iter()

	authr := &Authority{}
	for cursor.Next(authr) {
		authrs = append(authrs, authr)
		authr = &Authority{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, authrId bson.ObjectId) (err error) {
	coll := db.Authorities()

	_, err = coll.RemoveAll(&bson.M{
		"_id": authrId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
