package authority

import (
	"github.com/pritunl/pritunl-zero/database"
	"gopkg.in/mgo.v2/bson"
)

func Get(db *database.Database, authId bson.ObjectId) (
	auth *Authority, err error) {

	coll := db.Authorities()
	auth = &Authority{}

	err = coll.FindOneId(authId, auth)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database) (auths []*Authority, err error) {
	coll := db.Authorities()
	auths = []*Authority{}

	cursor := coll.Find(bson.M{}).Iter()

	auth := &Authority{}
	for cursor.Next(auth) {
		auths = append(auths, auth)
		auth = &Authority{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, authId bson.ObjectId) (err error) {
	coll := db.Authorities()

	_, err = coll.RemoveAll(&bson.M{
		"_id": authId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
