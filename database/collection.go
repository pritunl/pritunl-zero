package database

import (
	"github.com/dropbox/godropbox/container/set"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
)

type Collection struct {
	mgo.Collection
	Database *Database
}

func (c *Collection) FindOne(query interface{}, result interface{}) (
	err error) {

	err = c.Find(query).One(result)
	if err != nil {
		err = ParseError(err)
		return
	}

	return
}

func (c *Collection) FindOneId(id interface{}, result interface{}) (
	err error) {

	err = c.FindId(id).One(result)
	if err != nil {
		err = ParseError(err)
		return
	}

	return
}

func (c *Collection) Commit(id interface{}, data interface{}) (err error) {
	err = c.UpdateId(id, bson.M{
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

	err = c.UpdateId(id, bson.M{
		"$set": SelectFields(data, fields),
	})
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

		data[tag] = val
	}

	return
}
