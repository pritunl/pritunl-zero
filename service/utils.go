package service

import (
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/utils"
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

func GetAllName(db *database.Database) (services []*Service, err error) {
	coll := db.Services()
	services = []*Service{}

	cursor, err := coll.Find(
		db,
		&bson.M{},
		&options.FindOptions{
			Projection: &bson.D{
				{"name", 1},
			},
		},
	)
	defer cursor.Close(db)

	for cursor.Next(db) {
		srvc := &Service{}
		err = cursor.Decode(srvc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		services = append(services, srvc)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAllPaged(db *database.Database, query *bson.M,
	page, pageCount int64) (services []*Service, count int64, err error) {

	coll := db.Services()
	services = []*Service{}

	count, err = coll.CountDocuments(db, query)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	page = utils.Min64(page, count/pageCount)
	skip := utils.Min64(page*pageCount, count)

	cursor, err := coll.Find(
		db,
		query,
		&options.FindOptions{
			Sort: &bson.D{
				{"name", 1},
			},
			Skip:  &skip,
			Limit: &pageCount,
		},
	)
	defer cursor.Close(db)

	for cursor.Next(db) {
		srvce := &Service{}
		err = cursor.Decode(srvce)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		services = append(services, srvce)
		srvce = &Service{}
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

func RemoveMulti(db *database.Database, serviceIds []primitive.ObjectID) (
	err error) {

	coll := db.Services()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": &bson.M{
			"$in": serviceIds,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
