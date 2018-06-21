package event

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Listener struct {
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

func (l *Listener) sub(db *database.Database, cursorId bson.ObjectId) {
	defer db.Close()
	coll := db.Events()

	var channelBson interface{}
	if len(l.channels) == 1 {
		channelBson = l.channels[0]
	} else {
		channelBson = &bson.M{
			"$in": l.channels,
		}
	}

	query := &bson.M{
		"_id": &bson.M{
			"$gt": cursorId,
		},
		"channel": channelBson,
	}
	iter := coll.Find(query).Sort("$natural").Tail(10 * time.Second)
	defer func() {
		defer func() {
			recover()
		}()
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

			if !l.state {
				return
			}

			l.stream <- msg
		}

		if !l.state {
			return
		}

		if iter.Err() != nil {
			err := iter.Close()

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("event: Listener error")

			time.Sleep(constants.RetryDelay)
		} else if iter.Timeout() {
			continue
		}

		if !l.state {
			return
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
		iter = coll.Find(query).Sort("$natural").Tail(10 * time.Second)
	}
}

func (l *Listener) init() (err error) {
	db := database.GetDatabase()

	coll := db.Events()
	cursorId, err := getCursorId(coll, l.channels)
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
		l.sub(db, cursorId)
	}()

	return
}

func SubscribeListener(channels []string) (lst *Listener, err error) {
	lst = &Listener{
		channels: channels,
		stream:   make(chan *Event),
	}

	err = lst.init()
	if err != nil {
		return
	}

	return
}
