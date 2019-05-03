package device

import (
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
)

func Get(db *database.Database, devcId primitive.ObjectID) (
	devc *Device, err error) {

	coll := db.Devices()
	devc = &Device{}

	err = coll.FindOneId(devcId, devc)
	if err != nil {
		return
	}

	return
}

func GetUser(db *database.Database, devcId primitive.ObjectID,
	userId primitive.ObjectID) (devc *Device, err error) {

	coll := db.Devices()
	devc = &Device{}

	err = coll.FindOne(db, &bson.M{
		"_id":  devcId,
		"user": userId,
	}).Decode(devc)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAll(db *database.Database, userId primitive.ObjectID) (
	devices []*Device, err error) {

	coll := db.Devices()
	devices = []*Device{}

	cursor, err := coll.Find(db, &bson.M{
		"user": userId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		devc := &Device{}
		err = cursor.Decode(devc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		devices = append(devices, devc)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAllSorted(db *database.Database, userId primitive.ObjectID) (
	devices []*Device, err error) {

	coll := db.Devices()
	devices = []*Device{}

	cursor, err := coll.Find(db, &bson.M{
		"user": userId,
	}, &options.FindOptions{
		Sort: &bson.D{
			{"mode", 1},
			{"name", 1},
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		devc := &Device{}
		err = cursor.Decode(devc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		devices = append(devices, devc)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAllMode(db *database.Database, userId primitive.ObjectID,
	mode string) (devices []*Device, err error) {

	coll := db.Devices()
	devices = []*Device{}

	cursor, err := coll.Find(db, &bson.M{
		"user": userId,
		"mode": mode,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		devc := &Device{}
		err = cursor.Decode(devc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		devices = append(devices, devc)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Count(db *database.Database, userId primitive.ObjectID) (
	count int64, err error) {

	coll := db.Devices()

	count, err = coll.CountDocuments(db, &bson.M{
		"user": userId,
		"mode": Secondary,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func New(userId primitive.ObjectID, typ, mode string) (devc *Device) {
	devc = &Device{
		Id:         primitive.NewObjectID(),
		Type:       typ,
		Mode:       mode,
		User:       userId,
		Timestamp:  time.Now(),
		LastActive: time.Now(),
	}

	return
}

func Remove(db *database.Database, id primitive.ObjectID) (err error) {
	coll := db.Devices()

	_, err = coll.DeleteOne(db, &bson.M{
		"_id": id,
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

func RemoveUser(db *database.Database, id primitive.ObjectID,
	userId primitive.ObjectID) (err error) {

	coll := db.Devices()

	_, err = coll.DeleteOne(db, &bson.M{
		"_id":  id,
		"user": userId,
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

func RemoveAll(db *database.Database, userId primitive.ObjectID) (err error) {
	coll := db.Devices()

	_, err = coll.DeleteMany(db, &bson.M{
		"user": userId,
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
