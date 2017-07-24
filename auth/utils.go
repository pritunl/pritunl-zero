package auth

import (
	"github.com/pritunl/pritunl-zero/cookie"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/session"
	"net/http"
)

func Get(db *database.Database, state string) (tokn *Token, err error) {
	coll := db.Tokens()
	tokn = &Token{}

	err = coll.FindOneId(state, tokn)
	if err != nil {
		return
	}

	return
}

func CookieSession(db *database.Database,
	w http.ResponseWriter, r *http.Request) (
	cook *cookie.Cookie, sess *session.Session, err error) {

	cook, err = cookie.Get(w, r)
	if err != nil {
		return
	}

	sess, err = cook.GetSession(db)
	if err != nil {
		switch err.(type) {
		case *errortypes.NotFoundError:
			sess = nil
			err = nil
			break
		}
		return
	}

	return
}

func CookieSessionProxy(db *database.Database, srvc *service.Service,
	w http.ResponseWriter, r *http.Request) (
	cook *cookie.Cookie, sess *session.Session, err error) {

	cook, err = cookie.GetProxy(srvc, w, r)
	if err != nil {
		return
	}

	sess, err = cook.GetSession(db)
	if err != nil {
		switch err.(type) {
		case *errortypes.NotFoundError:
			sess = nil
			err = nil
			break
		}
		return
	}

	return
}
