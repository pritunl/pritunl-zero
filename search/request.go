package search

import (
	"net/http"
	"net/url"
	"time"

	"github.com/dropbox/godropbox/container/set"
)

var (
	RequestTypes = set.NewSet(
		"application/xml",
		"application/json",
		"application/x-www-form-urlencoded",
	)
)

type Request struct {
	User      string      `json:"user"`
	Username  string      `json:"username"`
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
	Index("zero-requests", r)
	return
}
