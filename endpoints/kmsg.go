package endpoints

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/pritunl/pritunl-zero/alert"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/settings"
)

type Kmsg struct {
	Id        bson.ObjectID `bson:"_id" json:"id"`
	Endpoint  bson.ObjectID `bson:"e" json:"e"`
	Timestamp time.Time     `bson:"t" json:"t"`

	Boot     int64  `bson:"b" json:"b"`
	Priortiy int    `bson:"p" json:"p"`
	Sequence int64  `bson:"s" json:"s"`
	Message  string `bson:"m" json:"m"`
}

func (d *Kmsg) generateId() bson.ObjectID {
	var b [12]byte

	hash := fnv.New64a()
	hash.Write(d.Endpoint[:])
	binary.Write(hash, binary.BigEndian, d.Sequence)
	sum := hash.Sum(nil)

	binary.BigEndian.PutUint32(b[0:4], uint32(d.Boot))
	copy(b[4:12], sum[:])

	return b
}

func (d *Kmsg) GetCollection(db *database.Database) *database.Collection {
	return db.EndpointsKmsg()
}

func (d *Kmsg) Format(id bson.ObjectID) time.Time {
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

func (d *Kmsg) Handle(db *database.Database) (handled bool, checkAlerts bool,
	err error) {

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
	endpoint bson.ObjectID) (logData LogData, err error) {

	logData = []string{}

	coll := db.EndpointsKmsg()

	limit := int64(settings.Endpoint.KmsgDisplayLimit)

	cursor, err := coll.Find(
		c,
		bson.M{
			"e": endpoint,
		},
		options.Find().
			SetLimit(limit).
			SetSort(bson.D{
				{"b", -1},
				{"s", -1},
			}),
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(c)

	logDataRervse := LogData{}
	for cursor.Next(c) {
		doc := &Kmsg{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		logDataRervse = append(logDataRervse, doc.FormattedLog())
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	for i := len(logDataRervse) - 1; i >= 0; i-- {
		logData = append(logData, logDataRervse[i])
	}

	return
}
