package bastion

import (
	"fmt"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

func DockerGetName(authrId bson.ObjectId) string {
	return fmt.Sprintf("pritunl-bastion-%s", authrId.Hex())
}
