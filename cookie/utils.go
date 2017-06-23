package cookie

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/pritunl/pritunl-zero/errortypes"
)

func Get(c *gin.Context) (cook *Cookie, err error) {
	store, err := Store.New(c.Request, "pritunl-zero-console")
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err.(securecookie.MultiError)[0],
				"cookie: Unknown cookie error"),
		}
		return
	}

	cook = &Cookie{
		store: store,
		con:   c,
	}

	return
}

func New(c *gin.Context) (cook *Cookie) {
	store, _ := Store.New(c.Request, "pritunl-zero-console")

	cook = &Cookie{
		store: store,
		con:   c,
	}

	return
}

func GetProxy(c *gin.Context) (cook *Cookie, err error) {
	store, err := ProxyStore.New(c.Request, "pritunl-zero")
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err.(securecookie.MultiError)[0],
				"cookie: Unknown cookie error"),
		}
		return
	}

	cook = &Cookie{
		store: store,
		con:   c,
	}

	return
}

func NewProxy(c *gin.Context) (cook *Cookie) {
	store, _ := ProxyStore.New(c.Request, "pritunl-zero")

	cook = &Cookie{
		store: store,
		con:   c,
	}

	return
}
