package database

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
)

var (
	indexes     = []string{}
	indexesLock = sync.Mutex{}
)

type Index struct {
	Collection *Collection
	Keys       *bson.D
	Unique     bool
	Partial    interface{}
	Expire     time.Duration
}

func GenerateIndexName(doc bson.D) (indexName string, err error) {
	name := bytes.NewBufferString("")
	first := true

	for _, elem := range doc {
		if !first {
			_, err = name.WriteRune('_')
			if err != nil {
				err = &UnknownError{
					errors.Wrap(err, "database: Write rune error"),
				}
				return
			}
		}

		_, err = name.WriteString(elem.Key)
		if err != nil {
			err = &UnknownError{
				errors.Wrap(err, "database: Write string error"),
			}
			return
		}

		_, err = name.WriteRune('_')
		if err != nil {
			err = &UnknownError{
				errors.Wrap(err, "database: Write rune error"),
			}
			return
		}

		value := ""
		switch val := elem.Value.(type) {
		case int, int32, int64:
			value = fmt.Sprintf("%d", val)
		case string:
			value = val
		default:
			err = &UnknownError{
				errors.New("database: Invalid index value"),
			}
			return
		}

		_, err = name.WriteString(value)
		if err != nil {
			err = &UnknownError{
				errors.Wrap(err, "database: Write string error"),
			}
			return
		}

		first = false
	}

	indexName = name.String()
	return
}

func (i *Index) Create() (err error) {
	opts := options.Index()

	if i.Unique {
		opts.SetUnique(true)
	}

	if i.Partial != nil {
		opts.SetPartialFilterExpression(i.Partial)
	}

	if i.Expire != 0 {
		opts.SetExpireAfterSeconds(int32(i.Expire.Seconds()))
	}

	indexName, err := GenerateIndexName(*i.Keys)
	if err != nil {
		return
	}
	opts.SetName(indexName)

	name, err := i.Collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    i.Keys,
			Options: opts,
		},
	)
	if err != nil {
		err = ParseError(err)
		if _, ok := err.(*IndexConflict); ok {
			err = nil

			err = i.Collection.Indexes().DropOne(
				context.Background(),
				indexName,
			)
			if err != nil {
				return
			}

			name, err = i.Collection.Indexes().CreateOne(
				context.Background(),
				mongo.IndexModel{
					Keys:    i.Keys,
					Options: opts,
				},
			)
			if err != nil {
				err = ParseError(err)
				return
			}
		} else {
			return
		}
	}

	indexesLock.Lock()
	indexes = append(indexes, name)
	indexesLock.Unlock()

	return
}
