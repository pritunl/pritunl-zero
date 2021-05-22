package endpoints

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
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
}

type ChartFloat struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
}

type ChartUint struct {
	X int64  `json:"x"`
	Y uint64 `json:"y"`
}

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
	case "network":
		return &Network{}
	default:
		return nil
	}
}

func GetChart(c context.Context, db *database.Database,
	endpoint primitive.ObjectID, typ string, start, end time.Time,
	interval time.Duration) (interface{}, error) {

	switch typ {
	case "system":
		return GetSystemChart(c, db, endpoint, start, end, interval)
	case "load":
		return GetLoadChart(c, db, endpoint, start, end, interval)
	case "disk":
		return GetDiskChart(c, db, endpoint, start, end, interval)
	case "network":
		return GetNetworkChart(c, db, endpoint, start, end, interval)
	default:
		return nil, &errortypes.UnknownError{
			errors.New("endpoints: Unknown resource type"),
		}
	}
}
