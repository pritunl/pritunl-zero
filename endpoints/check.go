package endpoints

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
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

type CheckAgg struct {
	Id struct {
		Endpoint  primitive.ObjectID `bson:"e"`
		Timestamp int64              `bson:"t"`
	} `bson:"_id"`
	TargetsUp   int     `bson:"u"`
	TargetsDown int     `bson:"d"`
	LatencyAvg  float64 `bson:"p"`
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

	return
}

func (d *Check) HandleOld(db *database.Database) (handled, checkAlerts bool,
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

func GetCheckChartSingle(c context.Context, db *database.Database,
	checkId primitive.ObjectID, start, end time.Time) (
	chartData ChartData, err error) {

	coll := db.EndpointsCheck()
	chart := NewChart(start, end, time.Minute)

	chck, err := check.Get(db, checkId)
	if err != nil {
		return
	}

	names, err := getRolesNameMapped(db, chck.Roles)
	if err != nil {
		return
	}

	timeQuery := bson.D{
		{"$gte", start},
	}
	if !end.IsZero() {
		timeQuery = append(timeQuery, bson.E{"$lte", end})
	}

	cursor, err := coll.Find(
		c,
		&bson.M{
			"c": checkId,
			"t": timeQuery,
		},
		&options.FindOptions{
			Sort: &bson.D{
				{"t", 1},
			},
		},
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		doc := &Check{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		timestamp := doc.Timestamp.UnixMilli()

		name := names[doc.Endpoint]
		if name == "" {
			println("skip1")
			continue
		}

		chart.Add(name+"-u", timestamp, doc.TargetsUp)
		chart.Add(name+"-d", timestamp, doc.TargetsDown)
		chart.Add(name+"-p", timestamp, doc.LatencyAvg)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	chartData = chart.Export()

	return
}

func GetCheckChart(c context.Context, db *database.Database,
	checkId primitive.ObjectID, start, end time.Time,
	interval time.Duration) (chartData ChartData, err error) {

	if interval == 1*time.Minute {
		chartData, err = GetCheckChartSingle(c, db, checkId, start, end)
		return
	}

	coll := db.EndpointsCheck()
	chart := NewChart(start, end, interval)

	chck, err := check.Get(db, checkId)
	if err != nil {
		return
	}

	names, err := getRolesNameMapped(db, chck.Roles)
	if err != nil {
		return
	}

	timeQuery := bson.D{
		{"$gte", start},
	}
	if !end.IsZero() {
		timeQuery = append(timeQuery, bson.E{"$lte", end})
	}

	cursor, err := coll.Aggregate(c, []*bson.M{
		&bson.M{
			"$match": &bson.M{
				"c": checkId,
				"t": timeQuery,
			},
		},
		&bson.M{
			"$group": &bson.M{
				"_id": &bson.D{
					{"t", &bson.M{
						"$let": &bson.M{
							"vars": &bson.M{
								"t": &bson.D{{"$toLong", "$t"}},
							},
							"in": &bson.M{
								"$subtract": &bson.A{
									"$$t",
									&bson.M{
										"$mod": &bson.A{
											"$$t",
											interval.Milliseconds(),
										},
									},
								},
							},
						},
					}},
					{"e", "$e"},
				},
				"u": &bson.D{
					{"$min", "$u"},
				},
				"d": &bson.D{
					{"$max", "$d"},
				},
				"p": &bson.D{
					{"$avg", "$p"},
				},
			},
		},
		&bson.M{
			"$sort": &bson.M{
				"_id": 1,
			},
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		doc := &CheckAgg{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		name := names[doc.Id.Endpoint]
		if name == "" {
			println("skip2")
			continue
		}

		chart.Add(name+"-u", doc.Id.Timestamp, doc.TargetsUp)
		chart.Add(name+"-d", doc.Id.Timestamp, doc.TargetsDown)
		chart.Add(name+"-p", doc.Id.Timestamp,
			int(math.Round(doc.LatencyAvg)))
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	chartData = chart.Export()

	return
}
