package database

import (
	"fmt"
	"strings"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/mongo"
	"github.com/pritunl/mongo-go-driver/mongo/options"
)

func FindProject(fields ...string) *options.FindOptions {
	prcj := bson.M{}

	for _, field := range fields {
		prcj[field] = 1
	}

	return &options.FindOptions{
		Projection: prcj,
	}
}

func FindOneProject(fields ...string) *options.FindOneOptions {
	prcj := bson.M{}

	for _, field := range fields {
		prcj[field] = 1
	}

	return &options.FindOneOptions{
		Projection: prcj,
	}
}

func GetErrorCodes(err error) (errCodes []int) {
	switch err := err.(type) {
	case mongo.CommandError:
		errCodes = []int{int(err.Code)}
		if strings.Contains(err.Name, "Conflict") {
			errCodes = append(errCodes, 85)
		}
		break
	case mongo.WriteError:
		errCodes = []int{err.Code}
		break
	case mongo.BulkWriteError:
		errCodes = []int{err.Code}
		break
	case mongo.WriteConcernError:
		errCodes = []int{err.Code}
		break
	case mongo.WriteException:
		errCodes = []int{}
		if err.WriteConcernError != nil {
			errCodes = append(errCodes, err.WriteConcernError.Code)
		}
		if err.WriteErrors != nil {
			for _, e := range err.WriteErrors {
				errCodes = append(errCodes, e.Code)
			}
		}
		break
	case mongo.WriteErrors:
		errCodes = []int{}
		for _, e := range err {
			eCodes := GetErrorCodes(e)
			errCodes = append(errCodes, eCodes...)
		}
		break
	case *mongo.WriteError:
		errCodes = []int{err.Code}
		break
	case *mongo.BulkWriteError:
		errCodes = []int{err.Code}
		break
	case *mongo.WriteConcernError:
		errCodes = []int{err.Code}
		break
	case *mongo.WriteException:
		errCodes = []int{}
		if err.WriteConcernError != nil {
			errCodes = append(errCodes, err.WriteConcernError.Code)
		}
		if err.WriteErrors != nil {
			for _, e := range err.WriteErrors {
				errCodes = append(errCodes, e.Code)
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

	errCodes := GetErrorCodes(err)
	for _, errCode := range errCodes {
		switch errCode {
		case 66:
			newErr = &ImmutableKeyError{
				errors.New("database: Immutable key"),
			}
			return
		case 85:
			newErr = &IndexConflict{
				errors.New("database: Index conflict"),
			}
			return
		case 11000, 11001, 12582, 16460:
			newErr = &DuplicateKeyError{
				errors.New("database: Duplicate key"),
			}
			return
		}
	}

	newErr = &UnknownError{
		errors.Wrap(err, fmt.Sprintf(
			"database: Unknown error %v", errCodes)),
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
