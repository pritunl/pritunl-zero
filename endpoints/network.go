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

func ParseInterface(i *Interface) *InterfaceStatic {
	return &InterfaceStatic{
		Name: i.Name,
	}
}

type NetworkAgg struct {
	Id struct {
		Interface string `bson:"n"`
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

func (d *Network) CheckAlerts(resources []*alert.Resource) (alerts []*Alert) {
	return
}

func GetNetworkChartSingle(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time) (
	chartData ChartData, err error) {

	coll := db.EndpointsNetwork()
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
		doc := &Network{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		for _, iface := range doc.Interfaces {
			timestamp := doc.Timestamp.UnixMilli()

			chart.Add(iface.Name+"-bs", timestamp, iface.BytesSent)
			chart.Add(iface.Name+"-br", timestamp, iface.BytesRecv)
			//chart.Add(iface.Name+"-ps", timestamp, iface.PacketsSent)
			//chart.Add(iface.Name+"-pr", timestamp, iface.PacketsRecv)
			//chart.Add(iface.Name+"-es", timestamp, iface.ErrorsSent)
			//chart.Add(iface.Name+"-er", timestamp, iface.ErrorsRecv)
			//chart.Add(iface.Name+"-ds", timestamp, iface.DropsSent)
			//chart.Add(iface.Name+"-dr", timestamp, iface.DropsRecv)
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

func GetNetworkChart(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time,
	interval time.Duration) (chartData ChartData, err error) {

	if interval == 1*time.Minute {
		chartData, err = GetNetworkChartSingle(c, db, endpoint, start, end)
		return
	}

	coll := db.EndpointsNetwork()
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
			"$unwind": "$i",
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
					{"n", "$i.n"},
				},
				"bs": &bson.D{
					{"$sum", "$i.bs"},
				},
				"br": &bson.D{
					{"$sum", "$i.br"},
				},
				//"ps": &bson.D{
				//	{"$avg", "$i.ps"},
				//},
				//"pr": &bson.D{
				//	{"$avg", "$i.pr"},
				//},
				//"es": &bson.D{
				//	{"$avg", "$i.es"},
				//},
				//"er": &bson.D{
				//	{"$avg", "$i.er"},
				//},
				//"ds": &bson.D{
				//	{"$avg", "$i.ds"},
				//},
				//"dr": &bson.D{
				//	{"$avg", "$i.dr"},
				//},
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

		chart.Add(doc.Id.Interface+"-bs", doc.Id.Timestamp, doc.BytesSent)
		chart.Add(doc.Id.Interface+"-br", doc.Id.Timestamp, doc.BytesRecv)
		//chart.Add(doc.Id.Interface+"-ps", doc.Id.Timestamp, doc.PacketsSent)
		//chart.Add(doc.Id.Interface+"-pr", doc.Id.Timestamp, doc.PacketsRecv)
		//chart.Add(doc.Id.Interface+"-es", doc.Id.Timestamp, doc.ErrorsSent)
		//chart.Add(doc.Id.Interface+"-er", doc.Id.Timestamp, doc.ErrorsRecv)
		//chart.Add(doc.Id.Interface+"-ds", doc.Id.Timestamp, doc.DropsSent)
		//chart.Add(doc.Id.Interface+"-dr", doc.Id.Timestamp, doc.DropsRecv)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	chartData = chart.Export()

	return
}
