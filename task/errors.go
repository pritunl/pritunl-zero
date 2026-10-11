package task

import (
	"github.com/dropbox/godropbox/errors"
)

type ReservationLostError struct {
	errors.DropboxError
}

type DurationElapsedError struct {
	errors.DropboxError
}

type ShutdownError struct {
	errors.DropboxError
}
