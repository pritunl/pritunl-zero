package utils

import (
	"fmt"

	"github.com/pritunl/mongo-go-driver/v2/bson"
)

func GetTempPath() string {
	return fmt.Sprintf("/tmp/pritunl-zero/%s", bson.NewObjectID().Hex())
}
