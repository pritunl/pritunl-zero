package search

import (
	"bytes"
	"container/list"
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
	NoRetry  bool
	attempts int
}

type searchBulkReqData struct {
	Index string `json:"_index"`
	Id    string `json:"_id"`
}

type searchBulkReq struct {
	Index *searchBulkReqData `json:"index"`
}

func Index(index string, data interface{}, noRetry bool) {
	clnt := Default

	if clnt == nil {
		return
	}

	doc := &Document{
		Index:   index + dateSuffix(),
		Id:      primitive.NewObjectID().Hex(),
		Data:    data,
		NoRetry: noRetry,
	}

	if len(buffer) <= BufferSize {
		buffer <- doc
	}

	return
}

func (c *Client) BulkDocuments(docs *list.List, log bool) (err error) {
	var resp *opensearchapi.Response

	for i := 0; i < 3; i++ {
		reqsData := bytes.NewBuffer(nil)

		for elem := docs.Front(); elem != nil; elem = elem.Next() {
			doc := elem.Value.(*Document)

			if !c.indexes.Contains(doc.Index) {
				err = c.AddIndex(doc.Index)
				if err != nil {
					return
				}
			}

			if doc.attempts >= RetryCount {
				continue
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

			reqsData.Write(indexReq)
			reqsData.Write([]byte("\n"))

			indexDoc, e := json.Marshal(doc.Data)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "search: Failed to marshal doc"),
				}
				return
			}

			reqsData.Write(indexDoc)
			reqsData.Write([]byte("\n"))
		}

		resp, err = c.clnt.Bulk(reqsData)
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
		for elem := docs.Front(); elem != nil; elem = elem.Next() {
			doc := elem.Value.(*Document)

			if doc.attempts >= RetryCount || doc.NoRetry {
				continue
			}

			if len(failedBuffer) <= BufferSize {
				failedBuffer <- doc
			}
		}

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
				if logLimit() {
					logrus.WithFields(logrus.Fields{
						"status_code": resp.StatusCode,
						"response":    respBody,
						"error":       err,
					}).Error("search: Bulk insert failed, moving to buffer")
				}
				err = nil
			}
		} else if log {
			if logLimit() {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("search: Bulk insert failed, moving to buffer")
			}
			err = nil
		}
		return
	}

	return
}
