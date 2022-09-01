package check

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/utils"
)

func Get(db *database.Database, checkId primitive.ObjectID) (
	chck *Check, err error) {

	coll := db.Checks()
	chck = &Check{}

	err = coll.FindOneId(checkId, chck)
	if err != nil {
		return
	}

	return
}

func GetMulti(db *database.Database, checkIds []primitive.ObjectID) (
	checks []*Check, err error) {

	coll := db.Checks()
	checks = []*Check{}

	cursor, err := coll.Find(
		db,
		&bson.M{
			"_id": &bson.M{
				"$in": checkIds,
			},
		},
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		chck := &Check{}
		err = cursor.Decode(chck)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		checks = append(checks, chck)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAll(db *database.Database) (checks []*Check, err error) {
	coll := db.Checks()
	checks = []*Check{}

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
		chck := &Check{}
		err = cursor.Decode(chck)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		checks = append(checks, chck)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAllPaged(db *database.Database, query *bson.M,
	page, pageCount int64) (checks []*Check, count int64, err error) {

	coll := db.Checks()
	checks = []*Check{}

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
	page = utils.Min64(page, maxPage)
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
		chck := &Check{}
		err = cursor.Decode(chck)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		checks = append(checks, chck)
		chck = &Check{}
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetRoles(db *database.Database, roles []string) (
	checks []*Check, err error) {

	coll := db.Checks()
	checks = []*Check{}

	if roles == nil {
		roles = []string{}
	}

	cursor, err := coll.Find(
		db,
		&bson.M{
			"roles": &bson.M{
				"$in": roles,
			},
		},
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		polcy := &Check{}
		err = cursor.Decode(polcy)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		checks = append(checks, polcy)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetRolesMapped(db *database.Database, rolesSet set.Set) (
	checksMap map[string][]*Check, err error) {

	checksMap = map[string][]*Check{}

	roles := []string{}
	for role := range rolesSet.Iter() {
		roles = append(roles, role.(string))
	}

	checks, err := GetRoles(db, roles)
	if err != nil {
		return
	}

	for _, chck := range checks {
		for _, role := range chck.Roles {
			roleAlrts := checksMap[role]
			if roleAlrts == nil {
				roleAlrts = []*Check{}
			}
			checksMap[role] = append(roleAlrts, chck)
		}
	}

	return
}

func Remove(db *database.Database,
	checkId primitive.ObjectID) (err error) {

	coll := db.Checks()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": checkId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func RemoveMulti(db *database.Database, checkIds []primitive.ObjectID) (
	err error) {

	coll := db.Checks()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": &bson.M{
			"$in": checkIds,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
