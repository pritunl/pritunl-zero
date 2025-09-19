package endpoints

import (
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
)

type endpointName struct {
	Id   bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string        `bson:"name" json:"name"`
}

func getRolesName(db *database.Database, roles []string) (
	endpts []*endpointName, err error) {

	coll := db.Endpoints()
	endpts = []*endpointName{}

	if roles == nil {
		roles = []string{}
	}

	cursor, err := coll.Find(
		db,
		bson.M{
			"roles": bson.M{
				"$in": roles,
			},
		},
		options.Find().
			SetProjection(bson.D{{"name", 1}}),
	)
	defer cursor.Close(db)

	for cursor.Next(db) {
		endpt := &endpointName{}
		err = cursor.Decode(endpt)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		endpts = append(endpts, endpt)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func getRolesNameMapped(db *database.Database, roles []string) (
	endptsMap map[bson.ObjectID]string, err error) {

	endptsMap = map[bson.ObjectID]string{}

	endpts, err := getRolesName(db, roles)
	if err != nil {
		return
	}

	for _, endpt := range endpts {
		endptsMap[endpt.Id] = endpt.Name
	}

	return
}
