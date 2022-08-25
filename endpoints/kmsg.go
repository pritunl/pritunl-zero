package endpoints

import (
	"context"
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/alert"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/settings"
)

type Kmsg struct {
	Id        primitive.Binary   `bson:"_id" json:"id"`
	Endpoint  primitive.ObjectID `bson:"e" json:"e"`
	Timestamp time.Time          `bson:"t" json:"t"`

	Boot     int64  `bson:"b" json:"b"`
	Priortiy int    `bson:"p" json:"p"`
	Sequence int64  `bson:"s" json:"s"`
	Message  string `bson:"m" json:"m"`
}

func (d *Kmsg) generateId() primitive.Binary {
	hash := md5.New()
	hash.Write([]byte(d.Endpoint.Hex()))
	hash.Write([]byte(strconv.FormatInt(d.Boot, 10)))
	hash.Write([]byte("-"))
	hash.Write([]byte(strconv.FormatInt(d.Sequence, 10)))

	return primitive.Binary{
		Subtype: BinaryMD5,
		Data:    hash.Sum(nil),
	}
}

func (d *Kmsg) GetCollection(db *database.Database) *database.Collection {
	return db.EndpointsKmsg()
}

func (d *Kmsg) Format(id primitive.ObjectID) time.Time {
	d.Endpoint = id
	d.Id = d.generateId()
	return d.Timestamp
}

func (d *Kmsg) StaticData() *bson.M {
	return nil
}

func (d *Kmsg) CheckAlerts(resources []*alert.Alert) (alerts []*Alert) {
	alerts = []*Alert{}

	for _, resource := range resources {
		switch resource.Resource {
		case alert.KmsgKeyword:
			if strings.Contains(strings.ToLower(d.Message),
				strings.ToLower(resource.ValueStr)) {

				alerts = []*Alert{
					NewAlert(resource, fmt.Sprintf(
						"Kmsg keyword match (%s): %s",
						resource.ValueStr,
						strings.Split(d.Message, "\n")[0],
					)),
				}
			}
			break
		}
	}

	return
}

func (d *Kmsg) FormattedLog() string {
	return fmt.Sprintf(
		"[%s] %s",
		d.Timestamp.Format("Mon Jan _2 15:04:05 2006"),
		d.Message,
	)
}

func GetKmsgLog(c context.Context, db *database.Database,
	endpoint primitive.ObjectID) (logData LogData, err error) {

	logData = []string{}

	coll := db.EndpointsKmsg()

	cursor, err := coll.Find(
		c,
		&bson.M{
			"e": endpoint,
		},
		&options.FindOptions{
			Limit: &settings.Endpoint.KmsgDisplayLimit,
			Sort: &bson.D{
				{"b", 1},
				{"s", 1},
			},
		},
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		doc := &Kmsg{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		logData = append(logData, doc.FormattedLog())
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
