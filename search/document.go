package search

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/sirupsen/logrus"
)

type Document struct {
	Index    string
	Id       string
	Data     []byte
	Size     int
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

	dataJson, err := json.Marshal(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "search: Failed to marshal doc"),
		}
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("search: Failed to marshal index")
		return
	}

	doc := &Document{
		Index:   index + dateSuffix(),
		Id:      bson.NewObjectID().Hex(),
		Data:    dataJson,
		NoRetry: noRetry,
	}

	if len(buffer) <= settings.Elastic.BufferLength {
		select {
		case buffer <- doc:
		default:
		}
	}

	return
}

func (c *Client) BulkDocuments(docs *IndexList, log bool) (err error) {
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

			reqsData.Write(doc.Data)
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

			if len(failedBuffer) <= settings.Elastic.BufferLength {
				select {
				case failedBuffer <- doc:
				default:
				}
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
