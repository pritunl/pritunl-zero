package task

import (
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/notification"
	"github.com/sirupsen/logrus"
)

var notificationCheck = &Task{
	Name:    "notification_check",
	Version: 1,
	Hours: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12,
		13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
	Minutes:    []int{15},
	Handler:    notificationCheckHandler,
	RunOnStart: true,
}

func notificationCheckHandler(db *database.Database) (err error) {
	err = notification.Check()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("task: Failed to check vulnerability alerts")
	}

	return
}

func init() {
	register(notificationCheck)
}
