package session

import (
	"net/http"
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/useragent"
	"github.com/pritunl/pritunl-zero/utils"
)

func GetExpire(typ string) time.Duration {
	switch typ {
	case Proxy:
		return time.Duration(settings.Auth.ProxyExpire) * time.Minute
	case User:
		return time.Duration(settings.Auth.UserExpire) * time.Minute
	default:
		return time.Duration(settings.Auth.AdminExpire) * time.Minute
	}
}

func GetMaxDuration(typ string) time.Duration {
	switch typ {
	case Proxy:
		return time.Duration(settings.Auth.ProxyMaxDuration) * time.Minute
	case User:
		return time.Duration(settings.Auth.UserMaxDuration) * time.Minute
	default:
		return time.Duration(settings.Auth.AdminMaxDuration) * time.Minute
	}
}

func Get(db *database.Database, sessId string) (
	sess *Session, err error) {

	coll := db.Sessions()
	sess = &Session{}

	err = coll.FindOneId(sessId, sess)
	if err != nil {
		return
	}

	return
}

func GetUpdate(db *database.Database, sessId string, r *http.Request,
	typ, sig string) (sess *Session, err error) {

	query := bson.M{
		"_id": sessId,
		"removed": &bson.M{
			"$ne": true,
		},
	}

	expire := GetExpire(typ)
	maxDuration := GetMaxDuration(typ)

	if expire != 0 {
		query["last_active"] = &bson.M{
			"$gte": time.Now().Add(-expire),
		}
	}

	if maxDuration != 0 {
		query["timestamp"] = &bson.M{
			"$gte": time.Now().Add(-maxDuration),
		}
	}

	coll := db.Sessions()
	sess = &Session{}
	timestamp := time.Now()

	err = coll.FindOneAndUpdate(
		db,
		query,
		&bson.M{
			"$set": &bson.M{
				"last_active": timestamp,
			},
		},
	).Decode(sess)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	sess.LastActive = timestamp

	valid, err := sess.CheckSignature(db, sig)
	if err != nil {
		return
	}

	if !valid {
		sess = nil
		return
	}

	agnt, err := useragent.Parse(db, r)
	if err != nil {
		return
	}

	if agnt != nil && (sess.Agent == nil || sess.Agent.Diff(agnt)) {
		sess.Agent = agnt
		err = coll.UpdateId(sess.Id, &bson.M{
			"$set": &bson.M{
				"agent": agnt,
			},
		})
		if err != nil {
			err = database.ParseError(err)
			return
		}
	}

	return
}

func GetAll(db *database.Database, userId primitive.ObjectID,
	includeRemoved bool) (sessions []*Session, err error) {

	coll := db.Sessions()
	sessions = []*Session{}

	cursor, err := coll.Find(db, &bson.M{
		"user": userId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		sess := &Session{}
		err = cursor.Decode(sess)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		if !sess.Active() {
			if !includeRemoved {
				continue
			}
			sess.Removed = true
		}
		sessions = append(sessions, sess)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func New(db *database.Database, r *http.Request, userId primitive.ObjectID,
	typ string) (sess *Session, sig string, err error) {

	id, err := utils.RandStr(32)
	if err != nil {
		return
	}

	agnt, err := useragent.Parse(db, r)
	if err != nil {
		return
	}

	coll := db.Sessions()
	sess = &Session{
		Id:         id,
		Type:       typ,
		User:       userId,
		Timestamp:  time.Now(),
		LastActive: time.Now(),
		Agent:      agnt,
	}

	sig, err = sess.GenerateSignature(db)
	if err != nil {
		return
	}

	_, err = coll.InsertOne(db, sess)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, id string) (err error) {
	coll := db.Sessions()

	err = coll.UpdateId(id, &bson.M{
		"$set": &bson.M{
			"removed": true,
		},
	})
	if err != nil {
		err = database.ParseError(err)

		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		default:
			return
		}
	}

	return
}

func RemoveAll(db *database.Database, userId primitive.ObjectID) (err error) {
	coll := db.Sessions()

	_, err = coll.UpdateMany(db, &bson.M{
		"user": userId,
	}, &bson.M{
		"$set": &bson.M{
			"removed": true,
		},
	})
	if err != nil {
		err = database.ParseError(err)

		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		default:
			return
		}
	}

	return
}
