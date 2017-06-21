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
		case *mgo.QueryError:
		errCode = err.Code
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
	default:
		newErr = &UnknownError{
			errors.Wrap(err, fmt.Sprintf(
				"database: Unknown error %d", errCode)),
		}
	}

	return
}
