package authorizer

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/auth"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/signature"
	"github.com/pritunl/pritunl-zero/user"
	"net/http"
)

func Authorize(db *database.Database, w http.ResponseWriter,
	r *http.Request) (authr *Authorizer, err error) {

	token := r.Header.Get("Pritunl-Zero-Token")
	sigStr := r.Header.Get("Pritunl-Zero-Signature")

	if token != "" && sigStr != "" {
		timestamp := r.Header.Get("Pritunl-Zero-Timestamp")
		nonce := r.Header.Get("Pritunl-Zero-Nonce")

		sig, e := signature.Parse(
			token,
			sigStr,
			timestamp,
			nonce,
			r.Method,
			r.URL.Path,
		)
		if e != nil {
			err = e
			return
		}

		err = sig.Validate(db)
		if err != nil {
			return
		}

		authr = &Authorizer{
			isProxy: false,
			sig:     sig,
		}
	} else {
		cook, sess, e := auth.CookieSession(db, w, r)
		if e != nil {
			err = e
			return
		}

		authr = &Authorizer{
			isProxy: false,
			cook:    cook,
			sess:    sess,
		}
	}

	return
}

func AuthorizeProxy(db *database.Database, srvc *service.Service,
	w http.ResponseWriter, r *http.Request) (authr *Authorizer, err error) {

	token := r.Header.Get("Pritunl-Zero-Token")
	sigStr := r.Header.Get("Pritunl-Zero-Signature")

	if token != "" && sigStr != "" {
		timestamp := r.Header.Get("Pritunl-Zero-Timestamp")
		nonce := r.Header.Get("Pritunl-Zero-Nonce")

		sig, e := signature.Parse(
			token,
			sigStr,
			timestamp,
			nonce,
			r.Method,
			r.URL.Path,
		)
		if e != nil {
			err = e
			return
		}

		authr = &Authorizer{
			isProxy: true,
			sig:     sig,
		}
		return
	}

	cook, sess, err := auth.CookieSessionProxy(db, srvc, w, r)
	if err != nil {
		return
	}

	authr = &Authorizer{
		isProxy: true,
		cook:    cook,
		sess:    sess,
	}

	return
}

func New() (authr *Authorizer) {
	authr = &Authorizer{
		isProxy: false,
	}

	return
}

func NewProxy() (authr *Authorizer) {
	authr = &Authorizer{
		isProxy: true,
	}

	return
}

func ValidateAdmin(db *database.Database, usr *user.User) (
	errData *errortypes.ErrorData, err error) {

	if usr.Disabled || usr.Administrator != "super" {
		errData = &errortypes.ErrorData{
			Error:   "unauthorized",
			Message: "Not authorized",
		}
		return
	}

	return
}

func Validate(db *database.Database, usr *user.User,
	authr *Authorizer, srvc *service.Service,
	r *http.Request) (errData *errortypes.ErrorData, err error) {

	if usr.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "unauthorized",
			Message: "Not authorized",
		}
		return
	}

	usrRoles := set.NewSet()
	for _, role := range usr.Roles {
		usrRoles.Add(role)
	}

	roleMatch := false
	for _, role := range srvc.Roles {
		if usrRoles.Contains(role) {
			roleMatch = true
			break
		}
	}

	if !roleMatch {
		errData = &errortypes.ErrorData{
			Error:   "service_unauthorized",
			Message: "Not authorized for service",
		}
		return
	}

	if !authr.IsApi() {
		polices, e := policy.GetService(db, srvc.Id)
		if e != nil {
			err = e
			return
		}

		for _, polcy := range polices {
			errData, err = polcy.ValidateUser(db, usr, r)
			if err != nil || errData != nil {
				return
			}
		}

		polices, err = policy.GetRoles(db, usr.Roles)
		if err != nil {
			return
		}

		for _, polcy := range polices {
			errData, err = polcy.ValidateUser(db, usr, r)
			if err != nil || errData != nil {
				return
			}
		}
	}

	return
}
