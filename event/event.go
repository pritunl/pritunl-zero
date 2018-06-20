// Pub/sub messaging system using mongodb tailable cursor.
package event

import (
	"github.com/Sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/requires"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	listeners = map[string][]func(*Event){}
)

type Event struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Channel   string        `bson:"channel" json:"channel"`
	Timestamp time.Time     `bson:"timestamp" json:"timestamp"`
	Data      interface{}   `bson:"data" json:"data"`
}

type Dispatch struct {
	Type string `bson:"type" json:"type"`
}

func getCursorId(coll *database.Collection, channels []string) (
	id bson.ObjectId, err error) {

	msg := &Event{}

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
		err = coll.Find(query).Sort("-$natural").One(msg)

		if err != nil {
			err = database.ParseError(err)
			if i > 0 {
				return
			}

			switch err.(type) {
			case *database.NotFoundError:
				// Cannot use client-side ObjectId for tailable collection
				err = Publish(coll.Database, channels[0], nil)
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

func getCursorIdRetry(channels []string) bson.ObjectId {
	db := database.GetDatabase()
	defer db.Close()

	for {
		coll := db.Events()

		cursorId, err := getCursorId(coll, channels)
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

	msg := &Event{
		Id:        bson.NewObjectId(),
		Channel:   channel,
		Timestamp: time.Now(),
		Data:      data,
	}

	err = coll.Insert(msg)
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
	onMsg func(*Event, error) bool) {

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

	query := bson.M{
		"_id": bson.M{
			"$gt": cursorId,
		},
		"channel": channelBson,
	}
	iter := coll.Find(query).Sort("$natural").Tail(duration)
	defer func() {
		iter.Close()
	}()

	for {
		msg := &Event{}
		for iter.Next(msg) {
			cursorId = msg.Id

			if msg.Data == nil {
				// Blank msg for cursor
				continue
			}

			if !onMsg(msg, nil) {
				return
			}
		}

		if iter.Err() != nil {
			err := iter.Close()

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("event: Subscribe error")

			if !onMsg(nil, err) {
				return
			}

			time.Sleep(constants.RetryDelay)
		} else if iter.Timeout() {
			if !onMsg(nil, nil) {
				return
			}
			continue
		}

		iter.Close()
		db.Close()
		db = database.GetDatabase()
		coll = db.Events()

		query := &bson.M{
			"_id": &bson.M{
				"$gt": cursorId,
			},
			"channel": channelBson,
		}
		iter = coll.Find(query).Sort("$natural").Tail(duration)
	}
}

func Register(channel string, callback func(*Event)) {
	callbacks := listeners[channel]

	if callbacks == nil {
		callbacks = []func(*Event){}
	}

	listeners[channel] = append(callbacks, callback)
}

func subscribe(channels []string) {
	Subscribe(channels, 10*time.Second,
		func(msg *Event, err error) bool {
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
