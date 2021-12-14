package search

import (
	"bytes"
	"container/list"
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/sirupsen/logrus"
)

var (
	ctx          = context.Background()
	client       *opensearch.Client
	buffer       = list.New()
	failedBuffer = list.New()
	lock         = sync.Mutex{}
	failedLock   = sync.Mutex{}
	reconnect    = false
)

type mapping struct {
	Field string
	Type  string
	Store bool
	Index bool
}

type document struct {
	Index string
	Id    string
	Data  interface{}
}

type bulkIndexReqData struct {
	Index string `json:"_index"`
	Id    string `json:"_id"`
}

type bulkIndexReq struct {
	Index *bulkIndexReqData `json:"index"`
}

func Index(index string, data interface{}) {
	clnt := client
	if clnt == nil {
		return
	}

	doc := &document{
		Index: index,
		Id:    primitive.NewObjectID().Hex(),
		Data:  data,
	}

	lock.Lock()
	buffer.PushBack(doc)
	lock.Unlock()

	return
}

func putIndex(clnt *opensearch.Client, index string,
	mappings []mapping) (err error) {

	exists, err := clnt.Indices.Exists([]string{index})
	if err != nil {
		err = &errortypes.DatabaseError{
			errors.Wrap(err, "search: Failed to check elastic index"),
		}
		return
	}

	if exists.StatusCode == 200 {
		return
	}

	properties := map[string]interface{}{}

	for _, mapping := range mappings {
		if mapping.Type == "object" {
			properties[mapping.Field] = struct {
				Enabled bool `json:"enabled"`
			}{
				Enabled: mapping.Index,
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

	data.Mappings["properties"] = properties

	body, err := json.Marshal(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "search: Failed to marshal data"),
		}
		return
	}

	resp, err := clnt.Indices.Create(
		index,
		func(r *opensearchapi.IndicesCreateRequest) {
			r.Body = bytes.NewBuffer(body)
		},
	)
	if err != nil {
		err = &errortypes.DatabaseError{
			errors.Wrap(err, "search: Failed to create elastic index"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody := ""
		respData, _ := ioutil.ReadAll(resp.Body)
		if respData != nil {
			respBody = string(respData)
		}

		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"response":    respBody,
			"error":       err,
		}).Error("search: Failed to create elastic index")

		err = &errortypes.DatabaseError{
			errors.New("search: Failed to create elastic index"),
		}
		return
	}

	return
}

func newClient(username, password string, addrs []string) (
	clnt *opensearch.Client, err error) {

	if len(addrs) == 0 {
		return
	}

	cfg := opensearch.Config{
		Addresses: addrs,
		Username:  username,
		Password:  password,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS13,
			},
		},
	}

	clnt, err = opensearch.NewClient(cfg)
	if err != nil {
		err = &errortypes.DatabaseError{
			errors.Wrap(err, "search: Failed to create elastic client"),
		}
		return
	}

	return
}

func hashConf(username, password string, addrs []string) []byte {
	hash := md5.New()

	io.WriteString(hash, username)
	io.WriteString(hash, password)

	for _, addr := range addrs {
		io.WriteString(hash, addr)
	}

	return hash.Sum(nil)
}

func update(username, password string, addrs []string) (err error) {
	clnt, err := newClient(username, password, addrs)
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
		Index: false,
	})

	mappings = append(mappings, mapping{
		Field: "header",
		Type:  "object",
		Index: false,
	})

	mappings = append(mappings, mapping{
		Field: "body",
		Type:  "text",
		Store: false,
		Index: false,
	})

	err = putIndex(clnt, "zero-requests", mappings)
	if err != nil {
		client = nil
		return
	}

	client = clnt

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
			err := update(username, password, addrs)
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

func sendDocs(clnt *opensearch.Client, docs []*document, log bool) {
	var err error
	var resp *opensearchapi.Response

	for i := 0; i < 5; i++ {
		reqsData := [][]byte{}

		for _, doc := range docs {
			indexReq, e := json.Marshal(&bulkIndexReq{
				Index: &bulkIndexReqData{
					Index: doc.Index,
					Id:    doc.Id,
				},
			})
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "search: Failed to marshal index"),
				}
				return
			}

			reqsData = append(reqsData, indexReq)

			indexDoc, e := json.Marshal(doc.Data)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "search: Failed to marshal doc"),
				}
				return
			}

			reqsData = append(reqsData, indexDoc)
		}

		data := bytes.Join(reqsData, []byte("\n"))

		data = append(data, []byte("\n")...)

		resp, err = clnt.Bulk(bytes.NewBuffer(data))
		if err != nil {
			err = &errortypes.DatabaseError{
				errors.Wrap(err, "search: Bulk insert failed"),
			}

			if i == 2 {
				reconnect = true
				time.Sleep(3 * time.Second)
			} else {
				time.Sleep(1 * time.Second)
			}
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			if i == 2 {
				reconnect = true
				time.Sleep(3 * time.Second)
			} else {
				time.Sleep(1 * time.Second)
			}
			continue
		}

		err = nil
		break
	}

	if err != nil || (resp != nil && resp.StatusCode != 200) {
		failedLock.Lock()
		for _, doc := range docs {
			failedBuffer.PushBack(doc)
		}
		failedLock.Unlock()

		if log {
			if resp != nil {
				respBody := ""
				respData, _ := ioutil.ReadAll(resp.Body)
				if respData != nil {
					respBody = string(respData)
				}

				err = &errortypes.DatabaseError{
					errors.New("search: Bulk insert failed"),
				}

				logrus.WithFields(logrus.Fields{
					"status_code": resp.StatusCode,
					"response":    respBody,
					"error":       err,
				}).Error("search: Bulk insert failed, moving to buffer")
			} else {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("search: Bulk insert failed, moving to buffer")
			}
		}
	}
}

func worker() {
	for {
		time.Sleep(1 * time.Second)

		group := []*document{}
		groups := [][]*document{}

		lock.Lock()
		for elem := buffer.Front(); elem != nil; elem = elem.Next() {
			if len(group) >= 100 {
				groups = append(groups, group)
				group = []*document{}
			}
			doc := elem.Value.(*document)
			group = append(group, doc)
		}
		buffer = list.New()
		lock.Unlock()

		if len(group) > 0 {
			groups = append(groups, group)
		}

		clnt := client
		if client == nil {
			continue
		}

		if len(groups) == 0 {
			continue
		}

		for _, group := range groups {
			go sendDocs(clnt, group, true)
		}
	}
}

func failedWorker() {
	for {
		time.Sleep(1 * time.Second)

		group := []*document{}
		groups := [][]*document{}

		lock.Lock()
		for elem := failedBuffer.Front(); elem != nil; elem = elem.Next() {
			if len(group) >= 10 {
				groups = append(groups, group)
				group = []*document{}
			}
			doc := elem.Value.(*document)
			group = append(group, doc)
		}
		buffer = list.New()
		lock.Unlock()

		if len(group) > 0 {
			groups = append(groups, group)
		}

		clnt := client
		if client == nil {
			continue
		}

		if len(groups) == 0 {
			continue
		}

		for _, group := range groups {
			go sendDocs(clnt, group, false)
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
