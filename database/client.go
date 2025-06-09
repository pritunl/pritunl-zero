package database

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pritunl/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

var (
	globalClient     atomic.Value
	globalClientLock sync.Mutex
	DefaultDatabase  string
)

func getClient() *mongo.Client {
	val := globalClient.Load()
	if val == nil {
		return nil
	}
	return val.(*mongo.Client)
}

func setClient(client *mongo.Client) {
	globalClientLock.Lock()
	curClientInf := globalClient.Load()
	if curClientInf != nil {
		curClient := curClientInf.(*mongo.Client)
		ctx, cancel := context.WithTimeout(
			context.Background(),
			30*time.Second,
		)
		defer cancel()

		err := curClient.Disconnect(ctx)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("database: Disconnect error")
		}
	}
	globalClient.Store(client)
	globalClientLock.Unlock()
}
