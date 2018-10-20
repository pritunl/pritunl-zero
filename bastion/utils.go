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

func DockerGetRunning() (running map[string]bson.ObjectId, err error) {
	running = map[string]bson.ObjectId{}

	output, err := utils.ExecOutput("",
		"docker", "ps", "-a", "--format", "{{.Names}}:{{.ID}}")
	if err != nil {
		return
	}

	for _, line := range strings.Split(output, "\n") {
		fields := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(fields) != 2 {
			continue
		}

		name := fields[0]
		containerId := fields[1]

		if len(name) != 40 || !strings.HasPrefix(name, "pritunl-bastion-") {
			continue
		}

		authrId := bson.ObjectIdHex(name[16:])

		running[containerId] = authrId
	}

	return
}

func DockerRemove(containerId string) (err error) {
	_, err = utils.ExecOutputLogged(nil, "docker", "rm", "-f", containerId)
	if err != nil {
		return
	}

	return
}
