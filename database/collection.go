package database

import (
	"reflect"
	"strings"

	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo"
	"github.com/pritunl/mongo-go-driver/mongo/options"
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

func (c *Collection) Upsert(id interface{}, data interface{}) (err error) {
	opts := &options.UpdateOptions{}
	opts.SetUpsert(true)

	_, err = c.UpdateOne(c.db, &bson.M{
		"_id": id,
	}, &bson.M{
		"$set": data,
	}, opts)
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
		field := val.Field(i)
		typ := val.Type().Field(i)

		if typ.PkgPath != "" {
			continue
		}

		tag := typ.Tag.Get("bson")
		if tag == "" || tag == "-" {
			continue
		}

		tag = strings.Split(tag, ",")[0]

		if fields.Contains(tag) {
			val := field.Interface()

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
		} else if (field.Kind() == reflect.Struct) ||
			(field.Kind() == reflect.Pointer &&
				field.Elem().Kind() == reflect.Struct) {

			var val reflect.Value
			if field.Kind() == reflect.Struct {
				val = field
			} else {
				val = reflect.ValueOf(field.Interface()).Elem()
			}

			x := val.NumField()
			for j := 0; j < x; j++ {
				nestedField := val.Field(j)
				nestedTyp := val.Type().Field(j)

				if nestedTyp.PkgPath != "" {
					continue
				}

				nestedTag := nestedTyp.Tag.Get("bson")
				if nestedTag == "" || nestedTag == "-" {
					continue
				}

				nestedTag = strings.Split(nestedTag, ",")[0]
				nestedTag = tag + "." + nestedTag

				if fields.Contains(nestedTag) {
					nestedVal := nestedField.Interface()

					switch nestedValTyp := nestedVal.(type) {
					case primitive.ObjectID:
						if nestedValTyp.IsZero() {
							data[nestedTag] = nil
						} else {
							data[nestedTag] = nestedVal
						}
						break
					default:
						data[nestedTag] = nestedVal
					}
				}
			}
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
		field := val.Field(i)
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

		if fields.Contains(tag) {
			val := field.Interface()

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
		} else if (field.Kind() == reflect.Struct) ||
			(field.Kind() == reflect.Pointer &&
				field.Elem().Kind() == reflect.Struct) {

			var val reflect.Value
			if field.Kind() == reflect.Struct {
				val = field
			} else {
				val = reflect.ValueOf(field.Interface()).Elem()
			}

			x := val.NumField()
			for j := 0; j < x; j++ {
				nestedField := val.Field(j)
				nestedTyp := val.Type().Field(j)

				if nestedTyp.PkgPath != "" {
					continue
				}

				nestedTag := nestedTyp.Tag.Get("bson")
				if nestedTag == "" || nestedTag == "-" {
					continue
				}

				nestedOmitempty := strings.Contains(nestedTag, "omitempty")
				nestedTag = strings.Split(nestedTag, ",")[0]
				nestedTag = tag + "." + nestedTag

				if fields.Contains(nestedTag) {
					nestedVal := nestedField.Interface()

					switch nestedValTyp := nestedVal.(type) {
					case primitive.ObjectID:
						if nestedValTyp.IsZero() {
							if nestedOmitempty {
								dataUnset[nestedTag] = 1
								dataUnseted = true
							} else {
								dataSet[nestedTag] = nil
							}
						} else {
							dataSet[nestedTag] = nestedVal
						}
						break
					default:
						dataSet[nestedTag] = nestedVal
					}
				}
			}
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
