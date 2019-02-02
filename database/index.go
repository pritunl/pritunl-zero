package database

import (
	"context"
	"sync"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/mongo"
	"github.com/pritunl/mongo-go-driver/mongo/options"
)

var (
	indexes     = []string{}
	indexesLock = sync.Mutex{}
)

type Index struct {
	Collection *Collection
	Keys       interface{}
	Unique     bool
	Expire     time.Duration
}

func (i *Index) Create() (err error) {
	opts := &options.IndexOptions{}
	opts.SetBackground(true)

	if i.Unique {
		opts.SetUnique(true)
	}

	if i.Expire != 0 {
		opts.SetExpireAfterSeconds(int32(i.Expire.Seconds()))
	}

	name, err := i.Collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    i.Keys,
			Options: opts,
		},
	)
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
		return
	}

	indexesLock.Lock()
	indexes = append(indexes, name)
	indexesLock.Unlock()

	return
}
