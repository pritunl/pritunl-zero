package database

import (
	"time"

	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/sirupsen/logrus"
)

var metricsCollections = []string{
	"endpoints_system",
	"endpoints_load",
	"endpoints_disk",
	"endpoints_diskio",
	"endpoints_network",
}

var (
	timeSeriesEnabled = map[string]bool{}
)

func IsTimeSeries(name string) bool {
	return timeSeriesEnabled[name]
}

func createTimeSeries(db *Database, name string) (created bool, err error) {
	err = db.database.RunCommand(
		db,
		bson.D{
			{"create", name},
			{"timeseries", bson.D{
				{"timeField", "t"},
				{"metaField", "e"},
				{"granularity", "minutes"},
			}},
			{"expireAfterSeconds", int64(2160 * time.Hour / time.Second)},
		},
	).Err()
	if err != nil {
		err = ParseError(err)
		switch err.(type) {
		case *UnsupportedError:
			err = nil
			return
		default:
			return
		}
	}

	created = true
	return
}

func addTimeSeriesCollections(db *Database,
	collTypes map[string]string) (err error) {

	for _, name := range metricsCollections {
		collType, exists := collTypes[name]

		if exists {
			if collType == "timeseries" {
				timeSeriesEnabled[name] = true
			}
			continue
		}

		created, e := createTimeSeries(db, name)
		if e != nil {
			err = e
			return
		}

		if created {
			timeSeriesEnabled[name] = true

			logrus.WithFields(logrus.Fields{
				"collection": name,
			}).Info("database: Created time series collection")
		}
	}

	return
}
