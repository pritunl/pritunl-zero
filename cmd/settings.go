package cmd

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
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
	}).Info("cmd: Set MongoDB URI")

	return
}

func ResetId() (err error) {
	err = config.Load()
	if err != nil {
		return
	}

	config.Config.NodeId = primitive.NewObjectID().Hex()

	err = config.Save()
	if err != nil {
		return
	}

	logrus.WithFields(logrus.Fields{
		"node_id": config.Config.NodeId,
	}).Info("cmd: Reset node ID")

	return
}

func DefaultPassword() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	usr, err := user.GetUsername(db, user.Local, "pritunl")
	if err != nil {
		return
	}

	if usr.DefaultPassword == "" {
		err = &errortypes.NotFoundError{
			errors.New("cmd: No default password available"),
		}
		return
	}

	logrus.Info("cmd: Get default password")

	fmt.Println("Username: pritunl")
	fmt.Println("Password: " + usr.DefaultPassword)

	return
}

func ResetPassword() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	coll := db.Users()

	_, err = coll.DeleteOne(db, &bson.M{
		"username": "pritunl",
	})
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			err = nil
		} else {
			return
		}
	}

	usr := user.User{
		Type:          user.Local,
		Username:      "pritunl",
		Administrator: "super",
	}

	_, err = usr.Validate(db)
	if err != nil {
		return
	}

	err = usr.GenerateDefaultPassword()
	if err != nil {
		return
	}

	err = usr.Insert(db)
	if err != nil {
		return
	}

	logrus.Info("cmd: Password reset")

	fmt.Println("Username: pritunl")
	fmt.Println("Password: " + usr.DefaultPassword)

	return
}

func DisablePolicies() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	coll := db.Policies()

	_, err = coll.UpdateMany(db, &bson.M{}, &bson.M{
		"$set": &bson.M{
			"disabled": true,
		},
	})
	if err != nil {
		if _, ok := err.(*database.NotFoundError); ok {
			err = nil
		} else {
			return
		}
	}

	logrus.Info("cmd: Policies disabled")

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

func SettingsUnset() (err error) {
	group := flag.Arg(1)
	key := flag.Arg(2)
	db := database.GetDatabase()
	defer db.Close()

	err = settings.Unset(db, group, key)
	if err != nil {
		return
	}

	return
}
