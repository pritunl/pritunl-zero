package search

import (
	"bytes"
	"container/list"
	"context"
	"crypto/md5"
	"io"
	"sync"
	"time"

	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/sirupsen/logrus"
)

var (
	ctx          = context.Background()
	Default      *Client
	buffer       = list.New()
	lock         = sync.Mutex{}
	failedBuffer = list.New()
	failedLock   = sync.Mutex{}
	reconnect    = false
)

func hashConf(username, password string, addrs []string) []byte {
	hash := md5.New()

	io.WriteString(hash, username)
	io.WriteString(hash, password)

	for _, addr := range addrs {
		io.WriteString(hash, addr)
	}

	return hash.Sum(nil)
}

func updateClient(username, password string, addrs []string) (err error) {
	clnt, err := NewClient(username, password, addrs)
	if err != nil || clnt == nil {
		Default = nil
		return
	}

	err = clnt.UpdateIndexes()
	if err != nil {
		return
	}

	Default = clnt

	return
}

func watchSearch() {
	hash := hashConf("", "", []string{})

	for {
		username := settings.Elastic.Username
		password := settings.Elastic.Password
		addrs := settings.Elastic.Addresses
		newHash := hashConf(username, password, addrs)

		if bytes.Compare(hash, newHash) != 0 || reconnect {
			err := updateClient(username, password, addrs)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("search: Failed to update search indexes")
				time.Sleep(3 * time.Second)
				continue
			}

			hash = newHash
			reconnect = false
		}

		time.Sleep(1 * time.Second)
	}
}

func worker() {
	for {
		time.Sleep(1 * time.Second)

		group := []*Document{}
		groups := [][]*Document{}

		lock.Lock()
		for elem := buffer.Front(); elem != nil; elem = elem.Next() {
			if len(group) >= 100 {
				groups = append(groups, group)
				group = []*Document{}
			}
			doc := elem.Value.(*Document)
			group = append(group, doc)
		}
		buffer = list.New()
		lock.Unlock()

		if len(group) > 0 {
			groups = append(groups, group)
		}

		clnt := Default
		if clnt == nil || len(groups) == 0 {
			continue
		}

		for _, group := range groups {
			go func(docs []*Document) {
				e := clnt.BulkDocuments(docs, false)
				if e != nil {
					logrus.WithFields(logrus.Fields{
						"error": e,
					}).Error("search: Bulk insert error")
				}
			}(group)
		}
	}
}

func failedWorker() {
	for {
		time.Sleep(1 * time.Second)

		group := []*Document{}
		groups := [][]*Document{}

		lock.Lock()
		for elem := failedBuffer.Front(); elem != nil; elem = elem.Next() {
			if len(group) >= 10 {
				groups = append(groups, group)
				group = []*Document{}
			}
			doc := elem.Value.(*Document)
			group = append(group, doc)
		}
		buffer = list.New()
		lock.Unlock()

		if len(group) > 0 {
			groups = append(groups, group)
		}

		clnt := Default
		if clnt == nil || len(groups) == 0 {
			continue
		}

		for _, group := range groups {
			go func(docs []*Document) {
				e := clnt.BulkDocuments(docs, false)
				if e != nil {
					logrus.WithFields(logrus.Fields{
						"error": e,
					}).Error("search: Bulk insert retry error")
				}
			}(group)
		}
	}
}

func init() {
	module := requires.New("search")
	module.After("settings")

	module.Handler = func() (err error) {
		go watchSearch()
		go worker()
		go failedWorker()
		return
	}
}
