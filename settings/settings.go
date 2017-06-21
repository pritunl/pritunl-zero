// Settings stored on mongodb.
package settings

import (
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/requires"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	registry = map[string]interface{}{}
)

func Commit(db *database.Database, group interface{}, fields set.Set) (
	err error) {
	coll := db.Settings()

	selector := database.SelectFields(group, set.NewSet("_id"))
	update := database.SelectFields(group, fields)

	_, err = coll.Upsert(selector, update)
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

	err = coll.Find(bson.M{
		"_id": group,
	}).Select(bson.M{
		key: 1,
	}).One(grp)
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

	_, err = coll.Upsert(bson.M{
		"_id": group,
	}, bson.M{"$set": bson.M{
		key: val,
	}})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func parseFindError(inErr error) (err error) {
	if inErr != nil {
		switch inErr.(type) {
		case *database.NotFoundError:
			err = nil
		default:
			err = inErr
		}
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

func update(group string, data interface{}) (err error) {
	db := database.GetDatabase()
	defer db.Close()
	coll := db.Settings()

	err = parseFindError(coll.FindOneId(group, data))
	if err != nil {
		return
	}

	setDefaults(data)

	return
}

func Update(name string) {
	group, ok := registry[name]
	if !ok {
		return
	}

	for {
		err := update(name, group)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("database: Update")
		} else {
			break
		}

		time.Sleep(constants.RetryDelay)
	}
}

func register(name string, group interface{}) {
	registry[name] = group
}

func init() {
	module := requires.New("settings")
	module.After("database")

	module.Handler = func() {
		for name := range registry {
			Update(name)
		}
	}
}
