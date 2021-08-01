package endpoints

import (
	"context"
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
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

type DiskIoDiskChart struct {
	Path string       `json:"path"`
	Data []*ChartUint `json:"data"`
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

type DiskIoChart struct {
	Disks []*DiskIoDiskChart `json:"disks"`
}

func GetDiskIoChartSingle(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time) (
	chart map[string][]*ChartUint, err error) {

	coll := db.EndpointsDiskIo()
	chart = map[string][]*ChartUint{}

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
			timestamp := doc.Timestamp.Unix() * 1000

			dskChart := chart[dsk.Name+"-br"]
			if dskChart == nil {
				dskChart = []*ChartUint{}
			}
			chart[dsk.Name+"-br"] = append(dskChart, &ChartUint{
				X: timestamp,
				Y: dsk.BytesRead,
			})
			dskChart = chart[dsk.Name+"-bw"]
			if dskChart == nil {
				dskChart = []*ChartUint{}
			}
			chart[dsk.Name+"-bw"] = append(dskChart, &ChartUint{
				X: timestamp,
				Y: dsk.BytesWrite,
			})

			//dskChart = chart[dsk.Name+"-cr"]
			//if dskChart == nil {
			//	dskChart = []*ChartUint{}
			//}
			//chart[dsk.Name+"-cr"] = append(dskChart, &ChartUint{
			//	X: timestamp,
			//	Y: dsk.CountRead,
			//})
			//dskChart = chart[dsk.Name+"-cw"]
			//if dskChart == nil {
			//	dskChart = []*ChartUint{}
			//}
			//chart[dsk.Name+"-cw"] = append(dskChart, &ChartUint{
			//	X: timestamp,
			//	Y: dsk.CountWrite,
			//})

			dskChart = chart[dsk.Name+"-tr"]
			if dskChart == nil {
				dskChart = []*ChartUint{}
			}
			chart[dsk.Name+"-tr"] = append(dskChart, &ChartUint{
				X: timestamp,
				Y: dsk.TimeRead,
			})
			dskChart = chart[dsk.Name+"-tw"]
			if dskChart == nil {
				dskChart = []*ChartUint{}
			}
			chart[dsk.Name+"-tw"] = append(dskChart, &ChartUint{
				X: timestamp,
				Y: dsk.TimeWrite,
			})
			dskChart = chart[dsk.Name+"-ti"]
			if dskChart == nil {
				dskChart = []*ChartUint{}
			}
			chart[dsk.Name+"-ti"] = append(dskChart, &ChartUint{
				X: timestamp,
				Y: dsk.TimeIo,
			})
		}
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetDiskIoChart(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time,
	interval time.Duration) (chart map[string][]*ChartUint, err error) {

	if interval == 1*time.Minute {
		chart, err = GetDiskIoChartSingle(c, db, endpoint, start, end)
		return
	}

	coll := db.EndpointsDiskIo()
	chart = map[string][]*ChartUint{}

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

		dskChart := chart[doc.Id.Disk+"-br"]
		if dskChart == nil {
			dskChart = []*ChartUint{}
		}
		chart[doc.Id.Disk+"-br"] = append(dskChart, &ChartUint{
			X: doc.Id.Timestamp,
			Y: doc.BytesRead,
		})
		dskChart = chart[doc.Id.Disk+"-bw"]
		if dskChart == nil {
			dskChart = []*ChartUint{}
		}
		chart[doc.Id.Disk+"-bw"] = append(dskChart, &ChartUint{
			X: doc.Id.Timestamp,
			Y: doc.BytesWrite,
		})

		//dskChart = chart[doc.Id.Disk+"-cr"]
		//if dskChart == nil {
		//	dskChart = []*ChartUint{}
		//}
		//chart[doc.Id.Disk+"-cr"] = append(dskChart, &ChartUint{
		//	X: doc.Id.Timestamp,
		//	Y: doc.CountRead,
		//})
		//dskChart = chart[doc.Id.Disk+"-cw"]
		//if dskChart == nil {
		//	dskChart = []*ChartUint{}
		//}
		//chart[doc.Id.Disk+"-cw"] = append(dskChart, &ChartUint{
		//	X: doc.Id.Timestamp,
		//	Y: doc.CountWrite,
		//})

		dskChart = chart[doc.Id.Disk+"-tr"]
		if dskChart == nil {
			dskChart = []*ChartUint{}
		}
		chart[doc.Id.Disk+"-tr"] = append(dskChart, &ChartUint{
			X: doc.Id.Timestamp,
			Y: doc.TimeRead,
		})
		dskChart = chart[doc.Id.Disk+"-tw"]
		if dskChart == nil {
			dskChart = []*ChartUint{}
		}
		chart[doc.Id.Disk+"-tw"] = append(dskChart, &ChartUint{
			X: doc.Id.Timestamp,
			Y: doc.TimeWrite,
		})
		dskChart = chart[doc.Id.Disk+"-ti"]
		if dskChart == nil {
			dskChart = []*ChartUint{}
		}
		chart[doc.Id.Disk+"-ti"] = append(dskChart, &ChartUint{
			X: doc.Id.Timestamp,
			Y: doc.TimeIo,
		})
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
