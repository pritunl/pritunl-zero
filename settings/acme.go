package settings

var Acme *acme

type acme struct {
	Id  string `bson:"_id"`
	Url string `bson:"url" default:"https://acme-v01.api.letsencrypt.org"`
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
