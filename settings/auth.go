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
	Tenant         string        `bson:"tenant" json:"tenant"`               // azure
	ClientId       string        `bson:"client_id" json:"client_id"`         // azure
	ClientSecret   string        `bson:"client_secret" json:"client_secret"` // azure
	Domain         string        `bson:"domain" json:"domain"`               // google
	GoogleKey      string        `bson:"google_key" json:"google_key"`       // google
	GoogleEmail    string        `bson:"google_email" json:"google_email"`   // google
	IssuerUrl      string        `bson:"issuer_url" json:"issuer_url"`       // saml
	SamlUrl        string        `bson:"saml_url" json:"saml_url"`           // saml
	SamlCert       string        `bson:"saml_cert" json:"saml_cert"`         // saml
}

type auth struct {
	Id               string      `bson:"_id"`
	Server           string      `bson:"server" default:"https://auth.pritunl.com"`
	Sync             int         `bson:"sync" json:"sync" default:"1800"`
	Providers        []*Provider `bson:"providers"`
	Window           int         `bson:"window" json:"window" default:"60"`
	AdminExpire      int         `bson:"admin_expire" json:"admin_expire" default:"1440"`
	AdminMaxDuration int         `bson:"admin_max_duration" json:"admin_max_duration" default:"4320"`
	ProxyExpire      int         `bson:"proxy_expire" json:"proxy_expire" default:"1440"`
	ProxyMaxDuration int         `bson:"proxy_max_duration" json:"proxy_max_duration" default:"4320"`
	UserExpire       int         `bson:"user_expire" json:"user_expire" default:"1440"`
	UserMaxDuration  int         `bson:"user_max_duration" json:"user_max_duration" default:"4320"`
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
