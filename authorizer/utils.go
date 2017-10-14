package authorizer

import (
	"github.com/pritunl/pritunl-zero/auth"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/service"
	"net/http"
)

func Authorize(db *database.Database, w http.ResponseWriter,
	r *http.Request) (authr *Authorizer, err error) {

	cook, sess, err := auth.CookieSession(db, w, r)
	if err != nil {
		return
	}

	authr = &Authorizer{
		isProxy: false,
		cook:    cook,
		sess:    sess,
	}

	return
}

func AuthorizeProxy(db *database.Database, srvc *service.Service,
	w http.ResponseWriter, r *http.Request) (authr *Authorizer, err error) {

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
