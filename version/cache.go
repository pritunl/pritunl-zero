package version

import (
	"sync"
	"time"
)

var (
	cacheStore = map[string]*cache{}
	cacheLock  = sync.Mutex{}
)

const (
	cacheTtl = 5 * time.Minute
)

type cache struct {
	Version   int
	Timestamp time.Time
}

func cacheCheck(module string, ver int) (supported bool) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	cach, ok := cacheStore[module]
	if !ok {
		return true
	}

	if time.Since(cach.Timestamp) > cacheTtl {
		delete(cacheStore, module)
		return true
	}

	return ver >= cach.Version
}

func cacheSet(module string, ver int) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	cacheStore[module] = &cache{
		Version:   ver,
		Timestamp: time.Now(),
	}

	existing, ok := cacheStore[module]
	if !ok || ver > existing.Version {
		cacheStore[module] = &cache{
			Version:   ver,
			Timestamp: time.Now(),
		}
	}
}
