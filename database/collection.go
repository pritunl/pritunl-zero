package database

import (
	"reflect"
	"strings"

	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo"
)

type Collection struct {
	db *Database
	*mongo.Collection
}

func (c *Collection) FindOneId(id interface{}, data interface{}) (err error) {
	err = c.FindOne(c.db, &bson.M{
		"_id": id,
	}).Decode(data)
	if err != nil {
		err = ParseError(err)
		return
	}

	return
}

func (c *Collection) UpdateId(id interface{}, data interface{}) (err error) {
	_, err = c.UpdateOne(c.db, &bson.M{
		"_id": id,
	}, data)
	if err != nil {
		err = ParseError(err)
		return
	}

	return
}

func (c *Collection) Commit(id interface{}, data interface{}) (err error) {
	_, err = c.UpdateOne(c.db, &bson.M{
		"_id": id,
	}, &bson.M{
		"$set": data,
	})
	if err != nil {
		err = ParseError(err)
		return
	}

	return
}

func (c *Collection) CommitFields(id interface{}, data interface{},
	fields set.Set) (err error) {

	_, err = c.UpdateOne(c.db, &bson.M{
		"_id": id,
	}, SelectFieldsAll(data, fields))
	if err != nil {
		err = ParseError(err)
		return
	}

	return
}

func SelectFields(obj interface{}, fields set.Set) (data bson.M) {
	val := reflect.ValueOf(obj).Elem()
	data = bson.M{}

	n := val.NumField()
	for i := 0; i < n; i++ {
		typ := val.Type().Field(i)

		if typ.PkgPath != "" {
			continue
		}

		tag := typ.Tag.Get("bson")
		if tag == "" || tag == "-" {
			continue
		}
		tag = strings.Split(tag, ",")[0]

		if !fields.Contains(tag) {
			continue
		}

		val := val.Field(i).Interface()

		switch valTyp := val.(type) {
		case primitive.ObjectID:
			if valTyp.IsZero() {
				data[tag] = nil
			} else {
				data[tag] = val
			}
			break
		default:
			data[tag] = val
		}
	}

	return
}

func SelectFieldsAll(obj interface{}, fields set.Set) (data bson.M) {
	val := reflect.ValueOf(obj).Elem()

	dataSet := bson.M{}
	dataUnset := bson.M{}
	dataUnseted := false

	n := val.NumField()
	for i := 0; i < n; i++ {
		typ := val.Type().Field(i)

		if typ.PkgPath != "" {
			continue
		}

		tag := typ.Tag.Get("bson")
		if tag == "" || tag == "-" {
			continue
		}

		omitempty := strings.Contains(tag, "omitempty")

		tag = strings.Split(tag, ",")[0]
		if !fields.Contains(tag) {
			continue
		}

		val := val.Field(i).Interface()

		switch valTyp := val.(type) {
		case primitive.ObjectID:
			if valTyp.IsZero() {
				if omitempty {
					dataUnset[tag] = 1
					dataUnseted = true
				} else {
					dataSet[tag] = nil
				}
			} else {
				dataSet[tag] = val
			}
			break
		default:
			dataSet[tag] = val
		}
	}

	data = bson.M{
		"$set": dataSet,
	}
	if dataUnseted {
		data["$unset"] = dataUnset
	}

	return
}
