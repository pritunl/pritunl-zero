package user

import (
	"github.com/dropbox/godropbox/container/set"
)

const (
	Local    = "local"
	Api      = "api"
	Azure    = "azure"
	Google   = "google"
	OneLogin = "onelogin"
	Okta     = "okta"
)

var (
	types = set.NewSet(
		Local,
		Api,
		Azure,
		Google,
		OneLogin,
		Okta,
	)
)
