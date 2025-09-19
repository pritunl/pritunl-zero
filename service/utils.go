package service

import (
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
)

func Get(db *database.Database, serviceId bson.ObjectID) (
	srvce *Service, err error) {

	coll := db.Services()
	srvce = &Service{}

	err = coll.FindOneId(serviceId, srvce)
	if err != nil {
		return
	}

	return
}

func GetOne(db *database.Database, query *bson.M) (
	srvce *Service, err error) {

	coll := db.Services()
	srvce = &Service{}

	err = coll.FindOne(db, query).Decode(srvce)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetMulti(db *database.Database, serviceIds []bson.ObjectID) (
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

func GetAllNames(db *database.Database) (
	services []*database.Named, err error) {

	coll := db.Services()
	services = []*database.Named{}

	cursor, err := coll.Find(
		db,
		&bson.M{},
		options.Find().
			SetSort(bson.D{{"name", 1}}).
			SetProjection(bson.D{{"name", 1}}),
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		srvc := &database.Named{}
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

	if len(*query) == 0 {
		count, err = coll.EstimatedDocumentCount(db)
		if err != nil {
			err = database.ParseError(err)
			return
		}
	} else {
		count, err = coll.CountDocuments(db, query)
		if err != nil {
			err = database.ParseError(err)
			return
		}
	}

	maxPage := count / pageCount
	if count == pageCount {
		maxPage = 0
	}
	page = min(page, maxPage)
	skip := min(page*pageCount, count)

	cursor, err := coll.Find(
		db,
		query,
		options.Find().
			SetSort(bson.D{{"name", 1}}).
			SetSkip(skip).
			SetLimit(pageCount),
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

func Remove(db *database.Database, serviceId bson.ObjectID) (err error) {
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

func RemoveMulti(db *database.Database, serviceIds []bson.ObjectID) (
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
