package setup

import (
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/Sirupsen/logrus"
)

func Setup() (err error) {
	db := database.GetDatabase()

	exists, err := user.HasSuper(db)
	if err != nil {
		return
	}

	if !exists {
		logrus.Info("setup: Creating default super user")

		usr := user.User{
			Type:          "local",
			Username:      "pritunl",
			Administrator: "super",
		}

		err = usr.SetPassword("pritunl")
		if err != nil {
			return
		}

		err = usr.Insert(db)
		if err != nil {
			return
		}
	}

	return
}
