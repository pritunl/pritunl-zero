package cmd

import (
	"encoding/json"
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"gopkg.in/mgo.v2/bson"
)

func Mongo() (err error) {
	mongodbUri := flag.Arg(1)

	err = config.Load()
	if err != nil {
		return
	}

	config.Config.MongoUri = mongodbUri

	err = config.Save()
	if err != nil {
		return
	}

	logrus.WithFields(logrus.Fields{
		"mongo_uri": config.Config.MongoUri,
	}).Info("cmd.settings: Set MongoDB URI")

	return
}

func ResetId() (err error) {
	err = config.Load()
	if err != nil {
		return
	}

	config.Config.NodeId = bson.NewObjectId().Hex()

	err = config.Save()
	if err != nil {
		return
	}

	logrus.WithFields(logrus.Fields{
		"node_id": config.Config.NodeId,
	}).Info("cmd.settings: Reset node ID")

	return
}

func SettingsSet() (err error) {
	group := flag.Arg(1)
	key := flag.Arg(2)
	val := flag.Arg(3)
	db := database.GetDatabase()
	defer db.Close()

	var valParsed interface{}
	err = json.Unmarshal([]byte(val), &valParsed)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "cmd.settings: Failed to parse value"),
		}
		return
	}

	err = settings.Set(db, group, key, valParsed)
	if err != nil {
		return
	}

	return
}
