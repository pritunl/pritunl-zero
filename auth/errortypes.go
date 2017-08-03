package auth

import (
	"github.com/dropbox/godropbox/errors"
)

type InvalidState struct {
	errors.DropboxError
}
