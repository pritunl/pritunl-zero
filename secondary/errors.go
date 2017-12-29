package secondary

import (
	"github.com/dropbox/godropbox/errors"
)

type IncompleteError struct {
	errors.DropboxError
}
