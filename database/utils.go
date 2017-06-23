package database

import (
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/mgo.v2"
)

// Get mongodb error code from error
func GetErrorCode(err error) (errCode int) {
	switch err := err.(type) {
	case *mgo.LastError:
		errCode = err.Code
		break
	case *mgo.QueryError:
		errCode = err.Code
		break
	}

	return
}

// Parse database error data and return error type
func ParseError(err error) (newErr error) {
	if err == mgo.ErrNotFound {
		newErr = &NotFoundError{
			errors.New("database: Not found"),
		}
		return
	}

	errCode := GetErrorCode(err)

	switch errCode {
	case 11000, 11001, 12582, 16460:
		newErr = &DuplicateKeyError{
			errors.New("database: Duplicate key"),
		}
		break
	default:
		newErr = &UnknownError{
			errors.Wrap(err, fmt.Sprintf(
				"database: Unknown error %d", errCode)),
		}
	}

	return
}

// Ignore not found error
func IgnoreNotFoundError(err error) (newErr error) {
	if err != nil {
		switch err.(type) {
		case *NotFoundError:
			newErr = nil
			break
		default:
			newErr = err
		}
	}

	return
}
