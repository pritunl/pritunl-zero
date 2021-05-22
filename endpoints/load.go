package endpoints

import (
	"context"
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
)

type Load struct {
	Id        primitive.Binary   `bson:"_id" json:"id"`
	Endpoint  primitive.ObjectID `bson:"e" json:"e"`
	Timestamp time.Time          `bson:"t" json:"t"`

	Load1  float64 `bson:"lx" json:"lx"`
	Load5  float64 `bson:"ly" json:"ly"`
	Load15 float64 `bson:"lz" json:"lz"`
}

type LoadAgg struct {
	Id     int64   `bson:"_id"`
	Load1  float64 `bson:"lx"`
	Load5  float64 `bson:"ly"`
	Load15 float64 `bson:"lz"`
}

func (d *Load) GetCollection(db *database.Database) *database.Collection {
	return db.EndpointsLoad()
}

func (d *Load) Format(id primitive.ObjectID) time.Time {
	d.Endpoint = id
	d.Timestamp = d.Timestamp.UTC().Truncate(1 * time.Minute)
	d.Id = GenerateId(id, d.Timestamp)
	return d.Timestamp
}

func (d *Load) StaticData() *bson.M {
	return nil
}

type LoadChart struct {
	Load1  []*ChartFloat `json:"load1"`
	Load5  []*ChartFloat `json:"load5"`
	Load15 []*ChartFloat `json:"load15"`
}

func GetLoadChartSingle(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time) (
	chart *LoadChart, err error) {

	coll := db.EndpointsLoad()
	cpuUsage := []*ChartFloat{}
	memUsage := []*ChartFloat{}
	swapUsage := []*ChartFloat{}

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
		doc := &Load{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		cpuUsage = append(cpuUsage, &ChartFloat{
			X: doc.Timestamp.Unix() * 1000,
			Y: doc.Load1,
		})
		memUsage = append(memUsage, &ChartFloat{
			X: doc.Timestamp.Unix() * 1000,
			Y: doc.Load5,
		})
		swapUsage = append(swapUsage, &ChartFloat{
			X: doc.Timestamp.Unix() * 1000,
			Y: doc.Load15,
		})
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	chart = &LoadChart{
		Load1:  cpuUsage,
		Load5:  memUsage,
		Load15: swapUsage,
	}

	return
}

func GetLoadChart(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time,
	interval time.Duration) (chart *LoadChart, err error) {

	if interval == 1*time.Minute {
		chart, err = GetLoadChartSingle(c, db, endpoint, start, end)
		return
	}

	coll := db.EndpointsLoad()
	cpuUsage := []*ChartFloat{}
	memUsage := []*ChartFloat{}
	swapUsage := []*ChartFloat{}

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
				"lx": &bson.D{
					{"$avg", "$lx"},
				},
				"ly": &bson.D{
					{"$avg", "$ly"},
				},
				"lz": &bson.D{
					{"$avg", "$lz"},
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
		doc := &LoadAgg{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		cpuUsage = append(cpuUsage, &ChartFloat{
			X: doc.Id,
			Y: doc.Load1,
		})
		memUsage = append(memUsage, &ChartFloat{
			X: doc.Id,
			Y: doc.Load5,
		})
		swapUsage = append(swapUsage, &ChartFloat{
			X: doc.Id,
			Y: doc.Load15,
		})
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	chart = &LoadChart{
		Load1:  cpuUsage,
		Load5:  memUsage,
		Load15: swapUsage,
	}

	return
}
