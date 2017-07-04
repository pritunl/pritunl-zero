package settings

var Auth *auth

type Provider struct {
	Type         string   `bson:"type" json:"type"`
	Label        string   `bson:"label" json:"label"`
	DefaultRoles []string `bson:"default_roles" json:"default_roles"`
	Domain       string   `bson:"domain" json:"domain"`
}

type auth struct {
	Id        string      `bson:"_id"`
	Providers []*Provider `bson:"providers"`
}

func newAuth() interface{} {
	return &auth{
		Id: "auth",
	}
}

func updateAuth(data interface{}) {
	Auth = data.(*auth)
}

func init() {
	register("auth", newAuth, updateAuth)
}
