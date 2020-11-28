package endpoint

import (
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/utils"
)

func Get(db *database.Database, endpointId primitive.ObjectID) (
	endpt *Endpoint, err error) {

	coll := db.Endpoints()
	endpt = &Endpoint{}

	err = coll.FindOneId(endpointId, endpt)
	if err != nil {
		return
	}

	return
}

func GetMulti(db *database.Database, endpointIds []primitive.ObjectID) (
	endpoints []*Endpoint, err error) {

	coll := db.Endpoints()
	endpoints = []*Endpoint{}

	cursor, err := coll.Find(
		db,
		&bson.M{
			"_id": &bson.M{
				"$in": endpointIds,
			},
		},
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		endpt := &Endpoint{}
		err = cursor.Decode(endpt)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		endpoints = append(endpoints, endpt)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAll(db *database.Database) (endpoints []*Endpoint, err error) {
	coll := db.Endpoints()
	endpoints = []*Endpoint{}

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
		endpt := &Endpoint{}
		err = cursor.Decode(endpt)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		endpoints = append(endpoints, endpt)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAllPaged(db *database.Database, query *bson.M,
	page, pageCount int64) (endpoints []*Endpoint, count int64, err error) {

	coll := db.Endpoints()
	endpoints = []*Endpoint{}

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
		endpt := &Endpoint{}
		err = cursor.Decode(endpt)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		endpoints = append(endpoints, endpt)
		endpt = &Endpoint{}
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database,
	endpointId primitive.ObjectID) (err error) {

	coll := db.Endpoints()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": endpointId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func RemoveMulti(db *database.Database, endpointIds []primitive.ObjectID) (
	err error) {

	coll := db.Endpoints()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": &bson.M{
			"$in": endpointIds,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
