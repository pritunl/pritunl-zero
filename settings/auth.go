package settings

import (
	"gopkg.in/mgo.v2/bson"
)

var Auth *auth

type Provider struct {
	Id             bson.ObjectId `bson:"id" json:"id"`
	Type           string        `bson:"type" json:"type"`
	Label          string        `bson:"label" json:"label"`
	DefaultRoles   []string      `bson:"default_roles" json:"default_roles"`
	AutoCreate     bool          `bson:"auto_create" json:"auto_create"`
	RoleManagement string        `bson:"role_management" json:"role_management"`
	Domain         string        `bson:"domain" json:"domain"`         // google
	IssuerUrl      string        `bson:"issuer_url" json:"issuer_url"` // saml
	SamlUrl        string        `bson:"saml_url" json:"saml_url"`     // saml
	SamlCert       string        `bson:"saml_cert" json:"saml_cert"`   // saml
}

type auth struct {
	Id          string      `bson:"_id"`
	Server      string      `bson:"server" default:"https://auth-test.pritunl.net"`
	Expire      int         `bson:"expire" json:"expire" default:"72"`
	MaxDuration int         `bson:"max_duration" json:"max_duration" default:"0"`
	Providers   []*Provider `bson:"providers"`
}

func (a *auth) GetProvider(id bson.ObjectId) *Provider {
	for _, provider := range a.Providers {
		if provider.Id == id {
			return provider
		}
	}

	return nil
}

func newAuth() interface{} {
	return &auth{
		Id:        "auth",
		Providers: []*Provider{},
	}
}

func updateAuth(data interface{}) {
	Auth = data.(*auth)
}

func init() {
	register("auth", newAuth, updateAuth)
}
