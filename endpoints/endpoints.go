package endpoints

import (
	"crypto/md5"
	"encoding/json"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
)

type Doc interface {
	GetCollection(*database.Database) *database.Collection
	SetEndpoint(primitive.ObjectID)
	Print()
}

func GenerateId(endpoint, clientId primitive.ObjectID) []byte {
	hash := md5.New()
	hash.Write([]byte(endpoint.Hex()))
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

func ProcessDoc(docType string, docData string) (err error) {
	docObj := GetObj(docType)

	err = json.Unmarshal([]byte(docData), docObj)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "endpoints: Failed to parse doc"),
		}
		return
	}

	docObj.Print()

	return
}
