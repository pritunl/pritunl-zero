package user

import (
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
)

func Get(db *database.Database, userId bson.ObjectId) (
	usr *User, err error) {

	coll := db.Users()
	usr = &User{}

	err = coll.FindOneId(userId, usr)
	if err != nil {
		return
	}

	return
}

func GetUsername(db *database.Database, typ, username string) (
	usr *User, err error) {

	coll := db.Users()
	usr = &User{}

	err = coll.FindOne(&bson.M{
		"type":     typ,
		"username": username,
	}, usr)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database, query *bson.M, page, pageCount int) (
	users []*User, count int, err error) {

	coll := db.Users()
	users = []*User{}

	qury := coll.Find(query)

	count, err = qury.Count()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	skip := utils.Min(page*pageCount, utils.Max(0, count-pageCount))

	cursor := qury.Skip(skip).Limit(pageCount).Iter()

	usr := &User{}
	for cursor.Next(usr) {
		users = append(users, usr)
		usr = &User{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func HasSuper(db *database.Database) (exists bool, err error) {
	coll := db.Users()

	count, err := coll.Find(bson.M{
		"administrator": "super",
	}).Count()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	if count > 0 {
		exists = true
	}

	return
}

func hasSuperSkip(db *database.Database, skipId bson.ObjectId) (
	exists bool, err error) {

	coll := db.Users()

	count, err := coll.Find(&bson.M{
		"_id": &bson.M{
			"$ne": skipId,
		},
		"administrator": "super",
	}).Count()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	if count > 0 {
		exists = true
	}

	return
}
