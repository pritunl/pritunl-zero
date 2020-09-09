package settings

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/utils"
)

func Commit(db *database.Database, group interface{}, fields set.Set) (
	err error) {

	coll := db.Settings()

	selector := database.SelectFields(group, set.NewSet("_id"))
	update := database.SelectFields(group, fields)
	opts := &options.UpdateOptions{}
	opts.SetUpsert(true)

	_, err = coll.UpdateOne(
		db,
		selector,
		&bson.M{
			"$set": update,
		},
		opts,
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Get(db *database.Database, group string, key string) (
	val interface{}, err error) {

	coll := db.Settings()

	grp := map[string]interface{}{}

	err = coll.FindOne(
		db,
		&bson.M{
			"_id": group,
		},
		&options.FindOneOptions{
			Projection: &bson.D{
				{key, 1},
			},
		},
	).Decode(grp)
	if err != nil {
		err = database.ParseError(err)

		switch err.(type) {
		case *database.NotFoundError:
			err = nil
			return
		default:
			err = &errortypes.DatabaseError{
				errors.Wrap(err, "settings: Database error"),
			}
			return
		}
	}

	val = grp[key]
	return
}

func Set(db *database.Database, group string, key string, val interface{}) (
	err error) {

	coll := db.Settings()
	opts := &options.UpdateOptions{}
	opts.SetUpsert(true)

	_, err = coll.UpdateOne(
		db,
		&bson.M{
			"_id": group,
		},
		&bson.M{
			"$set": &bson.M{
				key: val,
			},
		},
		opts,
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Unset(db *database.Database, group string, key string) (
	err error) {

	coll := db.Settings()
	opts := &options.UpdateOptions{}
	opts.SetUpsert(true)

	_, err = coll.UpdateOne(
		db,
		&bson.M{
			"_id": group,
		},
		&bson.M{
			"$unset": &bson.M{
				key: "",
			},
		},
		opts,
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func setDefaults(obj interface{}) {
	val := reflect.ValueOf(obj)
	elm := val.Elem()

	n := elm.NumField()
	for i := 0; i < n; i++ {
		fld := elm.Field(i)
		typ := elm.Type().Field(i)

		if typ.PkgPath != "" {
			continue
		}

		tag := typ.Tag.Get("default")
		if tag == "" {
			continue
		}

		switch fld.Kind() {
		case reflect.Bool:
			if !fld.IsNil() {
				break
			}

			parVal, err := strconv.ParseBool(tag)
			if err != nil {
				panic(err)
			}
			fld.SetBool(parVal)
		case reflect.Int:
			if fld.Int() != 0 {
				break
			}

			parVal, err := strconv.Atoi(tag)
			if err != nil {
				panic(err)
			}
			fld.SetInt(int64(parVal))
		case reflect.String:
			if fld.String() != "" {
				break
			}

			fld.SetString(tag)
		case reflect.Slice:
			if fld.Len() != 0 {
				break
			}

			sliceType := reflect.TypeOf(fld.Interface()).Elem()
			vals := strings.Split(tag, ",")
			n := len(vals)
			slice := reflect.MakeSlice(reflect.SliceOf(sliceType), n, n)

			switch sliceType.Kind() {
			case reflect.Bool:
				for i, val := range vals {
					parVal, err := strconv.ParseBool(val)
					if err != nil {
						panic(err)
					}
					slice.Index(i).SetBool(parVal)
				}
			case reflect.Int:
				for i, val := range vals {
					parVal, err := strconv.Atoi(val)
					if err != nil {
						panic(err)
					}
					slice.Index(i).SetInt(int64(parVal))
				}
			case reflect.String:
				for i, val := range vals {
					slice.Index(i).SetString(val)
				}
			}

			fld.Set(slice)
		}
	}

	return
}

func Update(name string) (err error) {
	db := database.GetDatabase()
	defer db.Close()

	coll := db.Settings()
	group := registry[name]
	data := group.New()

	err = database.IgnoreNotFoundError(coll.FindOneId(name, data))
	if err != nil {
		return
	}

	setDefaults(data)

	group.Update(data)

	return
}

func update() {
	for {
		time.Sleep(10 * time.Second)
		for name := range registry {
			err := Update(name)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("settings: Update error")
				return
			}
		}
	}
}

func init() {
	module := requires.New("settings")
	module.After("database")

	module.Handler = func() (err error) {
		for name := range registry {
			err = Update(name)
			if err != nil {
				return
			}
		}

		db := database.GetDatabase()
		defer db.Close()

		if System.DatabaseVersion == 0 {
			System.DatabaseVersion = constants.DatabaseVersion
			err = Commit(db, System, set.NewSet("database_version"))
			if err != nil {
				return
			}
		}

		if System.DatabaseVersion > constants.DatabaseVersion {
			logrus.WithFields(logrus.Fields{
				"database_version": System.DatabaseVersion,
				"software_version": constants.DatabaseVersion,
			}).Error("settings: Database version newer then software")

			err = &errortypes.DatabaseError{
				errors.New(
					"settings: Database version newer then software"),
			}
			return
		} else if System.DatabaseVersion != constants.DatabaseVersion {
			logrus.WithFields(logrus.Fields{
				"database_version":     System.DatabaseVersion,
				"new_database_version": constants.DatabaseVersion,
			}).Info("settings: Upgrading database version")

			System.DatabaseVersion = constants.DatabaseVersion
			err = Commit(db, System, set.NewSet("database_version"))
			if err != nil {
				return
			}
		}

		if System.Name == "" {
			System.Name = utils.RandName()
			err = Commit(db, System, set.NewSet("name"))
			if err != nil {
				return
			}
		}

		if Auth.Providers == nil {
			Auth.Providers = []*Provider{}
			err = Commit(db, Auth, set.NewSet("providers"))
			if err != nil {
				return
			}
		}
		if Auth.SecondaryProviders == nil {
			Auth.SecondaryProviders = []*SecondaryProvider{}
			err = Commit(db, Auth, set.NewSet("secondary_providers"))
			if err != nil {
				return
			}
		}

		go update()

		return
	}
}
