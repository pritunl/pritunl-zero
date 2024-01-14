package settings

var Acme *acme

type acme struct {
	Id           string `bson:"_id"`
	Url          string `bson:"url" default:"https://acme-v01.api.letsencrypt.org"`
	DnsRetryRate int    `bson:"dns_retry_rate" default:"3"`
	DnsTimeout   int    `bson:"dns_timeout" default:"45"`
	DnsDelay     int    `bson:"dns_delay" default:"15"`
}

func newAcme() interface{} {
	return &acme{
		Id: "acme",
	}
}

func updateAcme(data interface{}) {
	Acme = data.(*acme)
}

func init() {
	register("acme", newAcme, updateAcme)
}
