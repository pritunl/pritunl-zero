package utils

import (
	"fmt"

	"github.com/pritunl/mongo-go-driver/bson/primitive"
)

func GetTempPath() string {
	return fmt.Sprintf("/tmp/pritunl-zero/%s", primitive.NewObjectID().Hex())
}
