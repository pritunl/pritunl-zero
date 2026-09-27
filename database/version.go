package database

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/sirupsen/logrus"
)

type systemVersion struct {
	Id              string `bson:"_id"`
	DatabaseVersion int    `bson:"database_version"`
}

func GetDatabaseVersion(db *Database) (version int, err error) {
	coll := db.Settings()
	doc := &systemVersion{}

	err = coll.FindOne(
		db,
		&bson.M{
			"_id": "system",
		},
		FindOneProject("database_version"),
	).Decode(doc)
	if err != nil {
		err = ParseError(err)
		if _, ok := err.(*NotFoundError); ok {
			err = nil
		}
		return
	}

	version = doc.DatabaseVersion
	return
}

func SetDatabaseVersion(db *Database, version int) (err error) {
	coll := db.Settings()

	_, err = coll.UpdateOne(
		db,
		&bson.M{
			"_id": "system",
		},
		&bson.M{
			"$set": &bson.M{
				"database_version": version,
			},
		},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		err = ParseError(err)
		return
	}

	return
}

func CheckDatabaseVersion() (err error) {
	db := GetDatabase()
	defer db.Close()

	dbVersion, err := GetDatabaseVersion(db)
	if err != nil {
		return
	}

	if dbVersion == 0 {
		logrus.WithFields(logrus.Fields{
			"database_version": constants.DatabaseVersion,
		}).Info("database: Setting database version")

		err = SetDatabaseVersion(db, constants.DatabaseVersion)
		if err != nil {
			return
		}

		return
	}

	if dbVersion > constants.DatabaseVersion {
		logrus.WithFields(logrus.Fields{
			"database_version": dbVersion,
			"software_version": constants.DatabaseVersion,
		}).Error("database: Database version newer than software version")

		err = &VersionError{
			errors.Newf(
				"database: Database version %d newer than software version %d",
				dbVersion, constants.DatabaseVersion,
			),
		}
		return
	}

	if dbVersion < constants.DatabaseVersion {
		logrus.WithFields(logrus.Fields{
			"database_version":     dbVersion,
			"new_database_version": constants.DatabaseVersion,
		}).Info("database: Upgrading database version")

		err = SetDatabaseVersion(db, constants.DatabaseVersion)
		if err != nil {
			return
		}
	}

	return
}
