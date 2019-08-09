package rokey

import (
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/utils"
)

func Get(db *database.Database, typ string) (rkey *Rokey, err error) {
	timestamp := time.Now()
	timeblock := time.Date(
		timestamp.Year(),
		timestamp.Month(),
		timestamp.Day(),
		timestamp.Hour(),
		0,
		0,
		0,
		timestamp.Location(),
	)

	rkey = GetCache(typ, timeblock)
	if rkey != nil {
		return
	}

	secret, err := utils.RandStr(64)
	if err != nil {
		return
	}

	coll := db.Rokeys()
	rkey = &Rokey{
		Type:      typ,
		Timeblock: timeblock,
		Timestamp: timestamp,
		Secret:    secret,
	}

	opts := &options.FindOneAndUpdateOptions{}
	opts.SetUpsert(true)
	opts.SetReturnDocument(options.After)

	err = coll.FindOneAndUpdate(
		db,
		&bson.M{
			"type":      typ,
			"timeblock": timeblock,
		},
		&bson.M{
			"$setOnInsert": rkey,
		},
		opts,
	).Decode(rkey)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	PutCache(rkey)

	return
}

func GetId(db *database.Database, typ string,
	rkeyId primitive.ObjectID) (rkey *Rokey, err error) {

	rkey = GetCacheId(typ, rkeyId)
	if rkey != nil {
		return
	}

	coll := db.Rokeys()
	rkey = &Rokey{}

	err = coll.FindOneId(rkeyId, rkey)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			rkey = nil
			err = nil
		} else {
			return
		}
	}

	if rkey != nil {
		PutCache(rkey)
		if rkey.Type != typ {
			rkey = nil
		}
	}

	return
}
