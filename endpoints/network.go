package endpoints

import (
	"context"
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
)

type Network struct {
	Id        primitive.Binary   `bson:"_id" json:"id"`
	Endpoint  primitive.ObjectID `bson:"e" json:"e"`
	Timestamp time.Time          `bson:"t" json:"t"`

	Interfaces []*Interface `bson:"i" json:"i"`
}

type Interface struct {
	Name        string `bson:"n" json:"n"`
	BytesSent   uint64 `bson:"bs" json:"bs"`
	BytesRecv   uint64 `bson:"br" json:"br"`
	PacketsSent uint64 `bson:"ps" json:"ps"`
	PacketsRecv uint64 `bson:"pr" json:"pr"`
	ErrorsSent  uint64 `bson:"es" json:"es"`
	ErrorsRecv  uint64 `bson:"er" json:"er"`
	DropsSent   uint64 `bson:"ds" json:"ds"`
	DropsRecv   uint64 `bson:"dr" json:"dr"`
	FifoSent    uint64 `bson:"fs" json:"fs"`
	FifoRecv    uint64 `bson:"fr" json:"fr"`
}

type InterfaceStatic struct {
	Name string `bson:"n" json:"n"`
}

type InterfaceChart struct {
	Path string       `json:"path"`
	Data []*ChartUint `json:"data"`
}

func ParseInterface(i *Interface) *InterfaceStatic {
	return &InterfaceStatic{
		Name: i.Name,
	}
}

type NetworkAgg struct {
	Id struct {
		Interface string `bson:"i"`
		Timestamp int64  `bson:"t"`
	} `bson:"_id"`
	BytesSent   uint64 `bson:"bs"`
	BytesRecv   uint64 `bson:"br"`
	PacketsSent uint64 `bson:"ps"`
	PacketsRecv uint64 `bson:"pr"`
	ErrorsSent  uint64 `bson:"es"`
	ErrorsRecv  uint64 `bson:"er"`
	DropsSent   uint64 `bson:"ds"`
	DropsRecv   uint64 `bson:"dr"`
	FifoSent    uint64 `bson:"fs"`
	FifoRecv    uint64 `bson:"fr"`
}

func (d *Network) GetCollection(db *database.Database) *database.Collection {
	return db.EndpointsNetwork()
}

func (d *Network) Format(id primitive.ObjectID) time.Time {
	d.Endpoint = id
	d.Timestamp = d.Timestamp.UTC().Truncate(1 * time.Minute)
	d.Id = GenerateId(id, d.Timestamp)
	return d.Timestamp
}

func (d *Network) StaticData() *bson.M {
	ifaces := []*InterfaceStatic{}

	for _, iface := range d.Interfaces {
		ifaces = append(ifaces, ParseInterface(iface))
	}

	return &bson.M{
		"data.interfaces": ifaces,
	}
}

type NetworkChart struct {
	Interfaces []*InterfaceChart `json:"interfaces"`
}

func GetNetworkChartSingle(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time) (
	chart map[string][]*ChartUint, err error) {

	coll := db.EndpointsNetwork()
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
		doc := &Network{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		for _, iface := range doc.Interfaces {
			pathInterfaces := chart[iface.Name]
			if pathInterfaces == nil {
				pathInterfaces = []*ChartUint{}
			}
			chart[iface.Name+"-us"] = append(pathInterfaces, &ChartUint{
				X: doc.Timestamp.Unix() * 1000,
				Y: iface.BytesSent,
			})
			chart[iface.Name+"-ur"] = append(pathInterfaces, &ChartUint{
				X: doc.Timestamp.Unix() * 1000,
				Y: iface.BytesSent,
			})

			chart[iface.Name+"-ps"] = append(pathInterfaces, &ChartUint{
				X: doc.Timestamp.Unix() * 1000,
				Y: iface.PacketsSent,
			})
			chart[iface.Name+"-pr"] = append(pathInterfaces, &ChartUint{
				X: doc.Timestamp.Unix() * 1000,
				Y: iface.PacketsSent,
			})

			chart[iface.Name+"-es"] = append(pathInterfaces, &ChartUint{
				X: doc.Timestamp.Unix() * 1000,
				Y: iface.ErrorsSent,
			})
			chart[iface.Name+"-er"] = append(pathInterfaces, &ChartUint{
				X: doc.Timestamp.Unix() * 1000,
				Y: iface.ErrorsSent,
			})

			chart[iface.Name+"-ds"] = append(pathInterfaces, &ChartUint{
				X: doc.Timestamp.Unix() * 1000,
				Y: iface.DropsSent,
			})
			chart[iface.Name+"-dr"] = append(pathInterfaces, &ChartUint{
				X: doc.Timestamp.Unix() * 1000,
				Y: iface.DropsSent,
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

func GetNetworkChart(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time,
	interval time.Duration) (chart map[string][]*ChartUint, err error) {

	if interval == 1*time.Minute {
		chart, err = GetNetworkChartSingle(c, db, endpoint, start, end)
		return
	}

	coll := db.EndpointsNetwork()
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
			"$unwind": "$m",
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
					{"i", "$m.i"},
				},
				"bs": &bson.D{
					{"$avg", "$m.bs"},
				},
				"br": &bson.D{
					{"$avg", "$m.br"},
				},
				"ps": &bson.D{
					{"$avg", "$m.ps"},
				},
				"pr": &bson.D{
					{"$avg", "$m.pr"},
				},
				"es": &bson.D{
					{"$avg", "$m.es"},
				},
				"er": &bson.D{
					{"$avg", "$m.er"},
				},
				"ds": &bson.D{
					{"$avg", "$m.ds"},
				},
				"dr": &bson.D{
					{"$avg", "$m.dr"},
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
		doc := &NetworkAgg{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		pathInterfaces := chart[doc.Id.Interface]
		if pathInterfaces == nil {
			pathInterfaces = []*ChartUint{}
		}
		chart[doc.Id.Interface+"-bs"] = append(pathInterfaces, &ChartUint{
			X: doc.Id.Timestamp,
			Y: doc.BytesSent,
		})
		chart[doc.Id.Interface+"-br"] = append(pathInterfaces, &ChartUint{
			X: doc.Id.Timestamp,
			Y: doc.BytesRecv,
		})

		chart[doc.Id.Interface+"-ps"] = append(pathInterfaces, &ChartUint{
			X: doc.Id.Timestamp,
			Y: doc.PacketsSent,
		})
		chart[doc.Id.Interface+"-pr"] = append(pathInterfaces, &ChartUint{
			X: doc.Id.Timestamp,
			Y: doc.PacketsRecv,
		})

		chart[doc.Id.Interface+"-es"] = append(pathInterfaces, &ChartUint{
			X: doc.Id.Timestamp,
			Y: doc.ErrorsSent,
		})
		chart[doc.Id.Interface+"-er"] = append(pathInterfaces, &ChartUint{
			X: doc.Id.Timestamp,
			Y: doc.ErrorsRecv,
		})

		chart[doc.Id.Interface+"-ds"] = append(pathInterfaces, &ChartUint{
			X: doc.Id.Timestamp,
			Y: doc.DropsSent,
		})
		chart[doc.Id.Interface+"-dr"] = append(pathInterfaces, &ChartUint{
			X: doc.Id.Timestamp,
			Y: doc.DropsRecv,
		})
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
