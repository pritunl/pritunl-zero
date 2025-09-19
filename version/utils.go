package version

import (
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
)

func Check(db *database.Database, module string, ver int) (
	supported bool, err error) {

	if !cacheCheck(module, ver) {
		return false, nil
	}

	coll := db.Versions()
	vr := &Version{}

	err = coll.FindOneId(module, vr)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			vr = nil
			err = nil
		} else {
			return
		}
	}

	if vr == nil || ver >= vr.Version {
		supported = true
		return
	}

	cacheSet(module, vr.Version)
	return
}

func Set(db *database.Database, module string, ver int) (err error) {
	coll := db.Versions()

	_, err = coll.UpdateOne(
		db,
		&bson.M{
			"_id": module,
		},
		&bson.M{
			"$max": &bson.M{
				"version": ver,
			},
			"$setOnInsert": &bson.M{
				"_id": module,
			},
		},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
