package search

import (
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	User      string      `json:"user"`
	Session   string      `json:"session"`
	Address   string      `json:"address"`
	Timestamp time.Time   `json:"timestamp"`
	Scheme    string      `json:"scheme"`
	Host      string      `json:"host"`
	Path      string      `json:"path"`
	Query     url.Values  `json:"query"`
	Header    http.Header `json:"header"`
	Body      string      `json:"body"`
}

func (r *Request) Index() (err error) {
	err = Index("zero-requests", "request", r)
	if err != nil {
		return
	}

	return
}
