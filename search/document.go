package search

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/sirupsen/logrus"
)

type Document struct {
	Index    string
	Id       string
	Data     interface{}
	attempts int
}

type searchBulkReqData struct {
	Index string `json:"_index"`
	Id    string `json:"_id"`
}

type searchBulkReq struct {
	Index *searchBulkReqData `json:"index"`
}

func Index(index string, data interface{}) {
	clnt := Default

	if clnt == nil {
		return
	}

	doc := &Document{
		Index: index + dateSuffix(),
		Id:    primitive.NewObjectID().Hex(),
		Data:  data,
	}

	lock.Lock()
	buffer.PushBack(doc)
	lock.Unlock()

	return
}

func (c *Client) BulkDocuments(docs []*Document, log bool) (err error) {
	var resp *opensearchapi.Response

	for i := 0; i < 3; i++ {
		reqsData := [][]byte{}

		for _, doc := range docs {
			if !c.indexes.Contains(doc.Index) {
				err = c.AddIndex(doc.Index)
				if err != nil {
					return
				}
			}

			doc.attempts += 1

			indexReq, e := json.Marshal(&searchBulkReq{
				Index: &searchBulkReqData{
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

		resp, err = c.clnt.Bulk(bytes.NewBuffer(data))
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

		if resp != nil {
			respBody := ""
			respData, _ := ioutil.ReadAll(resp.Body)
			if respData != nil {
				respBody = string(respData)
			}

			err = &errortypes.DatabaseError{
				errors.New("search: Bulk insert failed"),
			}

			if log {
				logrus.WithFields(logrus.Fields{
					"status_code": resp.StatusCode,
					"response":    respBody,
					"error":       err,
				}).Error("search: Bulk insert failed, moving to buffer")
			}
		} else if log {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("search: Bulk insert failed, moving to buffer")
		}
		return
	}

	return
}
