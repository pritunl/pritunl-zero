package event

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
)

type Listener struct {
	db       *database.Database
	state    bool
	err      error
	channels []string
	stream   chan *Event
}

func (l *Listener) Listen() chan *Event {
	return l.stream
}

func (l *Listener) Close() {
	l.state = false
	close(l.stream)
}

func (l *Listener) sub(cursorId primitive.ObjectID) {
	coll := l.db.Events()

	var channelBson interface{}
	if len(l.channels) == 1 {
		channelBson = l.channels[0]
	} else {
		channelBson = &bson.M{
			"$in": l.channels,
		}
	}

	queryOpts := &options.FindOptions{
		Sort: &bson.D{
			{"$natural", 1},
		},
	}
	queryOpts.SetMaxAwaitTime(10 * time.Second)
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
			l.db,
			query,
			queryOpts,
		)
		if err != nil {
			err = database.ParseError(err)

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("event: Listener find error")
		} else {
			break
		}

		if !l.state {
			return
		}

		time.Sleep(constants.RetryDelay)

		if !l.state {
			return
		}
	}

	defer func() {
		defer func() {
			recover()
		}()
		cursor.Close(l.db)
	}()

	for {
		for cursor.Next(l.db) {
			msg := &Event{}
			err = cursor.Decode(msg)
			if err != nil {
				err = database.ParseError(err)

				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("event: Listener decode error")

				time.Sleep(constants.RetryDelay)
				break
			}

			cursorId = msg.Id

			if msg.Data == nil {
				// Blank msg for cursor
				continue
			}

			if !l.state {
				return
			}

			l.stream <- msg
		}

		if !l.state {
			return
		}

		err = cursor.Err()
		if err != nil {
			err = database.ParseError(err)

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("event: Listener cursor error")

			time.Sleep(constants.RetryDelay)
		}

		if !l.state {
			return
		}

		cursor.Close(l.db)
		coll = l.db.Events()

		query := &bson.M{
			"_id": &bson.M{
				"$gt": cursorId,
			},
			"channel": channelBson,
		}

		for {
			cursor, err = coll.Find(
				l.db,
				query,
				queryOpts,
			)
			if err != nil {
				err = database.ParseError(err)

				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("event: Listener find error")
			} else {
				break
			}

			if !l.state {
				return
			}

			time.Sleep(constants.RetryDelay)

			if !l.state {
				return
			}
		}
	}
}

func (l *Listener) init() (err error) {
	coll := l.db.Events()
	cursorId, err := getCursorId(l.db, coll, l.channels)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	l.state = true

	go func() {
		if r := recover(); r != nil {
			logrus.WithFields(logrus.Fields{
				"error": errors.New(fmt.Sprintf("%s", r)),
			}).Error("event: Listener panic")
		}
		l.sub(cursorId)
	}()

	return
}

func SubscribeListener(db *database.Database, channels []string) (
	lst *Listener, err error) {

	lst = &Listener{
		db:       db,
		channels: channels,
		stream:   make(chan *Event, 10),
	}

	err = lst.init()
	if err != nil {
		return
	}

	return
}
