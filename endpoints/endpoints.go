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
	Format(primitive.ObjectID)
	StaticData() *bson.M
	Print()
}

type Chart struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
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
	default:
		return nil, &errortypes.UnknownError{
			errors.New("endpoints: Unknown resource type"),
		}
	}
}
