package utils

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func GetTempPath() string {
	return fmt.Sprintf("/tmp/pritunl-zero/%s", bson.NewObjectId().Hex())
}
