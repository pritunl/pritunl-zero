package ssh

import (
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/requires"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

var (
	registry     = map[string]map[bson.ObjectId]func(){}
	registryLock = sync.Mutex{}
)

func Register(token string, callback func()) bson.ObjectId {
	listernerId := bson.NewObjectId()

	registryLock.Lock()
	defer registryLock.Unlock()

	callbacks, ok := registry[token]
	if !ok {
		callbacks = map[bson.ObjectId]func(){}
	}
	callbacks[listernerId] = callback
	registry[token] = callbacks

	return listernerId
}

func Unregister(token string, listenerId bson.ObjectId) {
	registryLock.Lock()
	defer registryLock.Unlock()

	callbacks, ok := registry[token]
	if ok {
		delete(callbacks, listenerId)
		if len(callbacks) == 0 {
			delete(registry, token)
		} else {
			registry[token] = callbacks
		}
	}
}

func callback(evt *event.Event) {
	token := evt.Data.(string)

	registryLock.Lock()
	defer registryLock.Unlock()

	callbacks, ok := registry[token]
	if ok {
		for _, callback := range callbacks {
			go callback()
		}
	}
}

func init() {
	module := requires.New("sshcert")
	module.After("settings")
	module.Before("event")

	module.Handler = func() (err error) {
		event.Register("ssh_challenge", callback)
		return
	}
}
