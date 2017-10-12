package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/log"
)

func ClearLogs() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	err = log.Clear(db)
	if err != nil {
		return
	}

	logrus.Info("cmd.log: Logs cleared")

	return
}
