package errortypes

import (
	"github.com/dropbox/godropbox/errors"
)

type UnknownError struct {
	errors.DropboxError
}

type NotFoundError struct {
	errors.DropboxError
}

type ReadError struct {
	errors.DropboxError
}

type WriteError struct {
	errors.DropboxError
}

type ParseError struct {
	errors.DropboxError
}

type AuthenticationError struct {
	errors.DropboxError
}

type ApiError struct {
	errors.DropboxError
}

type DatabaseError struct {
	errors.DropboxError
}

type RequestError struct {
	errors.DropboxError
}

type ErrorData struct {
	Error   string `json:"error"`
	Message string `json:"error_msg"`
}
