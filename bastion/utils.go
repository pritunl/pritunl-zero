package bastion

import (
	"fmt"
	"strings"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/utils"
)

func GetRuntime() string {
	exists, _ := utils.Exists("/usr/bin/podman")
	if exists {
		return "/usr/bin/podman"
	}
	return "docker"
}

func DockerMatchContainer(a, b string) bool {
	if len(b) > len(a) {
		a, b = b, a
	}
	return strings.HasPrefix(a, b)
}

func DockerGetName(authrId bson.ObjectID) string {
	return fmt.Sprintf("pritunl-bastion-%s", authrId.Hex())
}

func DockerGetRunning() (running map[string]bson.ObjectID, err error) {
	running = map[string]bson.ObjectID{}

	output, err := utils.ExecOutput("",
		GetRuntime(), "ps", "-a", "--format", "{{.Names}}:{{.ID}}")
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

		authrId, e := bson.ObjectIDFromHex(name[16:])
		if e != nil {
			err = &errortypes.ParseError{
				errors.Wrap(e, "bastion: Failed to parse ObjectID"),
			}
			return
		}

		running[containerId] = authrId
	}

	return
}

func DockerRemove(containerId string) (err error) {
	_, err = utils.ExecOutputLogged(nil, GetRuntime(), "rm", "-f", containerId)
	if err != nil {
		return
	}

	return
}
