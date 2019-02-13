package database

import (
	"fmt"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/mongo"
)

func GetErrorCode(err error) (errCode int) {
	switch err := err.(type) {
	case mongo.WriteError:
		errCode = err.Code
		break
	case mongo.BulkWriteError:
		errCode = err.Code
		break
	case mongo.WriteConcernError:
		errCode = err.Code
		break
	case mongo.WriteException:
		if err.WriteConcernError != nil {
			errCode = err.WriteConcernError.Code
		} else if err.WriteErrors != nil {
			for _, e := range err.WriteErrors {
				if e.Code != 0 {
					errCode = e.Code
					break
				}
			}
		}
		break
	case *mongo.WriteError:
		errCode = err.Code
		break
	case *mongo.BulkWriteError:
		errCode = err.Code
		break
	case *mongo.WriteConcernError:
		errCode = err.Code
		break
	case *mongo.WriteException:
		if err.WriteConcernError != nil {
			errCode = err.WriteConcernError.Code
		} else if err.WriteErrors != nil {
			for _, e := range err.WriteErrors {
				if e.Code != 0 {
					errCode = e.Code
					break
				}
			}
		}
		break
	}

	return
}

func ParseError(err error) (newErr error) {
	if err == mongo.ErrNoDocuments {
		newErr = &NotFoundError{
			errors.New("database: Not found"),
		}
		return
	}

	if errs, ok := err.(mongo.WriteErrors); ok {
		errCode := 0
		for _, e := range errs {
			errCode = GetErrorCode(&e)
			if errCode == 11000 || errCode == 11001 || errCode == 12582 ||
				errCode == 16460 {

				newErr = &DuplicateKeyError{
					errors.New("database: Duplicate key"),
				}
				return
			}
		}
		newErr = &UnknownError{
			errors.Wrap(err, fmt.Sprintf(
				"database: Unknown error %d", errCode)),
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
