package settings

var System *system

type system struct {
	Id                             string `bson:"_id"`
	Name                           string `bson:"name"`
	DatabaseVersion                int    `bson:"database_version"`
	Demo                           bool   `bson:"demo"`
	License                        string `bson:"license"`
	CookieAuthKey                  []byte `bson:"cookie_auth_key"`
	CookieCryptoKey                []byte `bson:"cookie_crypto_key"`
	ProxyCookieAuthKey             []byte `bson:"proxy_cookie_auth_key"`
	ProxyCookieCryptoKey           []byte `bson:"proxy_cookie_crypto_key"`
	UserCookieAuthKey              []byte `bson:"user_cookie_auth_key"`
	UserCookieCryptoKey            []byte `bson:"user_cookie_crypto_key"`
	NodeTimestampTtl               int    `bson:"node_timestamp_ttl" default:"10"`
	AcmeKeyAlgorithm               string `bson:"acme_key_algorithm" default:"rsa"`
	SshPubKeyLen                   int    `bson:"ssh_pub_key_len" default:"5000"`
	SshHostTokenLen                int    `bson:"ssh_host_token_len" default:"10"`
	HsmResponseTimeout             int    `bson:"hsm_response_timeout" default:"10"`
	DisableBastionHostCertificates bool   `bson:"disable_bastion_host_certificates"`
	BastionDockerImage             string `bson:"bastion_docker_image" default:"docker.io/pritunl/pritunl-bastion"`
	ClientCertCacheTtl             int    `bson:"client_cert_cache_ttl" default:"60"`
	TwilioAccount                  string `bson:"twilio_account"`
	TwilioSecret                   string `bson:"twilio_secret"`
	TwilioNumber                   string `bson:"twilio_number"`
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
