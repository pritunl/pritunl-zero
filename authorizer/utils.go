package authorizer

import (
	"github.com/pritunl/pritunl-zero/auth"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/signature"
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
