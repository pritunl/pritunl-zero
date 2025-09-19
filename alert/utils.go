package alert

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
)

func Get(db *database.Database, alertId bson.ObjectID) (
	alrt *Alert, err error) {

	coll := db.Alerts()
	alrt = &Alert{}

	err = coll.FindOneId(alertId, alrt)
	if err != nil {
		return
	}

	return
}

func GetMulti(db *database.Database, alertIds []bson.ObjectID) (
	alerts []*Alert, err error) {

	coll := db.Alerts()
	alerts = []*Alert{}

	cursor, err := coll.Find(
		db,
		&bson.M{
			"_id": &bson.M{
				"$in": alertIds,
			},
		},
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		alrt := &Alert{}
		err = cursor.Decode(alrt)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		alerts = append(alerts, alrt)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAll(db *database.Database) (alerts []*Alert, err error) {
	coll := db.Alerts()
	alerts = []*Alert{}

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
		alrt := &Alert{}
		err = cursor.Decode(alrt)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		alerts = append(alerts, alrt)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAllPaged(db *database.Database, query *bson.M,
	page, pageCount int64) (alerts []*Alert, count int64, err error) {

	coll := db.Alerts()
	alerts = []*Alert{}

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
		alrt := &Alert{}
		err = cursor.Decode(alrt)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		alerts = append(alerts, alrt)
		alrt = &Alert{}
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetRoles(db *database.Database, roles []string) (
	alerts []*Alert, err error) {

	coll := db.Alerts()
	alerts = []*Alert{}

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
		polcy := &Alert{}
		err = cursor.Decode(polcy)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		alerts = append(alerts, polcy)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetRolesMapped(db *database.Database, rolesSet set.Set) (
	alertsMap map[string][]*Alert, err error) {

	alertsMap = map[string][]*Alert{}

	roles := []string{}
	for role := range rolesSet.Iter() {
		roles = append(roles, role.(string))
	}

	alerts, err := GetRoles(db, roles)
	if err != nil {
		return
	}

	for _, alrt := range alerts {
		for _, role := range alrt.Roles {
			roleAlrts := alertsMap[role]
			if roleAlrts == nil {
				roleAlrts = []*Alert{}
			}
			alertsMap[role] = append(roleAlrts, alrt)
		}
	}

	return
}

func Remove(db *database.Database,
	alertId bson.ObjectID) (err error) {

	coll := db.Alerts()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": alertId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func RemoveMulti(db *database.Database, alertIds []bson.ObjectID) (
	err error) {

	coll := db.Alerts()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": &bson.M{
			"$in": alertIds,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
