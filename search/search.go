package search

import (
	"bytes"
	"container/list"
	"context"
	"crypto/md5"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/settings"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v6"
	"io"
	"sync"
	"time"
)

var (
	ctx          = context.Background()
	client       *elastic.Client
	buffer       = list.New()
	failedBuffer = list.New()
	lock         = sync.Mutex{}
	failedLock   = sync.Mutex{}
)

type mapping struct {
	Field string
	Type  string
	Store bool
	Index bool
}

func Index(index, typ string, data interface{}) {
	clnt := client
	if clnt == nil {
		return
	}

	id := bson.NewObjectId().Hex()

	request := elastic.NewBulkIndexRequest().Index(index).Type(typ).
		Id(id).Doc(data)

	lock.Lock()
	buffer.PushBack(request)
	lock.Unlock()

	return
}

func putIndex(clnt *elastic.Client, index string, typ string,
	mappings []mapping) (err error) {

	exists, err := clnt.IndexExists(index).Do(ctx)
	if err != nil {
		err = &errortypes.DatabaseError{
			errors.Wrap(err, "search: Failed to check elastic index"),
		}
		return
	}

	if exists {
		return
	}

	properties := map[string]interface{}{}

	for _, mapping := range mappings {
		if mapping.Type == "object" {
			properties[mapping.Field] = struct {
				Enabled bool `json:"enabled"`
			}{
				Enabled: false,
			}
		} else {
			properties[mapping.Field] = struct {
				Type  string `json:"type"`
				Store bool   `json:"store"`
				Index bool   `json:"index"`
			}{
				Type:  mapping.Type,
				Store: mapping.Store,
				Index: mapping.Index,
			}
		}
	}

	data := struct {
		Mappings map[string]interface{} `json:"mappings"`
	}{
		Mappings: map[string]interface{}{},
	}

	data.Mappings[typ] = struct {
		Properties map[string]interface{} `json:"properties"`
	}{
		Properties: properties,
	}

	_, err = clnt.CreateIndex(index).BodyJson(data).Do(ctx)
	if err != nil {
		err = &errortypes.DatabaseError{
			errors.Wrap(err, "search: Failed to create elastic index"),
		}
		return
	}

	return
}

func newClient(addrs []string) (clnt *elastic.Client, err error) {
	if len(addrs) == 0 {
		return
	}

	clnt, err = elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(addrs...),
	)
	if err != nil {
		err = &errortypes.DatabaseError{
			errors.Wrap(err, "search: Failed to create elastic client"),
		}
		return
	}

	return
}

func hashAddresses(addrs []string) []byte {
	hash := md5.New()

	for _, addr := range addrs {
		io.WriteString(hash, addr)
	}

	return hash.Sum(nil)
}

func update(addrs []string) (err error) {
	clnt, err := newClient(addrs)
	if err != nil {
		client = nil
		return
	}

	if clnt == nil {
		client = nil
		return
	}

	mappings := []mapping{}

	mappings = append(mappings, mapping{
		Field: "user",
		Type:  "keyword",
		Store: false,
		Index: true,
	})

	mappings = append(mappings, mapping{
		Field: "username",
		Type:  "keyword",
		Store: false,
		Index: true,
	})

	mappings = append(mappings, mapping{
		Field: "session",
		Type:  "keyword",
		Store: false,
		Index: true,
	})

	mappings = append(mappings, mapping{
		Field: "address",
		Type:  "ip",
		Store: false,
		Index: true,
	})

	mappings = append(mappings, mapping{
		Field: "timestamp",
		Type:  "date",
		Store: false,
		Index: true,
	})

	mappings = append(mappings, mapping{
		Field: "scheme",
		Type:  "keyword",
		Store: false,
		Index: true,
	})

	mappings = append(mappings, mapping{
		Field: "host",
		Type:  "keyword",
		Store: false,
		Index: true,
	})

	mappings = append(mappings, mapping{
		Field: "path",
		Type:  "keyword",
		Store: false,
		Index: true,
	})

	mappings = append(mappings, mapping{
		Field: "query",
		Type:  "object",
		Index: true,
	})

	mappings = append(mappings, mapping{
		Field: "header",
		Type:  "object",
		Index: true,
	})

	mappings = append(mappings, mapping{
		Field: "body",
		Type:  "text",
		Store: false,
		Index: false,
	})

	err = putIndex(clnt, "zero-requests", "request", mappings)
	if err != nil {
		client = nil
		return
	}

	client = clnt

	return
}

func watchSearch() {
	hash := hashAddresses([]string{})

	for {
		addrs := settings.Elastic.Addresses
		newHash := hashAddresses(addrs)

		if bytes.Compare(hash, newHash) != 0 {
			err := update(addrs)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("search: Failed to update search indexes")
				time.Sleep(3 * time.Second)
				continue
			}

			hash = newHash
		}

		time.Sleep(1 * time.Second)
	}
}

func sendRequests(clnt *elastic.Client,
	requests []*elastic.BulkIndexRequest, log bool) {
	var err error

	for i := 0; i < 10; i++ {
		bulk := clnt.Bulk()

		for _, request := range requests {
			bulk.Add(request)
		}

		_, err = bulk.Do(ctx)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		err = nil
		break
	}

	if err != nil {
		failedLock.Lock()
		for _, request := range requests {
			failedBuffer.PushBack(request)
		}
		failedLock.Unlock()

		if log {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("search: Bulk insert failed, moving to failed buffer")
		}
	}
}

func worker() {
	for {
		time.Sleep(1 * time.Second)

		group := []*elastic.BulkIndexRequest{}
		requests := [][]*elastic.BulkIndexRequest{}

		lock.Lock()
		for elem := buffer.Front(); elem != nil; elem = elem.Next() {
			if len(group) >= 100 {
				requests = append(requests, group)
				group = []*elastic.BulkIndexRequest{}
			}
			request := elem.Value.(*elastic.BulkIndexRequest)
			group = append(group, request)
		}
		buffer = list.New()
		lock.Unlock()

		if len(group) > 0 {
			requests = append(requests, group)
		}

		clnt := client
		if client == nil {
			continue
		}

		if len(requests) == 0 {
			continue
		}

		for _, group := range requests {
			go sendRequests(clnt, group, true)
		}
	}
}

func failedWorker() {
	for {
		time.Sleep(1 * time.Second)

		group := []*elastic.BulkIndexRequest{}
		requests := [][]*elastic.BulkIndexRequest{}

		lock.Lock()
		for elem := failedBuffer.Front(); elem != nil; elem = elem.Next() {
			if len(group) >= 10 {
				requests = append(requests, group)
				group = []*elastic.BulkIndexRequest{}
			}
			request := elem.Value.(*elastic.BulkIndexRequest)
			group = append(group, request)
		}
		buffer = list.New()
		lock.Unlock()

		if len(group) > 0 {
			requests = append(requests, group)
		}

		clnt := client
		if client == nil {
			continue
		}

		if len(requests) == 0 {
			continue
		}

		for _, group := range requests {
			go sendRequests(clnt, group, false)
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
