package settings

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
)

func findField(obj interface{}, key string) (
	reflect.Value, reflect.StructField, bool) {

	elm := reflect.ValueOf(obj).Elem()
	typ := elm.Type()

	n := typ.NumField()
	for i := 0; i < n; i++ {
		fld := typ.Field(i)

		tag := fld.Tag.Get("bson")
		if tag == "" {
			continue
		}

		name := strings.Split(tag, ",")[0]
		if name == key {
			return elm.Field(i), fld, true
		}
	}

	return reflect.Value{}, reflect.StructField{}, false
}

func ParseValue(group, key, val string) (parsed interface{}, err error) {
	grp, ok := registry[group]
	if !ok {
		err = &errortypes.NotFoundError{
			errors.Newf("settings: Group '%s' does not exist", group),
		}
		return
	}

	_, fld, found := findField(grp.New(), key)
	if !found {
		err = &errortypes.NotFoundError{
			errors.Newf(
				"settings: Key '%s' does not exist in group '%s'",
				key, group),
		}
		return
	}

	newVal := reflect.New(fld.Type)
	err = json.Unmarshal([]byte(val), newVal.Interface())
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrapf(err,
				"settings: Invalid value for '%s' '%s', expected type %s",
				group, key, fld.Type),
		}
		return
	}

	parsed = newVal.Elem().Interface()

	return
}

func Value(group, key string) (val interface{}, err error) {
	data, err := Load(group)
	if err != nil {
		return
	}

	fld, _, found := findField(data, key)
	if !found {
		err = &errortypes.NotFoundError{
			errors.Newf(
				"settings: Key '%s' does not exist in group '%s'",
				key, group),
		}
		return
	}

	val = fld.Interface()

	return
}
