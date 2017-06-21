package requires

import (
	"github.com/dropbox/godropbox/errors"
)

type InitError struct {
	errors.DropboxError
}
