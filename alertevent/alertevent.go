package alertevent

import (
	"fmt"
	"time"

	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/device"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/sirupsen/logrus"
)

type Alert struct {
	Id         string        `bson:"_id" json:"_id"`
	Name       string        `bson:"name" json:"name"`
	Timestamp  time.Time     `bson:"timestamp" json:"timestamp"`
	Roles      []string      `bson:"roles" json:"roles"`
	Source     bson.ObjectID `bson:"source" json:"source"`
	SourceName string        `bson:"source_name" json:"source_name"`
	Level      int           `bson:"level" json:"level"`
	Resource   string        `bson:"resource" json:"resource"`
	Message    string        `bson:"message" json:"message"`
	Frequency  time.Duration `bson:"frequency" json:"frequency"`
}

func (a *Alert) GetFrequency() (frequency time.Duration) {
	frequency = a.Frequency
	if frequency == 0 {
		frequency = 5 * time.Minute
	}
	return
}

func (a *Alert) DocId() string {
	timestamp := a.Timestamp.Unix()
	timekey := timestamp - (timestamp % int64(a.GetFrequency().Seconds()))

	return fmt.Sprintf(
		"%s-%s-%d",
		a.Source.Hex(),
		a.Resource,
		timekey,
	)
}

func (a *Alert) Key(devc *device.Device) string {
	timestamp := a.Timestamp.Unix()
	timekey := timestamp - (timestamp % int64(a.GetFrequency().Seconds()))

	return fmt.Sprintf(
		"%s-%s-%s-%d",
		a.Source.Hex(),
		a.Resource,
		devc.Id.Hex(),
		timekey,
	)
}

func (a *Alert) Lock(db *database.Database, devc *device.Device) (
	success bool, err error) {

	coll := db.AlertsEventLock()

	_, err = coll.InsertOne(db, &bson.M{
		"_id":       a.Key(devc),
		"timestamp": time.Now(),
	})
	if err != nil {
		err = database.ParseError(err)
		if _, ok := err.(*database.DuplicateKeyError); ok {
			err = nil
		}
		return
	}

	success = true

	return
}

func (a *Alert) FormattedTextMessage() string {
	return fmt.Sprintf("%s:%s == %s", a.Name, a.SourceName, a.Message)
}

func (a *Alert) FormattedCallMessage() string {
	return fmt.Sprintf("%s. %s", a.SourceName, a.Message)
}

func (a *Alert) Send(db *database.Database, roles []string) (err error) {
	coll := db.AlertsEvent()
	alrt := &Alert{}

	err = coll.FindOneId(a.Id, alrt)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			alrt = nil
			err = nil
		} else {
			return
		}
	}

	if alrt != nil && time.Since(alrt.Timestamp) < alrt.GetFrequency() {
		return
	}

	users, _, err := user.GetAll(db, &bson.M{
		"roles": &bson.D{
			{"$in", roles},
		},
	}, 0, 0)
	if err != nil {
		return
	}

	for _, usr := range users {
		devices, e := usr.GetDevices(db)
		if e != nil {
			err = e
			return
		}
		for _, devc := range devices {
			if devc.Mode != device.Phone || !devc.CheckLevel(a.Level) {
				continue
			}

			success, e := a.Lock(db, devc)
			if e != nil {
				err = e
				return
			}

			if !success {
				continue
			}

			msg := ""
			if devc.Type == device.Call {
				msg = a.FormattedCallMessage()
			} else {
				msg = a.FormattedTextMessage()
			}

			errData, e := Send(devc.Number, msg, devc.Type)
			if e != nil {
				if errData != nil {
					logrus.WithFields(logrus.Fields{
						"server_error":   errData.Error,
						"server_message": errData.Message,
						"error":          e,
					}).Error("alert: Failed to send alert")
				} else {
					logrus.WithFields(logrus.Fields{
						"error": e,
					}).Error("alert: Failed to send alert")
				}
			}
		}
	}

	_, err = coll.InsertOne(db, a)
	if err != nil {
		err = database.ParseError(err)
		if _, ok := err.(*database.DuplicateKeyError); ok {
			err = nil
		}
		return
	}

	return
}

func New(roles []string, source bson.ObjectID,
	name, sourceName, resource, message string, level int,
	frequency time.Duration) {

	db := database.GetDatabase()
	defer db.Close()

	alrt := &Alert{
		Name:       name,
		Timestamp:  time.Now(),
		Roles:      roles,
		Source:     source,
		SourceName: sourceName,
		Level:      level,
		Resource:   resource,
		Message:    message,
		Frequency:  frequency,
	}

	alrt.Id = alrt.DocId()

	err := alrt.Send(db, roles)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("alert: Failed to process alert")
	}

	return
}
