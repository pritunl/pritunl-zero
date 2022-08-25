package endpoints

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/alert"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
)

const (
	BinaryMD5 byte = 0x05
)

type Doc interface {
	GetCollection(*database.Database) *database.Collection
	Format(primitive.ObjectID) time.Time
	StaticData() *bson.M
	CheckAlerts(resources []*alert.Alert) []*Alert
}

type Point struct {
	X int64       `json:"x"`
	Y interface{} `json:"y"`
}

type ChartData = map[string][]*Point

type LogData = []string

func GenerateId(endpointId primitive.ObjectID,
	timestamp time.Time) primitive.Binary {

	hash := md5.New()
	hash.Write([]byte(endpointId.Hex()))
	hash.Write([]byte(fmt.Sprintf("%d", timestamp.Unix())))

	return primitive.Binary{
		Subtype: BinaryMD5,
		Data:    hash.Sum(nil),
	}
}

func GetObj(typ string) Doc {
	switch typ {
	case "system":
		return &System{}
	case "load":
		return &Load{}
	case "disk":
		return &Disk{}
	case "diskio":
		return &DiskIo{}
	case "network":
		return &Network{}
	case "kmsg":
		return &Kmsg{}
	default:
		return nil
	}
}

func GetChart(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, typ string, start, end time.Time,
	interval time.Duration) (ChartData, error) {

	start = start.Add(time.Duration(start.UnixMilli()%
		interval.Milliseconds()) * -time.Millisecond)
	end = end.Add(time.Duration(end.UnixMilli()%
		interval.Milliseconds()) * -time.Millisecond)

	switch typ {
	case "system":
		return GetSystemChart(c, db, endpoint, start, end, interval)
	case "load":
		return GetLoadChart(c, db, endpoint, start, end, interval)
	case "disk":
		return GetDiskChart(c, db, endpoint, start, end, interval)
	case "diskio":
		return GetDiskIoChart(c, db, endpoint, start, end, interval)
	case "network":
		return GetNetworkChart(c, db, endpoint, start, end, interval)
	default:
		return nil, &errortypes.UnknownError{
			errors.New("endpoints: Unknown resource type"),
		}
	}
}

func GetLog(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, typ string) (LogData, error) {

	switch typ {
	case "kmsg":
		return GetKmsgLog(c, db, endpoint)
	default:
		return nil, &errortypes.UnknownError{
			errors.New("endpoints: Unknown resource type"),
		}
	}
}

type Chart struct {
	start    int64
	end      int64
	intv     int64
	valType  int
	data     ChartData
	curTimes map[string]int64
}

func (c *Chart) add(resource string, timestamp int64, value interface{}) {
	c.data[resource] = append(c.data[resource], &Point{
		X: timestamp,
		Y: value,
	})
}

func (c *Chart) Add(resource string, timestamp int64, value interface{}) {
	cur := c.curTimes[resource]
	if cur == 0 {
		cur = c.start - c.intv
	}

	for timestamp-c.intv > cur {
		cur += c.intv
		c.add(resource, cur, 0)
	}

	c.add(resource, timestamp, value)
	c.curTimes[resource] = timestamp
}

func (c *Chart) Export() map[string][]*Point {
	for resource, cur := range c.curTimes {
		for c.end > cur {
			cur += c.intv
			c.add(resource, cur, 0)
		}
	}

	return c.data
}

func NewChart(start, end time.Time, interval time.Duration) (chrt *Chart) {
	chrt = &Chart{
		start:    start.UnixMilli(),
		end:      end.UnixMilli(),
		intv:     interval.Milliseconds(),
		data:     ChartData{},
		curTimes: map[string]int64{},
	}

	if interval == time.Minute {
		chrt.end -= time.Minute.Milliseconds()
	}

	return
}

type Alert struct {
	Name      string
	Resource  string
	Message   string
	Level     int
	Frequency time.Duration
}

func NewAlert(resource *alert.Alert, message string) (alrt *Alert) {
	alrt = &Alert{
		Name:      resource.Name,
		Resource:  resource.Resource,
		Message:   message,
		Level:     resource.Level,
		Frequency: time.Duration(resource.Frequency) * time.Second,
	}

	return
}
