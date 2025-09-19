package endpoint

import (
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
)

func Get(db *database.Database, endpointId bson.ObjectID) (
	endpt *Endpoint, err error) {

	coll := db.Endpoints()
	endpt = &Endpoint{}

	err = coll.FindOneId(endpointId, endpt)
	if err != nil {
		return
	}

	return
}

func GetMulti(db *database.Database, endpointIds []bson.ObjectID) (
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

func RemoveData(db *database.Database, endpointId bson.ObjectID) (
	err error) {

	coll := db.EndpointsSystem()
	_, err = coll.DeleteMany(db, &bson.M{
		"e": endpointId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	coll = db.EndpointsLoad()
	_, err = coll.DeleteMany(db, &bson.M{
		"e": endpointId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	coll = db.EndpointsDisk()
	_, err = coll.DeleteMany(db, &bson.M{
		"e": endpointId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	coll = db.EndpointsDiskIo()
	_, err = coll.DeleteMany(db, &bson.M{
		"e": endpointId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	coll = db.EndpointsNetwork()
	_, err = coll.DeleteMany(db, &bson.M{
		"e": endpointId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	coll = db.EndpointsCheck()
	_, err = coll.DeleteMany(db, &bson.M{
		"e": endpointId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	coll = db.EndpointsCheckLog()
	_, err = coll.DeleteMany(db, &bson.M{
		"e": endpointId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database,
	endpointId bson.ObjectID) (err error) {

	err = RemoveData(db, endpointId)
	if err != nil {
		return
	}

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

func RemoveMulti(db *database.Database, endpointIds []bson.ObjectID) (
	err error) {

	coll := db.Endpoints()

	for _, endpointId := range endpointIds {
		err = RemoveData(db, endpointId)
		if err != nil {
			return
		}
	}

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
