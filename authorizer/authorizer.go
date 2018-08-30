package authorizer

import (
	"github.com/pritunl/pritunl-zero/cookie"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/signature"
	"github.com/pritunl/pritunl-zero/user"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type Authorizer struct {
	typ  string
	cook *cookie.Cookie
	sess *session.Session
	sig  *signature.Signature
	srvc *service.Service
	usr  *user.User
}

func (a *Authorizer) IsApi() bool {
	return a.sig != nil
}

func (a *Authorizer) IsValid() bool {
	return a.sess != nil || a.sig != nil
}

func (a *Authorizer) Clear(db *database.Database, w http.ResponseWriter,
	r *http.Request) (err error) {

	a.sess = nil
	a.sig = nil

	if a.cook != nil {
		err = a.cook.Remove(db)
		if err != nil {
			return
		}
	}

	switch a.typ {
	case Admin:
		cookie.CleanAdmin(w, r)
		break
	case Proxy:
		cookie.CleanProxy(w, r)
		break
	case User:
		cookie.CleanUser(w, r)
		break
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

	if a.sess != nil {
		if a.usr != nil {
			usr = a.usr
			return
		}

		if db != nil {
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
		}

		if usr == nil {
			a.sess = nil
		} else {
			a.usr = usr
		}
	} else if a.sig != nil {
		if a.usr != nil {
			usr = a.usr
			return
		}

		if db != nil {
			usr, err = a.sig.GetUser(db)
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
		}

		if usr == nil {
			a.sig = nil
		} else {
			a.usr = usr
		}
	}

	return
}

func (a *Authorizer) ServiceId() bson.ObjectId {
	if a.srvc != nil {
		return a.srvc.Id
	}
	return ""
}

func (a *Authorizer) GetSession() *session.Session {
	return a.sess
}

func (a *Authorizer) SessionId() string {
	if a.sess != nil {
		return a.sess.Id
	}

	return ""
}
