package settings

import (
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/utils"
)

var Auth *auth

const (
	SetOnInsert = "set_on_insert"
	Merge       = "merge"
	Overwrite   = "overwrite"

	Azure     = "azure"
	AuthZero  = "authzero"
	Google    = "google"
	OneLogin  = "onelogin"
	Okta      = "okta"
	JumpCloud = "jumpcloud"

	Duo       = "duo"
	OneLogin2 = "one_login"
)

type Provider struct {
	Id              bson.ObjectID `bson:"id" json:"id"`
	Type            string        `bson:"type" json:"type"`
	Label           string        `bson:"label" json:"label"`
	DefaultRoles    []string      `bson:"default_roles" json:"default_roles"`
	AutoCreate      bool          `bson:"auto_create" json:"auto_create"`
	RoleManagement  string        `bson:"role_management" json:"role_management"`
	Region          string        `bson:"region" json:"region"`                     // azure
	Tenant          string        `bson:"tenant" json:"tenant"`                     // azure
	ClientId        string        `bson:"client_id" json:"client_id"`               // azure + authzero
	ClientSecret    string        `bson:"client_secret" json:"client_secret"`       // azure + authzero
	Domain          string        `bson:"domain" json:"domain"`                     // google + authzero
	GoogleKey       string        `bson:"google_key" json:"google_key"`             // google
	GoogleEmail     string        `bson:"google_email" json:"google_email"`         // google
	JumpCloudAppId  string        `bson:"jumpcloud_app_id" json:"jumpcloud_app_id"` // jumpcloud
	JumpCloudSecret string        `bson:"jumpcloud_secret" json:"jumpcloud_secret"` // jumpcloud
	IssuerUrl       string        `bson:"issuer_url" json:"issuer_url"`             // saml
	SamlUrl         string        `bson:"saml_url" json:"saml_url"`                 // saml
	SamlCert        string        `bson:"saml_cert" json:"saml_cert"`               // saml
}

func (p *Provider) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if p.Id.IsZero() {
		p.Id = bson.NewObjectID()
	}

	p.Label = utils.FilterStr(p.Label, 32)

	switch p.Type {
	case AuthZero:
		p.Region = ""
		p.Tenant = ""
		p.GoogleKey = ""
		p.GoogleEmail = ""
		p.JumpCloudAppId = ""
		p.JumpCloudSecret = ""
		p.IssuerUrl = ""
		p.SamlUrl = ""
		p.SamlCert = ""
		break
	case Azure:
		if p.Region == "" {
			p.Region = "global2"
		}
		p.Domain = ""
		p.GoogleKey = ""
		p.GoogleEmail = ""
		p.JumpCloudAppId = ""
		p.JumpCloudSecret = ""
		p.IssuerUrl = ""
		p.SamlUrl = ""
		p.SamlCert = ""
		break
	case Google:
		p.Region = ""
		p.Tenant = ""
		p.ClientId = ""
		p.ClientSecret = ""
		p.JumpCloudAppId = ""
		p.JumpCloudSecret = ""
		p.IssuerUrl = ""
		p.SamlUrl = ""
		p.SamlCert = ""
		break
	case OneLogin:
		p.Region = ""
		p.Tenant = ""
		p.ClientId = ""
		p.ClientSecret = ""
		p.Domain = ""
		p.GoogleKey = ""
		p.GoogleEmail = ""
		p.JumpCloudAppId = ""
		p.JumpCloudSecret = ""
		break
	case Okta:
		p.Region = ""
		p.Tenant = ""
		p.ClientId = ""
		p.ClientSecret = ""
		p.Domain = ""
		p.GoogleKey = ""
		p.GoogleEmail = ""
		p.JumpCloudAppId = ""
		p.JumpCloudSecret = ""
		break
	case JumpCloud:
		p.Region = ""
		p.Tenant = ""
		p.ClientId = ""
		p.ClientSecret = ""
		p.Domain = ""
		p.GoogleKey = ""
		p.GoogleEmail = ""
		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "unknown_provider_type",
			Message: "Unknown authentication provider type",
		}
		return
	}

	switch p.RoleManagement {
	case SetOnInsert, "":
		break
	case Merge:
		break
	case Overwrite:
		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "unknown_role_management",
			Message: "Unknown role management mode",
		}
		return
	}

	return
}

type SecondaryProvider struct {
	Id             bson.ObjectID `bson:"id" json:"id"`
	Type           string        `bson:"type" json:"type"`
	Name           string        `bson:"name" json:"name"`
	Label          string        `bson:"label" json:"label"`
	DuoHostname    string        `bson:"duo_hostname" json:"duo_hostname"`         // duo
	DuoKey         string        `bson:"duo_key" json:"duo_key"`                   // duo
	DuoSecret      string        `bson:"duo_secret" json:"duo_secret"`             // duo
	OneLoginRegion string        `bson:"one_login_region" json:"one_login_region"` // onelogin
	OneLoginId     string        `bson:"one_login_id" json:"one_login_id"`         // onelogin
	OneLoginSecret string        `bson:"one_login_secret" json:"one_login_secret"` // onelogin
	OktaDomain     string        `bson:"okta_domain" json:"okta_domain"`           // okta
	OktaToken      string        `bson:"okta_token" json:"okta_token"`             // okta
	PushFactor     bool          `bson:"push_factor" json:"push_factor"`           // duo + onelogin + okta
	PhoneFactor    bool          `bson:"phone_factor" json:"phone_factor"`         // duo + onelogin + okta
	PasscodeFactor bool          `bson:"passcode_factor" json:"passcode_factor"`   // duo + onelogin + okta
	SmsFactor      bool          `bson:"sms_factor" json:"sms_factor"`             // duo + onelogin + okta
}

func (p *SecondaryProvider) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if p.Id.IsZero() {
		p.Id = bson.NewObjectID()
	}

	p.Name = utils.FilterStr(p.Name, 32)
	p.Label = utils.FilterStr(p.Label, 32)

	switch p.Type {
	case Duo:
		p.OneLoginRegion = ""
		p.OneLoginId = ""
		p.OneLoginSecret = ""
		p.OktaDomain = ""
		p.OktaToken = ""
		break
	case OneLogin2:
		p.DuoHostname = ""
		p.DuoKey = ""
		p.DuoSecret = ""
		p.OktaDomain = ""
		p.OktaToken = ""
		if p.OneLoginRegion == "" {
			p.OneLoginRegion = "us"
		}
		break
	case Okta:
		p.DuoHostname = ""
		p.DuoKey = ""
		p.DuoSecret = ""
		p.OneLoginRegion = ""
		p.OneLoginId = ""
		p.OneLoginSecret = ""
		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "unknown_secondary_provider_type",
			Message: "Unknown secondary authentication provider type",
		}
		return
	}

	return
}

type auth struct {
	Id                    string               `bson:"_id"`
	Server                string               `bson:"server" default:"https://auth.pritunl.com"`
	Sync                  int                  `bson:"sync" json:"sync" default:"1800"`
	CookieAge             int                  `bson:"cookie_age" json:"cookie_age" default:"63072000"`
	Providers             []*Provider          `bson:"providers"`
	SecondaryProviders    []*SecondaryProvider `bson:"secondary_providers"`
	FastLogin             bool                 `bson:"fast_login" json:"fast_login"`
	ForceFastUserLogin    bool                 `bson:"force_fast_user_login" json:"force_fast_user_login"`
	ForceFastServiceLogin bool                 `bson:"force_fast_service_login" json:"force_fast_service_login"`
	Window                int                  `bson:"window" json:"window" default:"60"`
	WindowLong            int                  `bson:"window_long" json:"window_long" default:"300"`
	SecondaryExpire       int                  `bson:"secondary_expire" json:"secondary_expire" default:"90"`
	AdminExpire           int                  `bson:"admin_expire" json:"admin_expire" default:"1440"`
	AdminMaxDuration      int                  `bson:"admin_max_duration" json:"admin_max_duration" default:"4320"`
	ProxyExpire           int                  `bson:"proxy_expire" json:"proxy_expire" default:"1440"`
	ProxyMaxDuration      int                  `bson:"proxy_max_duration" json:"proxy_max_duration" default:"4320"`
	UserExpire            int                  `bson:"user_expire" json:"user_expire" default:"1440"`
	UserMaxDuration       int                  `bson:"user_max_duration" json:"user_max_duration" default:"4320"`
	DisaleGeo             bool                 `bson:"disable_geo" json:"disable_geo"`
	LimiterExpire         int                  `bson:"limiter_expire" json:"limiter_expire" default:"600"`
}

func (a *auth) GetProvider(id bson.ObjectID) *Provider {
	for _, provider := range a.Providers {
		if provider.Id == id {
			return provider
		}
	}

	return nil
}

func (a *auth) GetSecondaryProvider(id bson.ObjectID) *SecondaryProvider {
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
