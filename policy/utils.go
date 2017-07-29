package policy

import (
	"github.com/pritunl/pritunl-zero/database"
	"gopkg.in/mgo.v2/bson"
)

func Get(db *database.Database, policyId bson.ObjectId) (
	polcy *Policy, err error) {

	coll := db.Policies()
	polcy = &Policy{}

	err = coll.FindOneId(policyId, polcy)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database) (policies []*Policy, err error) {
	coll := db.Policies()
	policies = []*Policy{}

	cursor := coll.Find(bson.M{}).Iter()

	polcy := &Policy{}
	for cursor.Next(polcy) {
		policies = append(policies, polcy)
		polcy = &Policy{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, policyId bson.ObjectId) (err error) {
	coll := db.Policies()

	_, err = coll.RemoveAll(&bson.M{
		"_id": policyId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
