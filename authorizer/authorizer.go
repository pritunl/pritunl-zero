package authorizer

import (
	"github.com/pritunl/pritunl-zero/cookie"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/user"
	"net/http"
)

type Authorizer struct {
	isProxy bool
	cook    *cookie.Cookie
	sess    *session.Session
}

func (a *Authorizer) IsValid() bool {
	return a.sess != nil
}

func (a *Authorizer) Clear(db *database.Database, w http.ResponseWriter,
	r *http.Request) (err error) {

	a.sess = nil

	if a.cook != nil {
		err = a.cook.Remove(db)
		if err != nil {
			return
		}
	}

	if a.isProxy {
		cookie.CleanProxy(w, r)
	} else {
		cookie.Clean(w, r)
	}

	return
}

func (a *Authorizer) Remove(db *database.Database) error {
	if a.sess == nil {
		return nil
	}

	return a.sess.Remove(db)
}

func (a *Authorizer) GetUser(db *database.Database) (
	usr *user.User, err error) {

	if a.sess == nil {
		return
	}

	usr, err = a.sess.GetUser(db)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			usr = nil
			err = nil
			break
		default:
			return
		}
	}

	if usr == nil {
		a.sess = nil
	}

	return
}

func (a *Authorizer) SessionId() string {
	if a.sess != nil {
		return a.sess.Id
	}

	return ""
}
