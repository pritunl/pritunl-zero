package service

import (
	"github.com/pritunl/pritunl-zero/database"
	"gopkg.in/mgo.v2/bson"
)

func Get(db *database.Database, serviceId bson.ObjectId) (
	srvce *Service, err error) {

	coll := db.Services()
	srvce = &Service{}

	err = coll.FindOneId(serviceId, srvce)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database) (services []*Service, err error) {

	coll := db.Services()
	services = []*Service{}

	cursor := coll.Find(bson.M{}).Iter()

	srvce := &Service{}
	for cursor.Next(srvce) {
		services = append(services, srvce)
		srvce = &Service{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, serviceId bson.ObjectId) (err error) {
	coll := db.Services()

	_, err = coll.RemoveAll(&bson.M{
		"_id": serviceId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
