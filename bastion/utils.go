package bastion

import (
	"fmt"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

func DockerMatchContainer(a, b string) bool {
	if len(b) > len(a) {
		a, b = b, a
	}
	return strings.HasPrefix(a, b)
}

func DockerGetName(authrId bson.ObjectId) string {
	return fmt.Sprintf("pritunl-bastion-%s", authrId.Hex())
}
