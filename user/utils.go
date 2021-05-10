package user

import (
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/utils"
)

func Get(db *database.Database, userId primitive.ObjectID) (
	usr *User, err error) {

	coll := db.Users()
	usr = &User{}

	err = coll.FindOneId(userId, usr)
	if err != nil {
		return
	}

	return
}

func GetUpdate(db *database.Database, userId primitive.ObjectID) (
	usr *User, err error) {

	coll := db.Users()
	usr = &User{}
	timestamp := time.Now()

	err = coll.FindOneAndUpdate(
		db,
		&bson.M{
			"_id": userId,
		},
		&bson.M{
			"$set": &bson.M{
				"last_active": timestamp,
			},
		},
	).Decode(usr)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	usr.LastActive = timestamp

	return
}

func GetTokenUpdate(db *database.Database, token string) (
	usr *User, err error) {

	coll := db.Users()
	usr = &User{}
	timestamp := time.Now()

	err = coll.FindOneAndUpdate(
		db,
		&bson.M{
			"token": token,
		},
		&bson.M{
			"$set": &bson.M{
				"last_active": timestamp,
			},
		},
	).Decode(usr)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	usr.LastActive = timestamp

	return
}

func GetUsername(db *database.Database, typ, username string) (
	usr *User, err error) {

	coll := db.Users()
	usr = &User{}

	if username == "" {
		err = &errortypes.NotFoundError{
			errors.New("user: Username empty"),
		}
		return
	}

	err = coll.FindOne(db, &bson.M{
		"type":     typ,
		"username": username,
	}).Decode(usr)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAll(db *database.Database, query *bson.M, page, pageCount int64) (
	users []*User, count int64, err error) {

	coll := db.Users()
	users = []*User{}

	count, err = coll.CountDocuments(db, query)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	opts := options.FindOptions{
		Sort: &bson.D{
			{"username", 1},
		},
	}

	if pageCount != 0 {
		maxPage := count / pageCount
		if count == pageCount {
			maxPage = 0
		}
		page = utils.Min64(page, maxPage)
		skip := utils.Min64(page*pageCount, count)
		opts.Skip = &skip
		opts.Limit = &pageCount
	}

	cursor, err := coll.Find(
		db,
		query,
		&opts,
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		usr := &User{}
		err = cursor.Decode(usr)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		users = append(users, usr)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, userIds []primitive.ObjectID) (
	errData *errortypes.ErrorData, err error) {

	coll := db.Users()
	opts := &options.CountOptions{}
	opts.SetLimit(1)

	count, err := coll.CountDocuments(
		db,
		&bson.M{
			"_id": &bson.M{
				"$nin": userIds,
			},
			"administrator": "super",
		},
		opts,
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	if count == 0 {
		errData = &errortypes.ErrorData{
			Error:   "user_remove_super",
			Message: "Cannot remove all super administrators",
		}
		return
	}

	coll = db.Sessions()

	_, err = coll.DeleteMany(db, &bson.M{
		"user": &bson.M{
			"$in": userIds,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	coll = db.Users()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": &bson.M{
			"$in": userIds,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Count(db *database.Database) (count int64, err error) {
	coll := db.Users()

	count, err = coll.CountDocuments(db, &bson.M{})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func hasSuperSkip(db *database.Database, skipId primitive.ObjectID) (
	exists bool, err error) {

	coll := db.Users()
	opts := &options.CountOptions{}
	opts.SetLimit(1)

	count, err := coll.CountDocuments(
		db,
		&bson.M{
			"_id": &bson.M{
				"$ne": skipId,
			},
			"administrator": "super",
		},
		opts,
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	if count > 0 {
		exists = true
	}

	return
}
