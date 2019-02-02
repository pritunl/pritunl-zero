package sync

import (
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/authority"
	"github.com/pritunl/pritunl-zero/bastion"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
)

func bastionEnabled() bool {
	return node.Self.Authorities != nil &&
		len(node.Self.Authorities) != 0 &&
		strings.Contains(node.Self.Type, node.Bastion)
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

func bastionSync() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	authrs := []*authority.Authority{}
	nodeAuthrs := node.Self.Authorities

	if bastionEnabled() {
		authrs, err = authority.GetMulti(db, nodeAuthrs)
		if err != nil {
			return
		}
	}

	curAuthrs := set.NewSet()
	curContainers := []string{}

	for _, authr := range authrs {
		curAuthrs.Add(authr.Id)
	}

	for _, bast := range bastion.GetAll() {
		curContainers = append(curContainers, bast.Container)

		if !curAuthrs.Contains(bast.Authority) {
			e := bast.Stop()
			if e != nil {
				logrus.WithFields(logrus.Fields{
					"error": e,
				}).Error("sync: Failed to stop bastion server")
			}
		}
	}

	containers, err := bastion.DockerGetRunning()
	if err != nil {
		return
	}

	for containerId, authrId := range containers {
		match := false
		for _, containerId2 := range curContainers {
			if bastion.DockerMatchContainer(containerId, containerId2) {
				match = true
				break
			}
		}

		if !match {
			logrus.WithFields(logrus.Fields{
				"authority_id": authrId.Hex(),
				"container_id": containerId,
			}).Info("sync: Removing unknown bastion server")

			e := bastion.DockerRemove(containerId)
			if e != nil {
				logrus.WithFields(logrus.Fields{
					"error": e,
				}).Error("sync: Failed to remove bastion server")
			}
		}
	}

	for _, authr := range authrs {
		bast := bastion.Get(authr.Id)
		if bast == nil || !bast.State() {
			bast = bastion.New(authr.Id)

			e := bast.Start(db, authr)
			if e != nil {
				logrus.WithFields(logrus.Fields{
					"error": e,
				}).Error("sync: Failed to start bastion")
			}
		} else if bast.Diff(authr) {
			e := bast.Stop()
			if e != nil {
				logrus.WithFields(logrus.Fields{
					"error": e,
				}).Error("sync: Failed to stop bastion")
			}
		}
	}

	return
}

func bastionRunner() {
	time.Sleep(1 * time.Second)

	for {
		if bastionEnabled() {
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

		time.Sleep(1 * time.Second)
	}

	for {
		time.Sleep(1 * time.Second)

		err := bastionSync()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("sync: Failed to sync bastion host")
		}
	}
}

func initBastion() {
	go bastionRunner()
}
