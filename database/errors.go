package database

import (
	"github.com/dropbox/godropbox/errors"
)

type ConnectionError struct {
	errors.DropboxError
}

type IndexError struct {
	errors.DropboxError
}

type NotFoundError struct {
	errors.DropboxError
}

type DuplicateKeyError struct {
	errors.DropboxError
}

type UnknownError struct {
	errors.DropboxError
}

type CertificateError struct {
	errors.DropboxError
}
