package endpoints

import (
	"context"
	"fmt"
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/alert"
	"github.com/pritunl/pritunl-zero/database"
)

type System struct {
	Id        primitive.Binary   `bson:"_id" json:"id"`
	Endpoint  primitive.ObjectID `bson:"e" json:"e"`
	Timestamp time.Time          `bson:"t" json:"t"`

	Hostname       string           `bson:"-" json:"h"`
	Uptime         uint64           `bson:"-" json:"u"`
	Virtualization string           `bson:"-" json:"v"`
	Platform       string           `bson:"-" json:"p"`
	Processes      uint64           `bson:"pc" json:"pc"`
	CpuCores       int              `bson:"-" json:"cc"`
	CpuUsage       float64          `bson:"cu" json:"cu"`
	MemTotal       int              `bson:"-" json:"mt"`
	MemUsage       float64          `bson:"mu" json:"mu"`
	HugeTotal      int              `bson:"-" json:"ht"`
	HugeUsage      float64          `bson:"hu" json:"hu"`
	SwapTotal      int              `bson:"-" json:"st"`
	SwapUsage      float64          `bson:"su" json:"su"`
	MdStat         []*SystemMdState `bson:"-" json:"ra"`
}

type SystemMdState struct {
	Name   string `bson:"name" json:"n"`
	State  string `bson:"state" json:"s"`
	Level  string `bson:"level" json:"l"`
	Failed int    `bson:"failed" json:"f"`
	Spare  int    `bson:"spare" json:"x"`
	Total  int    `bson:"total" json:"t"`
}

type SystemAgg struct {
	Id        int64   `bson:"_id"`
	CpuUsage  float64 `bson:"cu"`
	MemUsage  float64 `bson:"mu"`
	SwapUsage float64 `bson:"su"`
	HugeUsage float64 `bson:"hu"`
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
		"data.hostname":       d.Hostname,
		"data.uptime":         d.Uptime,
		"data.virtualization": d.Virtualization,
		"data.platform":       d.Platform,
		"data.cpu_cores":      d.CpuCores,
		"data.mem_total":      d.MemTotal,
		"data.swap_total":     d.SwapTotal,
		"data.huge_total":     d.HugeTotal,
		"data.md_stat":        d.MdStat,
	}
}

func (d *System) CheckAlerts(resources []*alert.Alert) (alerts []*Alert) {
	alerts = []*Alert{}

	for _, resource := range resources {
		switch resource.Resource {
		case alert.SystemCpuLevel:
			if d.CpuUsage > float64(resource.ValueInt) {
				alerts = []*Alert{
					NewAlert(resource, fmt.Sprintf(
						"System cpu high usage (%.2f%%)",
						d.MemUsage,
					)),
				}
			}
			break
		case alert.SystemMemoryLevel:
			if d.MemUsage > float64(resource.ValueInt) {
				alerts = []*Alert{
					NewAlert(resource, fmt.Sprintf(
						"System low on memory (%.2f%%)",
						d.MemUsage,
					)),
				}
			}
			break
		case alert.SystemSwapLevel:
			if d.SwapUsage > float64(resource.ValueInt) {
				alerts = []*Alert{
					NewAlert(resource, fmt.Sprintf(
						"System low on swap (%.2f%%)",
						d.SwapUsage,
					)),
				}
			}
			break
		case alert.SystemHugePagesLevel:
			if d.SwapUsage > float64(resource.ValueInt) {
				alerts = []*Alert{
					NewAlert(resource, fmt.Sprintf(
						"System low on hugepages (%.2f%%)",
						d.SwapUsage,
					)),
				}
			}
			break
		case alert.SystemMdFailed:
			if d.MdStat != nil {
				for _, md := range d.MdStat {
					if md.Failed > 0 {
						alerts = []*Alert{
							NewAlert(resource, fmt.Sprintf(
								"System MD RAID device failed (%s %s)",
								md.Name,
								md.Level,
							)),
						}
					}
				}
			}
			break
		}
	}

	return
}

func (d *System) Handle(db *database.Database) (handled, checkAlerts bool,
	err error) {

	return
}

func GetSystemChartSingle(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time) (
	chartData ChartData, err error) {

	coll := db.EndpointsSystem()
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
		doc := &System{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		timestamp := doc.Timestamp.UnixMilli()

		chart.Add("cpu_usage", timestamp, doc.CpuUsage)
		chart.Add("mem_usage", timestamp, doc.MemUsage)
		chart.Add("swap_usage", timestamp, doc.SwapUsage)
		chart.Add("huge_usage", timestamp, doc.HugeUsage)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	chartData = chart.Export()

	return
}

func GetSystemChart(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time,
	interval time.Duration) (chartData ChartData, err error) {

	if interval == 1*time.Minute {
		chartData, err = GetSystemChartSingle(c, db, endpoint, start, end)
		return
	}

	coll := db.EndpointsSystem()
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
				"hu": &bson.D{
					{"$avg", "$hu"},
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

		chart.Add("cpu_usage", doc.Id, doc.CpuUsage)
		chart.Add("mem_usage", doc.Id, doc.MemUsage)
		chart.Add("swap_usage", doc.Id, doc.SwapUsage)
		chart.Add("huge_usage", doc.Id, doc.HugeUsage)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	chartData = chart.Export()

	return
}
