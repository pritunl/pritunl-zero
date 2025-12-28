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
	groups           *IndexList
	groupsLock       = sync.Mutex{}
	failedGroups     *IndexList
	failedGroupsLock = sync.Mutex{}
)

const (
	BufferLenMax = 2048
	RetryCount   = 5
	ThreadLimit  = 10
)

type IndexList struct {
	*list.List
	Size int
}

func NewIndexList() *IndexList {
	return &IndexList{
		List: list.New(),
		Size: 0,
	}
}

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
		entrySize := len(entry.Data)

		groupsLock.Lock()

		group := groups.Back().Value.(*IndexList)
		if group.Len()+1 > settings.Elastic.GroupLength ||
			group.Size+entrySize > settings.Elastic.GroupSize {

			group = NewIndexList()
			groups.PushBack(group)
		}

		group.PushBack(entry)
		group.Size += entrySize

		groupsLock.Unlock()
	}
}

func workerFailedBuffer() {
	for {
		entry := <-failedBuffer
		entrySize := len(entry.Data)

		failedGroupsLock.Lock()

		failedGroup := failedGroups.Back().Value.(*IndexList)
		if failedGroup.Len()+1 > settings.Elastic.GroupLength ||
			failedGroup.Size+entrySize > settings.Elastic.GroupSize {

			failedGroup = NewIndexList()
			failedGroups.PushBack(failedGroup)
		}

		failedGroup.PushBack(entry)
		failedGroup.Size += entrySize

		failedGroupsLock.Unlock()
	}
}

func workerGroup() {
	for {
		groupsLock.Lock()
		curGrps := groups
		groups = NewIndexList()
		groups.PushBack(NewIndexList())
		groupsLock.Unlock()

		clnt := Default
		if clnt == nil || curGrps.Front().Value.(*IndexList).Len() == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		for {
			waiters := sync.WaitGroup{}
			count := 0

			for elem := curGrps.Front(); elem != nil; elem = curGrps.Front() {
				group := curGrps.Remove(elem).(*IndexList)
				if group.Len() == 0 {
					continue
				}

				count += 1
				waiters.Add(1)
				go func(group *IndexList) {
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
		failedGroups = NewIndexList()
		failedGroups.PushBack(NewIndexList())
		failedGroupsLock.Unlock()

		clnt := Default
		if clnt == nil || curGrps.Front().Value.(*IndexList).Len() == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		for {
			waiters := sync.WaitGroup{}
			count := 0

			for elem := curGrps.Front(); elem != nil; elem = curGrps.Front() {
				group := curGrps.Remove(elem).(*IndexList)
				if group.Len() == 0 {
					continue
				}

				count += 1
				waiters.Add(1)
				go func(group *IndexList) {
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
	buffer = make(chan *Document, BufferLenMax)
	failedBuffer = make(chan *Document, BufferLenMax)
	groups = NewIndexList()
	groups.PushBack(NewIndexList())
	failedGroups = NewIndexList()
	failedGroups.PushBack(NewIndexList())

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
