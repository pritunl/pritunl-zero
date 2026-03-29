package database

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/sirupsen/logrus"
)

var (
	indexes     = map[string]set.Set{}
	indexesLock = sync.Mutex{}
)

type bsonIndex struct {
	Name string `bson:"name"`
}

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

	collName := i.Collection.Name()
	indexesLock.Lock()
	collIndexes, ok := indexes[collName]
	if !ok {
		collIndexes = set.NewSet()
		indexes[collName] = collIndexes
	}
	collIndexes.Add(name)
	indexesLock.Unlock()

	return
}

func CleanIndexes(db *Database) (err error) {
	indexesLock.Lock()
	curIndexes := indexes
	indexesLock.Unlock()

	for collName, collIndexes := range curIndexes {
		coll := db.GetCollection(collName)

		cursor, e := coll.Indexes().List(db)
		if e != nil {
			err = e
			return
		}

		for cursor.Next(db) {
			index := &bsonIndex{}

			err = cursor.Decode(index)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"collection": collName,
					"error":      err,
				}).Error("database: Failed to decode index")
				err = nil
				continue
			}

			if index.Name == "_id" || index.Name == "_id_" {
				continue
			}

			if collIndexes.Contains(index.Name) {
				continue
			}

			logrus.WithFields(logrus.Fields{
				"collection": collName,
				"index":      index.Name,
			}).Info("database: Dropping unused index")

			err = coll.Indexes().DropOne(
				db,
				index.Name,
			)
			if err != nil {
				cursor.Close(db)
				return
			}
		}

		err = cursor.Err()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"collection": collName,
				"error":      err,
			}).Error("database: Cursor error listing indexes")
		}

		cursor.Close(db)
	}

	return
}
