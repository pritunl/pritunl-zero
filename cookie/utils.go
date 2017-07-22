package cookie

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/gorilla/securecookie"
	"github.com/pritunl/pritunl-zero/errortypes"
	"net/http"
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

func GetProxy(w http.ResponseWriter, r *http.Request) (
	cook *Cookie, err error) {

	store, err := ProxyStore.New(r, "pritunl-zero")
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

func NewProxy(w http.ResponseWriter, r *http.Request) (cook *Cookie) {
	store, _ := ProxyStore.New(r, "pritunl-zero")

	cook = &Cookie{
		store: store,
		w:     w,
		r:     r,
	}

	return
}
