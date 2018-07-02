package device

import (
	"github.com/pritunl/pritunl-zero/database"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func Get(db *database.Database, devcId bson.ObjectId) (
	devc *Device, err error) {

	coll := db.Devices()
	devc = &Device{}

	err = coll.FindOneId(devcId, devc)
	if err != nil {
		return
	}

	return
}

func GetUser(db *database.Database, devcId bson.ObjectId,
	userId bson.ObjectId) (devc *Device, err error) {

	coll := db.Devices()
	devc = &Device{}

	err = coll.FindOne(&bson.M{
		"_id":  devcId,
		"user": userId,
	}, devc)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database, userId bson.ObjectId) (
	devices []*Device, err error) {

	coll := db.Devices()
	devices = []*Device{}

	cursor := coll.Find(&bson.M{
		"user": userId,
	}).Iter()

	devc := &Device{}
	for cursor.Next(devc) {
		devices = append(devices, devc)
		devc = &Device{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAllMode(db *database.Database, userId bson.ObjectId,
	mode string) (devices []*Device, err error) {

	coll := db.Devices()
	devices = []*Device{}

	cursor := coll.Find(&bson.M{
		"user": userId,
		"mode": mode,
	}).Iter()

	devc := &Device{}
	for cursor.Next(devc) {
		devices = append(devices, devc)
		devc = &Device{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Count(db *database.Database, userId bson.ObjectId) (
	count int, err error) {

	coll := db.Devices()

	count, err = coll.Find(&bson.M{
		"user": userId,
		"mode": Secondary,
	}).Count()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func New(userId bson.ObjectId, typ, mode string) (devc *Device) {
	devc = &Device{
		Id:         bson.NewObjectId(),
		Type:       typ,
		Mode:       mode,
		User:       userId,
		Timestamp:  time.Now(),
		LastActive: time.Now(),
	}

	return
}

func Remove(db *database.Database, id bson.ObjectId) (err error) {
	coll := db.Devices()

	err = coll.RemoveId(id)
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

func RemoveUser(db *database.Database, id bson.ObjectId,
	userId bson.ObjectId) (err error) {

	coll := db.Devices()

	err = coll.Remove(&bson.M{
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

func RemoveAll(db *database.Database, userId bson.ObjectId) (err error) {
	coll := db.Devices()

	_, err = coll.RemoveAll(&bson.M{
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
