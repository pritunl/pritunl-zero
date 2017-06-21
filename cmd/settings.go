package cmd

import (
	"flag"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/settings"
	"strconv"
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

	return
}

func SettingsSet() (err error) {
	group := flag.Arg(1)
	key := flag.Arg(2)
	val := flag.Arg(3)
	db := database.GetDatabase()
	defer db.Close()

	var valParsed interface{}
	if x, err := strconv.Atoi(val); err == nil {
		valParsed = x
	} else {
		valParsed = val
	}

	err = settings.Set(db, group, key, valParsed)
	if err != nil {
		return
	}

	return
}
