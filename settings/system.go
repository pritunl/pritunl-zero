package settings

var System *system

type system struct {
	Id                   string `bson:"_id"`
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
	register("system", newSystem, updateSystem)
}
