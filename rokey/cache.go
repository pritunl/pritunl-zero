package rokey

import (
	"fmt"
	"sync"
	"time"

	"github.com/pritunl/mongo-go-driver/v2/bson"
)

var (
	cache         = map[bson.ObjectID]*Rokey{}
	cacheLock     = sync.RWMutex{}
	cacheTime     = map[string]*Rokey{}
	cacheTimeLock = sync.RWMutex{}
)

func GetCache(typ string, timeblock time.Time) *Rokey {
	cacheTimeLock.RLock()
	rkey := cacheTime[fmt.Sprintf("%s-%d", typ, timeblock.Unix())]
	cacheTimeLock.RUnlock()
	if rkey != nil && rkey.Type == typ {
		return rkey
	}
	return nil
}

func GetCacheId(typ string, rkeyId bson.ObjectID) *Rokey {
	cacheLock.RLock()
	rkey := cache[rkeyId]
	cacheLock.RUnlock()
	if rkey != nil && rkey.Type == typ {
		return rkey
	}
	return nil
}

func PutCache(rkey *Rokey) {
	cacheLock.Lock()
	cache[rkey.Id] = rkey
	cacheLock.Unlock()
	cacheTimeLock.Lock()
	cacheTime[fmt.Sprintf("%s-%d", rkey.Type, rkey.Timeblock.Unix())] = rkey
	cacheTimeLock.Unlock()
}

func CleanCache() {
	cacheLock.Lock()
	for key, rkey := range cache {
		if time.Since(rkey.Timestamp) >= 721*time.Hour {
			delete(cache, key)
		}
	}
	cacheLock.Unlock()

	cacheTimeLock.Lock()
	for key, rkey := range cacheTime {
		if time.Since(rkey.Timestamp) >= 721*time.Hour {
			delete(cacheTime, key)
		}
	}
	cacheTimeLock.Unlock()
}

func init() {
	go func() {
		time.Sleep(1 * time.Hour)
		CleanCache()
	}()
}
