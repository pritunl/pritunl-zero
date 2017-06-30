package settings

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/utils"
)

var System *system

type system struct {
	Id                   string `bson:"_id"`
	Name                 string `bson:"name"`
	DatabaseVersion      string `bson:"database_version"`
	CookieAuthKey        []byte `bson:"cookie_auth_key"`
	CookieCryptoKey      []byte `bson:"cookie_crypto_key"`
	ProxyCookieAuthKey   []byte `bson:"proxy_cookie_auth_key"`
	ProxyCookieCryptoKey []byte `bson:"proxy_cookie_crypto_key"`
}

func newSystem() interface{} {
	return &system{
		Id: "system",
	}
}

func updateSystem(data interface{}) {
	System = data.(*system)
}

func init() {
	module := requires.New("settings.system")
	module.After("settings")

	module.Handler = func() (err error) {
		if System.Name == "" {
			db := database.GetDatabase()
			defer db.Close()

			System.Name = utils.RandName()
			err = Commit(db, System, set.NewSet("name"))
			if err != nil {
				return
			}
		}

		return
	}

	register("system", newSystem, updateSystem)
}
