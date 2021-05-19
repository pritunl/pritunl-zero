package endpoints

import (
	"context"
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
)

type System struct {
	Id        primitive.Binary   `bson:"_id" json:"id"`
	Endpoint  primitive.ObjectID `bson:"e" json:"e"`
	Timestamp time.Time          `bson:"t" json:"t"`

	CpuUsage  float64 `bson:"cu" json:"cu"`
	MemTotal  int     `bson:"-" json:"mt"`
	MemUsage  float64 `bson:"mu" json:"mu"`
	SwapTotal int     `bson:"-" json:"st"`
	SwapUsage float64 `bson:"su" json:"su"`
}

type SystemAgg struct {
	Id        int64   `bson:"_id"`
	CpuUsage  float64 `bson:"cu"`
	MemUsage  float64 `bson:"mu"`
	SwapUsage float64 `bson:"su"`
}

func (d *System) GetCollection(db *database.Database) *database.Collection {
	return db.EndpointsSystem()
}

func (d *System) Format(id primitive.ObjectID) time.Time {
	d.Endpoint = id
	d.Timestamp = d.Timestamp.UTC().Truncate(1 * time.Minute)
	d.Id = GenerateId(id, d.Timestamp)
	return d.Timestamp
}

func (d *System) StaticData() *bson.M {
	return &bson.M{
		"data.mem_total":  d.MemTotal,
		"data.swap_total": d.SwapTotal,
	}
}

type SystemChart struct {
	CpuUsage  []*Chart `json:"cpu_usage"`
	MemUsage  []*Chart `json:"mem_usage"`
	SwapUsage []*Chart `json:"swap_usage"`
}

func GetSystemChartSingle(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time) (
	chart *SystemChart, err error) {

	coll := db.EndpointsSystem()
	cpuUsage := []*Chart{}
	memUsage := []*Chart{}
	swapUsage := []*Chart{}

	timeQuery := bson.D{
		{"$gte", start},
	}
	if !end.IsZero() {
		timeQuery = append(timeQuery, bson.E{"$lte", end})
	}

	cursor, err := coll.Find(
		c,
		&bson.M{
			"e": endpoint,
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
		doc := &System{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		cpuUsage = append(cpuUsage, &Chart{
			X: doc.Timestamp.Unix() * 1000,
			Y: doc.CpuUsage,
		})
		memUsage = append(memUsage, &Chart{
			X: doc.Timestamp.Unix() * 1000,
			Y: doc.MemUsage,
		})
		swapUsage = append(swapUsage, &Chart{
			X: doc.Timestamp.Unix() * 1000,
			Y: doc.MemUsage,
		})
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	chart = &SystemChart{
		CpuUsage:  cpuUsage,
		MemUsage:  memUsage,
		SwapUsage: swapUsage,
	}

	return
}

func GetSystemChart(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time,
	interval time.Duration) (chart *SystemChart, err error) {

	if interval == 1*time.Minute {
		chart, err = GetSystemChartSingle(c, db, endpoint, start, end)
		return
	}

	coll := db.EndpointsSystem()
	cpuUsage := []*Chart{}
	memUsage := []*Chart{}
	swapUsage := []*Chart{}

	timeQuery := bson.D{
		{"$gte", start},
	}
	if !end.IsZero() {
		timeQuery = append(timeQuery, bson.E{"$lte", end})
	}

	cursor, err := coll.Aggregate(c, []*bson.M{
		&bson.M{
			"$match": &bson.M{
				"e": endpoint,
				"t": timeQuery,
			},
		},
		&bson.M{
			"$group": &bson.M{
				"_id": &bson.M{
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
				},
				"cu": &bson.D{
					{"$avg", "$cu"},
				},
				"mu": &bson.D{
					{"$avg", "$mu"},
				},
				"su": &bson.D{
					{"$avg", "$su"},
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
		doc := &SystemAgg{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		cpuUsage = append(cpuUsage, &Chart{
			X: doc.Id,
			Y: doc.CpuUsage,
		})
		memUsage = append(memUsage, &Chart{
			X: doc.Id,
			Y: doc.MemUsage,
		})
		swapUsage = append(swapUsage, &Chart{
			X: doc.Id,
			Y: doc.MemUsage,
		})
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	chart = &SystemChart{
		CpuUsage:  cpuUsage,
		MemUsage:  memUsage,
		SwapUsage: swapUsage,
	}

	return
}
