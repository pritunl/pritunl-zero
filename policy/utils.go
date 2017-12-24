package policy

import (
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/user"
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

func GetService(db *database.Database, serviceId bson.ObjectId) (
	policies []*Policy, err error) {

	coll := db.Policies()
	policies = []*Policy{}

	cursor := coll.Find(bson.M{
		"services": serviceId,
	}).Iter()

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

func GetRoles(db *database.Database, roles []string) (
	policies []*Policy, err error) {

	coll := db.Policies()
	policies = []*Policy{}

	cursor := coll.Find(bson.M{
		"roles": &bson.M{
			"$in": roles,
		},
	}).Iter()

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

func KeybaseMode(policies []*Policy) (mode string) {
	mode = Optional

	for _, polcy := range policies {
		switch polcy.KeybaseMode {
		case Disabled:
			if mode == Optional {
				mode = Disabled
			}
			break
		case Required:
			mode = Required
			break
		}
	}

	return
}
