package device

import (
	"time"

	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
)

func Get(db *database.Database, devcId bson.ObjectID) (
	devc *Device, err error) {

	coll := db.Devices()
	devc = &Device{}

	err = coll.FindOneId(devcId, devc)
	if err != nil {
		return
	}

	return
}

func GetUser(db *database.Database, devcId bson.ObjectID,
	userId bson.ObjectID) (devc *Device, err error) {

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

func GetAll(db *database.Database, userId bson.ObjectID) (
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

func GetAllSorted(db *database.Database, userId bson.ObjectID) (
	devices []*Device, err error) {

	coll := db.Devices()
	devices = []*Device{}

	cursor, err := coll.Find(db, bson.M{
		"user": userId,
	}, options.Find().
		SetSort(bson.D{
			{"mode", 1},
			{"name", 1},
		}))
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

func GetAllMode(db *database.Database, userId bson.ObjectID,
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

func CountSecondary(db *database.Database, userId bson.ObjectID) (
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

func New(userId bson.ObjectID, typ, mode string) (devc *Device) {
	devc = &Device{
		Id:         bson.NewObjectID(),
		Type:       typ,
		Mode:       mode,
		User:       userId,
		Timestamp:  time.Now(),
		LastActive: time.Now(),
	}

	return
}

func Remove(db *database.Database, id bson.ObjectID) (err error) {
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

func RemoveUser(db *database.Database, id bson.ObjectID,
	userId bson.ObjectID) (err error) {

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

func RemoveAll(db *database.Database, userId bson.ObjectID) (err error) {
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
