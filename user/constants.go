package user

import (
	"github.com/dropbox/godropbox/container/set"
)

const (
	Local    = "local"
	Api      = "api"
	Azure    = "azure"
	AuthZero = "authzero"
	Google   = "google"
	OneLogin = "onelogin"
	Okta     = "okta"
)

var (
	types = set.NewSet(
		Local,
		Api,
		Azure,
		AuthZero,
		Google,
		OneLogin,
		Okta,
	)
)
