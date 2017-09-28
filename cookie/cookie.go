package cookie

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/settings"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

var (
	Store *sessions.CookieStore
)

type Cookie struct {
	Id    bson.ObjectId
	store *sessions.Session
	w     http.ResponseWriter
	r     *http.Request
}

func (c *Cookie) Get(key string) string {
	valInf := c.store.Values[key]
	if valInf == nil {
		return ""
	}
	return valInf.(string)
}

func (c *Cookie) Set(key string, val string) {
	c.store.Values[key] = val
}

func (c *Cookie) GetSession(db *database.Database, r *http.Request) (
	sess *session.Session, err error) {

	sessId := c.Get("id")
	if sessId == "" {
		err = &errortypes.NotFoundError{
			errors.New("cookie: Session not found"),
		}
		return
	}

	sess, err = session.GetUpdate(db, sessId, r)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			err = &errortypes.NotFoundError{
				errors.New("cookie: Session not found"),
			}
		default:
			err = &errortypes.UnknownError{
				errors.Wrap(err, "cookie: Unknown session error"),
			}
		}
		return
	}

	return
}

func (c *Cookie) NewSession(db *database.Database, r *http.Request,
	id bson.ObjectId, remember bool) (sess *session.Session, err error) {

	sess, err = session.New(db, r, id)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "cookie: Unknown session error"),
		}
		return
	}

	c.Set("id", sess.Id)
	maxAge := 0

	if remember {
		maxAge = 15778500
	}

	c.store.Options.MaxAge = maxAge

	err = c.Save()
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "cookie: Unknown session error"),
		}
		return
	}

	return
}

func (c *Cookie) Remove(db *database.Database) (err error) {
	sessId := c.Get("id")
	if sessId == "" {
		return
	}

	err = session.Remove(db, sessId)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "cookie: Unknown session error"),
		}
		return
	}

	c.Set("id", "")
	err = c.Save()
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "cookie: Unknown session error"),
		}
		return
	}

	return
}

func (c *Cookie) Save() (err error) {
	err = c.store.Save(c.r, c.w)
	return
}

func init() {
	module := requires.New("cookie")
	module.After("settings")

	module.Handler = func() (err error) {
		cookieAuthKey := settings.System.CookieAuthKey
		cookieCryptoKey := settings.System.CookieCryptoKey
		proxyCookieAuthKey := settings.System.ProxyCookieAuthKey
		proxyCookieCryptoKey := settings.System.ProxyCookieCryptoKey

		if len(cookieAuthKey) == 0 || len(cookieCryptoKey) == 0 ||
			len(proxyCookieAuthKey) == 0 || len(proxyCookieCryptoKey) == 0 {

			db := database.GetDatabase()
			defer db.Close()

			cookieAuthKey = securecookie.GenerateRandomKey(64)
			cookieCryptoKey = securecookie.GenerateRandomKey(32)
			proxyCookieAuthKey = securecookie.GenerateRandomKey(64)
			proxyCookieCryptoKey = securecookie.GenerateRandomKey(32)
			settings.System.CookieAuthKey = cookieAuthKey
			settings.System.CookieCryptoKey = cookieCryptoKey
			settings.System.ProxyCookieAuthKey = proxyCookieAuthKey
			settings.System.ProxyCookieCryptoKey = proxyCookieCryptoKey

			fields := set.NewSet(
				"cookie_auth_key",
				"cookie_crypto_key",
				"proxy_cookie_auth_key",
				"proxy_cookie_crypto_key",
			)

			err = settings.Commit(db, settings.System, fields)
			if err != nil {
				return
			}
		}

		Store = sessions.NewCookieStore(
			cookieAuthKey, cookieCryptoKey)
		Store.Options.Secure = true
		Store.Options.HttpOnly = true

		return
	}
}
