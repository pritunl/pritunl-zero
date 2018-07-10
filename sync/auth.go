package sync

import (
	"github.com/Sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func authSync() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	coll := db.Users()

	count, err := coll.Find(&bson.M{
		"type": user.Local,
	}).Limit(1).Count()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	settings.Local.HasLocalAuth = count > 0

	return
}

func authRunner() {
	time.Sleep(1 * time.Second)

	for {
		time.Sleep(10 * time.Second)

		err := authSync()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("sync: Failed to sync authentication status")
		}
	}
}

func initAuth() {
	go authRunner()
}
