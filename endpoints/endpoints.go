package endpoints

import (
	"crypto/md5"

	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
)

type Doc interface {
	GetCollection(*database.Database) *database.Collection
	Format(primitive.ObjectID)
	Print()
}

func GenerateId(endpointId, clientId primitive.ObjectID) []byte {
	hash := md5.New()
	hash.Write([]byte(endpointId.Hex()))
	hash.Write([]byte(clientId.Hex()))
	return hash.Sum(nil)
}

func GetObj(typ string) Doc {
	switch typ {
	case "system":
		return &System{}
	default:
		return nil
	}
}
