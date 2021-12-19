package search

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/sirupsen/logrus"
)

const (
	Keyword = "keyword"
	Ip      = "ip"
	Date    = "date"
	Text    = "text"
	Object  = "object"
)

type Mapping struct {
	Field string
	Type  string
	Store bool
	Index bool
}

type searchMappingData struct {
	Type  string `json:"type"`
	Store bool   `json:"store"`
	Index bool   `json:"index"`
}

type searchMappingObject struct {
	Enabled bool `json:"enabled"`
}

type searchMappings map[string]interface{}

type searchProperties struct {
	Properties searchMappings `json:"properties"`
}

type searchIndexReq struct {
	Mappings *searchProperties `json:"mappings"`
}

func (c *Client) UpdateIndexes() (err error) {
	for index, mappings := range mappingsRegistry {
		err = c.CreateIndex(index+dateSuffix(), mappings)
		if err != nil {
			return
		}
	}

	return
}

func (c *Client) AddIndex(index string) (err error) {
	baseIndexSpl := strings.Split(index, "-")
	if len(baseIndexSpl) < 4 {
		err = &errortypes.DatabaseError{
			errors.Newf("search: Invalid index %s", index),
		}
		return
	}

	baseIndex := strings.Join(baseIndexSpl[:len(baseIndexSpl)-3], "-")

	mappings := mappingsRegistry[baseIndex]
	if mappings == nil {
		err = &errortypes.DatabaseError{
			errors.Newf("search: Mappings not found %s", baseIndex),
		}
		return
	}

	err = c.CreateIndex(index, mappings)
	if err != nil {
		return
	}

	return
}

func (c *Client) CreateIndex(index string, mappings []*Mapping) (err error) {
	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)

	exists, err := c.clnt.Indices.Exists([]string{index})
	if err != nil {
		err = &errortypes.DatabaseError{
			errors.Wrap(err, "search: Failed to check elastic index"),
		}
		return
	}

	if exists.StatusCode == 200 {
		c.indexes.Add(index)
		return
	}

	data := &searchIndexReq{
		Mappings: &searchProperties{
			Properties: searchMappings{},
		},
	}

	for _, mapping := range mappings {
		if strings.Contains(mapping.Field, ".") {
			continue
		}

		if mapping.Type == "object" {
			data.Mappings.Properties[mapping.Field] = &searchMappingObject{
				Enabled: mapping.Index,
			}
		} else {
			data.Mappings.Properties[mapping.Field] = &searchMappingData{
				Type:  mapping.Type,
				Store: mapping.Store,
				Index: mapping.Index,
			}
		}
	}

	for _, mapping := range mappings {
		if !strings.Contains(mapping.Field, ".") {
			continue
		}

		fields := strings.Split(mapping.Field, ".")

		curMapping := data.Mappings.Properties
		for _, field := range fields[:len(fields)-1] {
			if newProp, ok := curMapping[field].(searchProperties); ok {
				curMapping = newProp.Properties
			} else {
				newProp := &searchProperties{
					Properties: searchMappings{},
				}
				curMapping[field] = newProp
				curMapping = newProp.Properties
			}
		}

		field := fields[len(fields)-1]

		if mapping.Type == "object" {
			curMapping[field] = &searchMappingObject{
				Enabled: mapping.Index,
			}
		} else {
			curMapping[field] = &searchMappingData{
				Type:  mapping.Type,
				Store: mapping.Store,
				Index: mapping.Index,
			}
		}
	}

	body, err := json.Marshal(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "search: Failed to marshal data"),
		}
		return
	}

	logrus.WithFields(logrus.Fields{
		"index": index,
	}).Info("search: Create index")

	resp, err := c.clnt.Indices.Create(
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

		err = &errortypes.DatabaseError{
			errors.New("search: Failed to create elastic index"),
		}

		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"response":    respBody,
			"error":       err,
		}).Error("search: Failed to create elastic index")

		return
	}

	c.indexes.Add(index)

	return
}
