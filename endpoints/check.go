package endpoints

import (
	"fmt"
	"math"
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/alert"
	"github.com/pritunl/pritunl-zero/check"
	"github.com/pritunl/pritunl-zero/database"
)

type Check struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Check     primitive.ObjectID `bson:"c" json:"c"`
	Endpoint  primitive.ObjectID `bson:"e" json:"e"`
	Timestamp time.Time          `bson:"t" json:"t"`

	TargetsUp   int `bson:"u" json:"u"`
	TargetsDown int `bson:"d" json:"d"`
	LatencyAvg  int `bson:"p" json:"p"`

	TargetsIn []string `bson:"-" json:"x"`
	LatencyIn []int    `bson:"-" json:"l"`
	ErrorsIn  []string `bson:"-" json:"r"`

	checkName string `bson:"-" json:"-"`
}

func (d *Check) GetCollection(db *database.Database) *database.Collection {
	return db.EndpointsCheck()
}

func (d *Check) Format(id primitive.ObjectID) time.Time {
	d.Endpoint = id
	d.Timestamp = d.Timestamp.UTC().Truncate(10 * time.Second)
	d.Id = GenerateId(id, d.Timestamp)

	if d.TargetsIn == nil {
		d.TargetsIn = []string{}
	}
	if d.LatencyIn == nil {
		d.LatencyIn = []int{}
	}
	if d.ErrorsIn == nil {
		d.ErrorsIn = []string{}
	}

	count := 0
	for _, e := range d.ErrorsIn {
		if e == "" {
			count += 1
		}
	}
	d.TargetsUp = count

	count = 0
	for _, e := range d.ErrorsIn {
		if e != "" {
			count += 1
		}
	}
	d.TargetsDown = count

	count = 0
	avg := 0
	for _, lat := range d.LatencyIn {
		if lat > 0 {
			avg += lat
			count += 1
		}
	}

	d.LatencyAvg = int(math.Round(float64(avg) / float64(count)))

	return d.Timestamp
}

func (d *Check) StaticData() *bson.M {
	return nil
}

func (d *Check) CheckAlerts(resources []*alert.Alert) (alerts []*Alert) {
	alerts = []*Alert{}

	for _, resource := range resources {
		switch resource.Resource {
		case alert.CheckHttpFailed:
			for _, er := range d.ErrorsIn {
				if er != "" {
					alerts = []*Alert{
						NewAlert(resource, fmt.Sprintf(
							"Check HTTP error: %s %s",
							d.checkName,
							er,
						)),
					}
					break
				}
			}
			break
		}
	}

	return
}

func (d *Check) Handle(db *database.Database) (handled, checkAlerts bool,
	err error) {

	handled = true

	chck, err := check.Get(db, d.Check)
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			err = nil
		}
		return
	}

	d.checkName = chck.Name

	state := &check.State{
		Endpoint:  d.Endpoint,
		Timestamp: d.Timestamp,
		Targets:   d.TargetsIn,
		Latency:   d.LatencyIn,
		Errors:    d.ErrorsIn,
	}

	checkAlerts, err = chck.UpdateState(db, state)
	if err != nil {
		return
	}

	return
}
