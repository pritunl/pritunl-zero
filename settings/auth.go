package settings

import (
	"gopkg.in/mgo.v2/bson"
)

var Auth *auth

type Provider struct {
	Id           bson.ObjectId `bson:"id" json:"id"`
	Type         string        `bson:"type" json:"type"`
	Label        string        `bson:"label" json:"label"`
	DefaultRoles []string      `bson:"default_roles" json:"default_roles"`
	AutoCreate   bool          `bson:"auto_create" json:"auto_create"`
	Domain       string        `bson:"domain" json:"domain"` // google
}

type auth struct {
	Id        string      `bson:"_id"`
	Server    string      `bson:"server" default:"https://auth-test.pritunl.net"`
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
