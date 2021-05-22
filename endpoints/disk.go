package endpoints

import (
	"context"
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
)

type Disk struct {
	Id        primitive.Binary   `bson:"_id" json:"id"`
	Endpoint  primitive.ObjectID `bson:"e" json:"e"`
	Timestamp time.Time          `bson:"t" json:"t"`

	Mounts []*Mount `bson:"m" json:"m"`
}

type Mount struct {
	Path   string  `bson:"p" json:"p"`
	Format string  `bson:"-" json:"f"`
	Size   uint64  `bson:"-" json:"s"`
	Used   float64 `bson:"u" json:"u"`
}

type MountStatic struct {
	Path   string `bson:"p" json:"p"`
	Format string `bson:"f" json:"f"`
	Size   uint64 `bson:"s" json:"s"`
}

type MountChart struct {
	Path string        `json:"path"`
	Data []*ChartFloat `json:"data"`
}

func ParseMount(mn *Mount) *MountStatic {
	return &MountStatic{
		Path:   mn.Path,
		Format: mn.Format,
		Size:   mn.Size,
	}
}

type DiskAgg struct {
	Id struct {
		Path      string `bson:"p"`
		Timestamp int64  `bson:"t"`
	} `bson:"_id"`
	Used float64 `bson:"u"`
}

func (d *Disk) GetCollection(db *database.Database) *database.Collection {
	return db.EndpointsDisk()
}

func (d *Disk) Format(id primitive.ObjectID) time.Time {
	d.Endpoint = id
	d.Timestamp = d.Timestamp.UTC().Truncate(1 * time.Minute)
	d.Id = GenerateId(id, d.Timestamp)
	return d.Timestamp
}

func (d *Disk) StaticData() *bson.M {
	mounts := []*MountStatic{}

	for _, mount := range d.Mounts {
		mounts = append(mounts, ParseMount(mount))
	}

	return &bson.M{
		"data.mounts": mounts,
	}
}

type DiskChart struct {
	Mounts []*MountChart `json:"mounts"`
}

func GetDiskChartSingle(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time) (
	chart map[string][]*ChartFloat, err error) {

	coll := db.EndpointsDisk()
	chart = map[string][]*ChartFloat{}

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
		doc := &Disk{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		for _, mount := range doc.Mounts {
			pathMounts := chart[mount.Path]
			if pathMounts == nil {
				pathMounts = []*ChartFloat{}
			}
			chart[mount.Path] = append(pathMounts, &ChartFloat{
				X: doc.Timestamp.Unix() * 1000,
				Y: mount.Used,
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

func GetDiskChart(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, start, end time.Time,
	interval time.Duration) (chart map[string][]*ChartFloat, err error) {

	if interval == 1*time.Minute {
		chart, err = GetDiskChartSingle(c, db, endpoint, start, end)
		return
	}

	coll := db.EndpointsDisk()
	chart = map[string][]*ChartFloat{}

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
					{"p", "$m.p"},
				},
				"u": &bson.D{
					{"$avg", "$m.u"},
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
		doc := &DiskAgg{}
		err = cursor.Decode(doc)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		pathMounts := chart[doc.Id.Path]
		if pathMounts == nil {
			pathMounts = []*ChartFloat{}
		}
		chart[doc.Id.Path] = append(pathMounts, &ChartFloat{
			X: doc.Id.Timestamp,
			Y: doc.Used,
		})
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
