package sync

import (
	"github.com/Sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"time"
)

func bastionEnabled() bool {
	return node.Self.Authorities != nil && len(node.Self.Authorities) != 0
}

func bastionInit() (err error) {
	logrus.WithFields(logrus.Fields{
		"docker_image": settings.System.BastionDockerImage,
	}).Info("sync: Pulling bastion server docker image")

	_, err = utils.ExecCombinedOutputLogged(nil, "docker",
		"pull", settings.System.BastionDockerImage)
	if err != nil {
		return
	}

	return
}

func bastionRunner() {
	time.Sleep(1 * time.Second)

	for {
		time.Sleep(1 * time.Second)

		if bastionEnabled() {
			break
		}
	}

	for {
		err := bastionInit()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("sync: Failed to init bastion host")

			time.Sleep(10 * time.Second)

			continue
		}

		break
	}
}

func initBastion() {
	go bastionRunner()
}
