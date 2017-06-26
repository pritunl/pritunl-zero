package user

import (
	"github.com/dropbox/godropbox/container/set"
)

const (
	Local    = "local"
	Google   = "google"
	OneLogin = "onelogin"
	Okta     = "okta"
)

var (
	types = set.NewSet(
		Local,
		Google,
		OneLogin,
		Okta,
	)
)
