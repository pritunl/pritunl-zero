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

type SecondaryProvider struct {
	Id             bson.ObjectId `bson:"id" json:"id"`
	Type           string        `bson:"type" json:"type"`
	Name           string        `bson:"name" json:"name"`
	Label          string        `bson:"label" json:"label"`
	DuoHostname    string        `bson:"duo_hostname" json:"duo_hostname"`         // duo
	DuoKey         string        `bson:"duo_key" json:"duo_key"`                   // duo
	DuoSecret      string        `bson:"duo_secret" json:"duo_secret"`             // duo
	OneLoginRegion string        `bson:"one_login_region" json:"one_login_region"` // onelogin
	OneLoginId     string        `bson:"one_login_id" json:"one_login_id"`         // onelogin
	OneLoginSecret string        `bson:"one_login_secret" json:"one_login_secret"` // onelogin
	PushFactor     bool          `bson:"push_factor" json:"push_factor"`           // duo + onelogin
	PhoneFactor    bool          `bson:"phone_factor" json:"phone_factor"`         // duo + onelogin
	PasscodeFactor bool          `bson:"passcode_factor" json:"passcode_factor"`   // duo + onelogin
	SmsFactor      bool          `bson:"sms_factor" json:"sms_factor"`             // duo + onelogin
}

type auth struct {
	Id                 string               `bson:"_id"`
	Server             string               `bson:"server" default:"https://auth.pritunl.com"`
	Sync               int                  `bson:"sync" json:"sync" default:"1800"`
	Providers          []*Provider          `bson:"providers"`
	SecondaryProviders []*SecondaryProvider `bson:"secondary_providers"`
	Window             int                  `bson:"window" json:"window" default:"60"`
	AdminExpire        int                  `bson:"admin_expire" json:"admin_expire" default:"1440"`
	AdminMaxDuration   int                  `bson:"admin_max_duration" json:"admin_max_duration" default:"4320"`
	ProxyExpire        int                  `bson:"proxy_expire" json:"proxy_expire" default:"1440"`
	ProxyMaxDuration   int                  `bson:"proxy_max_duration" json:"proxy_max_duration" default:"4320"`
	UserExpire         int                  `bson:"user_expire" json:"user_expire" default:"1440"`
	UserMaxDuration    int                  `bson:"user_max_duration" json:"user_max_duration" default:"4320"`
}

func (a *auth) GetProvider(id bson.ObjectId) *Provider {
	for _, provider := range a.Providers {
		if provider.Id == id {
			return provider
		}
	}

	return nil
}

func (a *auth) GetSecondaryProvider(id bson.ObjectId) *SecondaryProvider {
	for _, provider := range a.SecondaryProviders {
		if provider.Id == id {
			return provider
		}
	}

	return nil
}

func newAuth() interface{} {
	return &auth{
		Id:                 "auth",
		Providers:          []*Provider{},
		SecondaryProviders: []*SecondaryProvider{},
	}
}

func updateAuth(data interface{}) {
	Auth = data.(*auth)
}

func init() {
	register("auth", newAuth, updateAuth)
}
