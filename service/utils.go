package service

import (
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
)

func Get(db *database.Database, serviceId primitive.ObjectID) (
	srvce *Service, err error) {

	coll := db.Services()
	srvce = &Service{}

	err = coll.FindOneId(serviceId, srvce)
	if err != nil {
		return
	}

	return
}

func GetMulti(db *database.Database, serviceIds []primitive.ObjectID) (
	services []*Service, err error) {

	coll := db.Services()
	services = []*Service{}

	cursor, err := coll.Find(
		db,
		&bson.M{
			"_id": &bson.M{
				"$in": serviceIds,
			},
		},
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		srvce := &Service{}
		err = cursor.Decode(srvce)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		services = append(services, srvce)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAll(db *database.Database) (services []*Service, err error) {
	coll := db.Services()
	services = []*Service{}

	cursor, err := coll.Find(
		db,
		&bson.M{},
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		srvce := &Service{}
		err = cursor.Decode(srvce)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		services = append(services, srvce)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, serviceId primitive.ObjectID) (err error) {
	coll := db.Services()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": serviceId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
