package cookie

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"net/http"
	"strconv"
	"strings"
)

func Get(w http.ResponseWriter, r *http.Request) (cook *Cookie, err error) {
	store, err := Store.New(r, "pritunl-zero-console")
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err.(securecookie.MultiError)[0],
				"cookie: Unknown cookie error"),
		}
		return
	}

	cook = &Cookie{
		store: store,
		w:     w,
		r:     r,
	}

	return
}

func New(w http.ResponseWriter, r *http.Request) (cook *Cookie) {
	store, _ := Store.New(r, "pritunl-zero-console")

	cook = &Cookie{
		store: store,
		w:     w,
		r:     r,
	}

	return
}

func getCookieTopDomain(r *http.Request) string {
	host := utils.StripPort(r.Host)
	if strings.Count(host, ".") >= 2 {
		i := strings.LastIndex(host, ".")
		tld := host[i+1:]
		if _, err := strconv.Atoi(tld); err != nil {
			host = "." + strings.SplitN(host, ".", 2)[1]
			return host
		}
	}

	return ""
}

func newProxyStore(srvc *service.Service,
	r *http.Request) *sessions.CookieStore {

	cookieStore := sessions.NewCookieStore(
		settings.System.ProxyCookieAuthKey,
		settings.System.ProxyCookieCryptoKey,
	)
	cookieStore.Options.Secure = true

	if srvc.ShareSession {
		cookieStore.Options.Domain = getCookieTopDomain(r)
	}

	return cookieStore
}

func GetProxy(srvc *service.Service, w http.ResponseWriter, r *http.Request) (
	cook *Cookie, err error) {

	cookStore := newProxyStore(srvc, r)

	store, err := cookStore.New(r, "pritunl-zero")
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err.(securecookie.MultiError)[0],
				"cookie: Unknown cookie error"),
		}
		return
	}

	cook = &Cookie{
		store: store,
		w:     w,
		r:     r,
	}

	return
}

func NewProxy(srvc *service.Service, w http.ResponseWriter, r *http.Request) (
	cook *Cookie) {

	cookStore := newProxyStore(srvc, r)

	store, _ := cookStore.New(r, "pritunl-zero")

	cook = &Cookie{
		store: store,
		w:     w,
		r:     r,
	}

	return
}
