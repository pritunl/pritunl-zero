package policy

import (
	"github.com/pritunl/pritunl-zero/database"
	"gopkg.in/mgo.v2/bson"
)

func Get(db *database.Database, certId bson.ObjectId) (
	cert *Policy, err error) {

	coll := db.Policies()
	cert = &Policy{}

	err = coll.FindOneId(certId, cert)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database) (certs []*Policy, err error) {
	coll := db.Policies()
	certs = []*Policy{}

	cursor := coll.Find(bson.M{}).Iter()

	cert := &Policy{}
	for cursor.Next(cert) {
		certs = append(certs, cert)
		cert = &Policy{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, certId bson.ObjectId) (err error) {
	coll := db.Policies()

	_, err = coll.RemoveAll(&bson.M{
		"_id": certId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
