package authority

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"io"
	"sync"
	"time"

	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/settings"
)

var (
	clientCertCache     = map[primitive.ObjectID]*clientCertCacheData{}
	clientCertCacheLock sync.Mutex
)

type clientCertCacheData struct {
	HashKey     []byte
	Timestamp   time.Time
	Certificate *tls.Certificate
}

func clientCertAuthrHashKey(authr *Authority) []byte {
	hash := md5.New()

	_, _ = io.WriteString(hash, authr.PrivateKey)
	_, _ = io.WriteString(hash, authr.RootCertificate)

	return hash.Sum(nil)
}

func clientCertCacheGet(authr *Authority) (cert *tls.Certificate) {
	clientCertCacheLock.Lock()
	curCache := clientCertCache[authr.Id]
	clientCertCacheLock.Unlock()

	if curCache == nil {
		return
	}

	hashKey := clientCertAuthrHashKey(authr)
	if bytes.Compare(hashKey, curCache.HashKey) != 0 {
		return
	}

	if time.Since(curCache.Timestamp) > time.Duration(
		settings.System.ClientCertCacheTtl)*time.Second {

		return
	}

	cert = curCache.Certificate

	return
}

func clientCertCacheSet(authr *Authority, cert *tls.Certificate) {
	hashKey := clientCertAuthrHashKey(authr)

	cache := &clientCertCacheData{
		HashKey:     hashKey,
		Timestamp:   time.Now(),
		Certificate: cert,
	}

	clientCertCacheLock.Lock()
	clientCertCache[authr.Id] = cache
	clientCertCacheLock.Unlock()

	return
}

func clientCertCacheWatch() {
	go func() {
		for {
			time.Sleep(300 * time.Second)

			clientCertCacheLock.Lock()
			for key, authr := range clientCertCache {
				if time.Since(authr.Timestamp) > time.Duration(
					settings.System.ClientCertCacheTtl)*time.Second {

					delete(clientCertCache, key)
				}
			}
			clientCertCacheLock.Unlock()
		}
	}()
}
