package endpoints

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/dropbox/godropbox/errors"
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
	Print()
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
