package user

import (
	"github.com/dropbox/godropbox/container/set"
)

const (
	Local    = "local"
	Azure    = "azure"
	Google   = "google"
	OneLogin = "onelogin"
	Okta     = "okta"
)

var (
	types = set.NewSet(
		Local,
		Azure,
		Google,
		OneLogin,
		Okta,
	)
)
