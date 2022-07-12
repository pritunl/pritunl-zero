package endpoints

import (
	"context"
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/alert"
	"github.com/pritunl/pritunl-zero/database"
)

type DiskIo struct {
	Id        primitive.Binary   `bson:"_id" json:"id"`
	Endpoint  primitive.ObjectID `bson:"e" json:"e"`
	Timestamp time.Time          `bson:"t" json:"t"`

	Disks []*DiskIoDisk `bson:"d" json:"d"`
}

type DiskIoDisk struct {
	Name       string `bson:"n" json:"n"`
	BytesRead  uint64 `bson:"br" json:"br"`
	BytesWrite uint64 `bson:"bw" json:"bw"`
	CountRead  uint64 `bson:"cr" json:"cr"`
	CountWrite uint64 `bson:"cw" json:"cw"`
	TimeRead   uint64 `bson:"tr" json:"tr"`
	TimeWrite  uint64 `bson:"tw" json:"tw"`
	TimeIo     uint64 `bson:"ti" json:"ti"`
}

type DiskStatic struct {
	Name string `bson:"n" json:"n"`
}

func ParseDisk(i *DiskIoDisk) *DiskStatic {
	return &DiskStatic{
		Name: i.Name,
	}
}

type DiskIoAgg struct {
	Id struct {
		Disk      string `bson:"n"`
		Timestamp int64  `bson:"t"`
	} `bson:"_id"`
	BytesRead  uint64 `bson:"br"`
	BytesWrite uint64 `bson:"bw"`
	CountRead  uint64 `bson:"cr"`
	CountWrite uint64 `bson:"cw"`
	TimeRead   uint64 `bson:"tr"`
	TimeWrite  uint64 `bson:"tw"`
	TimeIo     uint64 `bson:"ti"`
}

func (d *DiskIo) GetCollection(db *database.Database) *database.Collection {
	return db.EndpointsDiskIo()
}

func (d *DiskIo) Format(id primitive.ObjectID) time.Time {
	d.Endpoint = id
	d.Timestamp = d.Timestamp.UTC().Truncate(1 * time.Minute)
	d.Id = GenerateId(id, d.Timestamp)
	return d.Timestamp
}

func (d *DiskIo) StaticData() *bson.M {
	disks := []*DiskStatic{}

	for _, dsk := range d.Disks {
		disks = append(disks, ParseDisk(dsk))
	}

	return &bson.M{
		"data.disks": disks,
	}
}

func (d *DiskIo) CheckAlerts(resources []*alert.Resource) (alerts []*Alert) {
	return
}

func GetDiskIoChartSingle(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time) (
	chartData ChartData, err error) {

	coll := db.EndpointsDiskIo()
	chart := NewChart(start, end, time.Minute)

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
		doc := &DiskIo{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		for _, dsk := range doc.Disks {
			timestamp := doc.Timestamp.UnixMilli()

			chart.Add(dsk.Name+"-br", timestamp, dsk.BytesRead)
			chart.Add(dsk.Name+"-bw", timestamp, dsk.BytesWrite)
			//chart.Add(dsk.Name+"-cr", timestamp, dsk.CountRead)
			//chart.Add(dsk.Name+"-cw", timestamp, dsk.CountWrite)
			chart.Add(dsk.Name+"-tr", timestamp, dsk.TimeRead)
			chart.Add(dsk.Name+"-tw", timestamp, dsk.TimeWrite)
			chart.Add(dsk.Name+"-ti", timestamp, dsk.TimeIo)
		}
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	chartData = chart.Export()

	return
}

func GetDiskIoChart(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time,
	interval time.Duration) (chartData ChartData, err error) {

	if interval == 1*time.Minute {
		chartData, err = GetDiskIoChartSingle(c, db, endpoint, start, end)
		return
	}

	coll := db.EndpointsDiskIo()
	chart := NewChart(start, end, interval)

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
			"$unwind": "$d",
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
					{"n", "$d.n"},
				},
				"br": &bson.D{
					{"$sum", "$d.br"},
				},
				"bw": &bson.D{
					{"$sum", "$d.bw"},
				},
				//"cr": &bson.D{
				//	{"$sum", "$d.cr"},
				//},
				//"cw": &bson.D{
				//	{"$sum", "$d.cw"},
				//},
				"tr": &bson.D{
					{"$sum", "$d.tr"},
				},
				"tw": &bson.D{
					{"$sum", "$d.tw"},
				},
				"ti": &bson.D{
					{"$sum", "$d.ti"},
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
		doc := &DiskIoAgg{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		chart.Add(doc.Id.Disk+"-br", doc.Id.Timestamp, doc.BytesRead)
		chart.Add(doc.Id.Disk+"-bw", doc.Id.Timestamp, doc.BytesWrite)
		//chart.Add(doc.Id.Disk+"-cr", doc.Id.Timestamp, doc.CountRead)
		//chart.Add(doc.Id.Disk+"-cw", doc.Id.Timestamp, doc.CountWrite)
		chart.Add(doc.Id.Disk+"-tr", doc.Id.Timestamp, doc.TimeRead)
		chart.Add(doc.Id.Disk+"-tw", doc.Id.Timestamp, doc.TimeWrite)
		chart.Add(doc.Id.Disk+"-ti", doc.Id.Timestamp, doc.TimeIo)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	chartData = chart.Export()

	return
}
