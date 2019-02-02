// Pub/sub messaging system using mongodb tailable cursor.
package event

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/requires"
)

var (
	listeners = map[string][]func(*EventPublish){}
)

type Event struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Channel   string             `bson:"channel" json:"channel"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Data      bson.M             `bson:"data" json:"data"`
}

type EventPublish struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Channel   string             `bson:"channel" json:"channel"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Data      interface{}        `bson:"data" json:"data"`
}

type CustomEvent interface {
	GetId() primitive.ObjectID
	GetData() interface{}
}

type Dispatch struct {
	Type string `bson:"type" json:"type"`
}

func getCursorId(db *database.Database, coll *database.Collection,
	channels []string) (id primitive.ObjectID, err error) {

	msg := &EventPublish{}

	var query *bson.M
	if len(channels) == 1 {
		query = &bson.M{
			"channel": channels[0],
		}
	} else {
		query = &bson.M{
			"channel": &bson.M{
				"$in": channels,
			},
		}
	}

	for i := 0; i < 2; i++ {
		err = coll.FindOne(
			db,
			query,
			&options.FindOneOptions{
				Sort: &bson.D{
					{"$natural", -1},
				},
			},
		).Decode(msg)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		if err != nil {
			err = database.ParseError(err)
			if i > 0 {
				return
			}

			switch err.(type) {
			case *database.NotFoundError:
				// Cannot use client-side ObjectId for tailable collection
				err = Publish(db, channels[0], nil)
				if err != nil {
					err = database.ParseError(err)
					return
				}
				continue
			default:
				return
			}
		} else {
			break
		}
	}

	id = msg.Id

	return
}

func getCursorIdRetry(channels []string) primitive.ObjectID {
	db := database.GetDatabase()
	defer db.Close()

	for {
		coll := db.Events()

		cursorId, err := getCursorId(db, coll, channels)
		if err != nil {
			err = database.ParseError(err)

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("event: Subscribe cursor error")

			db.Close()
			db = database.GetDatabase()

			time.Sleep(constants.RetryDelay)

			continue
		}

		return cursorId
	}
}

func Publish(db *database.Database, channel string, data interface{}) (
	err error) {

	coll := db.Events()

	msg := &EventPublish{
		Id:        primitive.NewObjectID(),
		Channel:   channel,
		Timestamp: time.Now(),
		Data:      data,
	}

	_, err = coll.InsertOne(db, msg)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func PublishDispatch(db *database.Database, typ string) (
	err error) {

	evt := &Dispatch{
		Type: typ,
	}

	err = Publish(db, "dispatch", evt)
	if err != nil {
		return
	}

	return
}

func Subscribe(channels []string, duration time.Duration,
	onMsg func(*EventPublish, error) bool) {

	db := database.GetDatabase()
	defer db.Close()
	coll := db.Events()

	cursorId := getCursorIdRetry(channels)

	var channelBson interface{}
	if len(channels) == 1 {
		channelBson = channels[0]
	} else {
		channelBson = &bson.M{
			"$in": channels,
		}
	}

	queryOpts := &options.FindOptions{
		Sort: &bson.D{
			{"$natural", 1},
		},
	}
	queryOpts.SetMaxAwaitTime(duration)
	queryOpts.SetCursorType(options.TailableAwait)

	query := &bson.M{
		"_id": &bson.M{
			"$gt": cursorId,
		},
		"channel": channelBson,
	}

	var cursor mongo.Cursor
	var err error
	for {
		cursor, err = coll.Find(
			db,
			query,
			queryOpts,
		)
		if err != nil {
			err = database.ParseError(err)

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("event: Listener find error")

			if !onMsg(nil, err) {
				return
			}
		} else {
			break
		}

		time.Sleep(constants.RetryDelay)
	}
	defer func() {
		defer func() {
			recover()
		}()
		cursor.Close(db)
	}()

	for {
		for cursor.Next(db) {
			msg := &EventPublish{}
			err = cursor.Decode(msg)
			if err != nil {
				err = database.ParseError(err)

				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("event: Listener decode error")

				if !onMsg(nil, err) {
					return
				}

				time.Sleep(constants.RetryDelay)
				break
			}

			cursorId = msg.Id

			if msg.Data == nil {
				// Blank msg for cursor
				continue
			}

			if !onMsg(msg, nil) {
				return
			}
		}

		err = cursor.Err()
		if err != nil {
			err = database.ParseError(err)

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("event: Listener cursor error")

			if !onMsg(nil, err) {
				return
			}

			time.Sleep(constants.RetryDelay)
		}

		cursor.Close(db)
		db.Close()
		db = database.GetDatabase()
		coll = db.Events()

		query := &bson.M{
			"_id": &bson.M{
				"$gt": cursorId,
			},
			"channel": channelBson,
		}
		for {
			cursor, err = coll.Find(
				db,
				query,
				queryOpts,
			)
			if err != nil {
				err = database.ParseError(err)

				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("event: Listener find error")

				if !onMsg(nil, err) {
					return
				}
			} else {
				break
			}

			time.Sleep(constants.RetryDelay)
		}
	}
}

func SubscribeType(channels []string, duration time.Duration,
	newEvent func() CustomEvent, onMsg func(CustomEvent, error) bool) {

	db := database.GetDatabase()
	defer db.Close()
	coll := db.Events()

	cursorId := getCursorIdRetry(channels)

	var channelBson interface{}
	if len(channels) == 1 {
		channelBson = channels[0]
	} else {
		channelBson = &bson.M{
			"$in": channels,
		}
	}

	queryOpts := &options.FindOptions{
		Sort: &bson.D{
			{"$natural", 1},
		},
	}
	queryOpts.SetMaxAwaitTime(duration)
	queryOpts.SetCursorType(options.TailableAwait)

	query := &bson.M{
		"_id": &bson.M{
			"$gt": cursorId,
		},
		"channel": channelBson,
	}

	var cursor mongo.Cursor
	var err error
	for {
		cursor, err = coll.Find(
			db,
			query,
			queryOpts,
		)
		if err != nil {
			err = database.ParseError(err)

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("event: Listener find error")

			if !onMsg(nil, err) {
				return
			}
		} else {
			break
		}

		time.Sleep(constants.RetryDelay)
	}
	defer func() {
		defer func() {
			recover()
		}()
		cursor.Close(db)
	}()

	for {
		for cursor.Next(db) {
			msg := newEvent()
			err = cursor.Decode(msg)
			if err != nil {
				err = database.ParseError(err)

				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("event: Listener decode error")

				if !onMsg(nil, err) {
					return
				}

				time.Sleep(constants.RetryDelay)
				break
			}

			cursorId = msg.GetId()

			if msg.GetData() == nil {
				// Blank msg for cursor
				continue
			}

			if !onMsg(msg, nil) {
				return
			}
		}

		err = cursor.Err()
		if err != nil {
			err = database.ParseError(err)

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("event: Listener cursor error")

			if !onMsg(nil, err) {
				return
			}

			time.Sleep(constants.RetryDelay)
		}

		cursor.Close(db)
		db.Close()
		db = database.GetDatabase()
		coll = db.Events()

		query := &bson.M{
			"_id": &bson.M{
				"$gt": cursorId,
			},
			"channel": channelBson,
		}
		for {
			cursor, err = coll.Find(
				db,
				query,
				queryOpts,
			)
			if err != nil {
				err = database.ParseError(err)

				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("event: Listener find error")

				if !onMsg(nil, err) {
					return
				}
			} else {
				break
			}

			time.Sleep(constants.RetryDelay)
		}
	}
}

func Register(channel string, callback func(*EventPublish)) {
	callbacks := listeners[channel]

	if callbacks == nil {
		callbacks = []func(*EventPublish){}
	}

	listeners[channel] = append(callbacks, callback)
}

func subscribe(channels []string) {
	Subscribe(channels, 10*time.Second,
		func(msg *EventPublish, err error) bool {
			if msg == nil || err != nil {
				return true
			}

			for _, listener := range listeners[msg.Channel] {
				listener(msg)
			}

			return true
		})
}

func init() {
	module := requires.New("event")
	module.After("settings")

	module.Handler = func() (err error) {
		go func() {
			channels := []string{}

			for channel := range listeners {
				channels = append(channels, channel)
			}

			if len(channels) > 0 {
				subscribe(channels)
			}
		}()

		return
	}
}
