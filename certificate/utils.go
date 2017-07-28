package certificate

import (
	"github.com/pritunl/pritunl-zero/database"
	"gopkg.in/mgo.v2/bson"
)

func Get(db *database.Database, certId bson.ObjectId) (
	cert *Certificate, err error) {

	coll := db.Certificates()
	cert = &Certificate{}

	err = coll.FindOneId(certId, cert)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database) (certs []*Certificate, err error) {
	coll := db.Certificates()
	certs = []*Certificate{}

	cursor := coll.Find(bson.M{}).Iter()

	cert := &Certificate{}
	for cursor.Next(cert) {
		certs = append(certs, cert)
		cert = &Certificate{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, certId bson.ObjectId) (err error) {
	coll := db.Certificates()

	_, err = coll.RemoveAll(&bson.M{
		"_id": certId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
