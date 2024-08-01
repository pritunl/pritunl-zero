package dns

import (
	"github.com/dropbox/godropbox/errors"
)

type NotFoundError struct {
	errors.DropboxError
}

type ServiceError struct {
	errors.DropboxError
}

type UnknownError struct {
	errors.DropboxError
}
