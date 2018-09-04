package config

import (
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"os"
	"time"
)

var (
	Config            = &ConfigData{}
	StaticRoot        = ""
	StaticTestingRoot = ""
)

type ConfigData struct {
	path     string `json:"-"`
	loaded   bool   `json:"-"`
	MongoUri string `json:"mongo_uri"`
	NodeId   string `json:"node_id"`
}

func (c *ConfigData) Save() (err error) {
	if !c.loaded {
		err = &errortypes.WriteError{
			errors.New("config: Config file has not been loaded"),
		}
		return
	}

	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "config: File marshal error"),
		}
		return
	}

	err = ioutil.WriteFile(constants.ConfPath, data, 0600)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "config: File write error"),
		}
		return
	}

	return
}

func Load() (err error) {
	data := &ConfigData{}

	_, err = os.Stat(constants.ConfPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
			data.loaded = true
			Config = data
		} else {
			err = &errortypes.ReadError{
				errors.Wrap(err, "config: File stat error"),
			}
		}
		return
	}

	file, err := ioutil.ReadFile(constants.ConfPath)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "config: File read error"),
		}
		return
	}

	err = json.Unmarshal(file, data)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "config: File unmarshal error"),
		}
		return
	}

	data.loaded = true

	Config = data

	return
}

func Save() (err error) {
	err = Config.Save()
	if err != nil {
		return
	}

	return
}

func GetModTime() (mod time.Time, err error) {
	stat, err := os.Stat(constants.ConfPath)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "config: Failed to stat conf file"),
		}
		return
	}

	mod = stat.ModTime()

	return
}

func init() {
	module := requires.New("config")

	module.Handler = func() (err error) {
		for _, pth := range constants.StaticRoot {
			exists, _ := utils.ExistsDir(pth)
			if exists {
				StaticRoot = pth
			}
		}
		if StaticRoot == "" {
			StaticRoot = constants.StaticRoot[len(constants.StaticRoot)-1]
		}

		for _, pth := range constants.StaticTestingRoot {
			exists, _ := utils.ExistsDir(pth)
			if exists {
				StaticTestingRoot = pth
			}
		}
		if StaticTestingRoot == "" {
			StaticTestingRoot = constants.StaticTestingRoot[len(
				constants.StaticTestingRoot)-1]
		}

		err = utils.ExistsMkdir(constants.TempPath, 0700)
		if err != nil {
			return
		}

		err = Load()
		if err != nil {
			return
		}

		save := false
		nodeId := os.Getenv("NODE_ID")
		mongoUri := os.Getenv("MONGO_URI")

		if Config.NodeId == "" && nodeId == "" {
			save = true
			Config.NodeId = bson.NewObjectId().Hex()
		}

		if Config.MongoUri == "" && mongoUri == "" {
			save = true
			Config.MongoUri = "mongodb://localhost:27017/pritunl-zero"
		}

		if save {
			err = Save()
			if err != nil {
				return
			}
		}

		if nodeId != "" {
			Config.NodeId = nodeId
		}
		if mongoUri == "" {
			Config.MongoUri = mongoUri
		}

		return
	}
}
