package settings

const (
	SetOnInsert = "set_on_insert"
	Merge       = "merge"
	Overwrite   = "overwrite"
)

var System *system

type system struct {
	Id                   string `bson:"_id"`
	Name                 string `bson:"name"`
	DatabaseVersion      int    `bson:"database_version"`
	License              string `bson:"license"`
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
	register("system", newSystem, updateSystem)
}
