package search

import (
	"bytes"
	"container/list"
	"crypto/md5"
	"io"
	"sync"
	"time"

	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/sirupsen/logrus"
)

var (
	Default          *Client
	buffer           chan *Document
	failedBuffer     chan *Document
	reconnect        = false
	groups           *list.List
	groupsLock       = sync.Mutex{}
	failedGroups     *list.List
	failedGroupsLock = sync.Mutex{}
)

const (
	BufferSize       = 2048
	GroupLimit       = 100
	FailedGroupLimit = 10
	RetryCount       = 5
	ThreadLimit      = 10
)

func hashConf(username, password string, addrs []string) []byte {
	hash := md5.New()

	_, _ = io.WriteString(hash, username)
	_, _ = io.WriteString(hash, password)

	for _, addr := range addrs {
		_, _ = io.WriteString(hash, addr)
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

func workerBuffer() {
	for {
		entry := <-buffer

		groupsLock.Lock()

		group := groups.Back().Value.(*list.List)
		if group.Len() >= GroupLimit {
			group = list.New()
			groups.PushBack(group)
		}

		group.PushBack(entry)

		groupsLock.Unlock()
	}
}

func workerFailedBuffer() {
	for {
		entry := <-failedBuffer

		failedGroupsLock.Lock()

		failedGroup := failedGroups.Back().Value.(*list.List)
		if failedGroup.Len() >= FailedGroupLimit {
			failedGroup = list.New()
			failedGroups.PushBack(failedGroup)
		}

		failedGroup.PushBack(entry)

		failedGroupsLock.Unlock()
	}
}

func workerGroup() {
	for {
		groupsLock.Lock()
		curGrps := groups
		groups = list.New()
		groups.PushBack(list.New())
		groupsLock.Unlock()

		clnt := Default
		if clnt == nil || curGrps.Front().Value.(*list.List).Len() == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		for {
			waiters := sync.WaitGroup{}
			count := 0

			for elem := curGrps.Front(); elem != nil; elem = curGrps.Front() {
				group := curGrps.Remove(elem).(*list.List)
				if group.Len() == 0 {
					continue
				}

				count += 1
				waiters.Add(1)
				go func(group *list.List) {
					err := clnt.BulkDocuments(group, true)
					if err != nil {
						if logLimit() {
							logrus.WithFields(logrus.Fields{
								"error": err,
							}).Error("search: Bulk insert error")
						}
					}

					waiters.Done()
				}(group)

				if count >= ThreadLimit {
					break
				}
			}

			waiters.Wait()

			if curGrps.Front() == nil {
				break
			}
		}
	}
}

func workerFailedGroup() {
	for {
		failedGroupsLock.Lock()
		curGrps := failedGroups
		failedGroups = list.New()
		failedGroups.PushBack(list.New())
		failedGroupsLock.Unlock()

		clnt := Default
		if clnt == nil || curGrps.Front().Value.(*list.List).Len() == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		for {
			waiters := sync.WaitGroup{}
			count := 0

			for elem := curGrps.Front(); elem != nil; elem = curGrps.Front() {
				group := curGrps.Remove(elem).(*list.List)
				if group.Len() == 0 {
					continue
				}

				count += 1
				waiters.Add(1)
				go func(group *list.List) {
					err := clnt.BulkDocuments(group, false)
					if err != nil {
						if logLimit() {
							logrus.WithFields(logrus.Fields{
								"error": err,
							}).Error("search: Bulk insert retry error")
						}
					}

					waiters.Done()
				}(group)

				if count >= ThreadLimit {
					break
				}
			}

			waiters.Wait()

			if curGrps.Front() == nil {
				break
			}
		}
	}
}

func init() {
	buffer = make(chan *Document, BufferSize+500)
	failedBuffer = make(chan *Document, BufferSize+500)
	groups = list.New()
	groups.PushBack(list.New())
	failedGroups = list.New()
	failedGroups.PushBack(list.New())

	module := requires.New("search")
	module.After("settings")

	module.Handler = func() (err error) {
		go watchSearch()
		go workerBuffer()
		go workerFailedBuffer()
		go workerGroup()
		go workerFailedGroup()

		return
	}
}
