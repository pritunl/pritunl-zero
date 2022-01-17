package searches

import (
	"net/http"
	"net/url"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/search"
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

func (r *Request) Index() {
	search.Index("zero-requests", r, false)
	return
}

func init() {
	mappings := []*search.Mapping{
		&search.Mapping{
			Field: "user",
			Type:  search.Keyword,
			Store: false,
			Index: true,
		},
		&search.Mapping{
			Field: "username",
			Type:  search.Keyword,
			Store: false,
			Index: true,
		},
		&search.Mapping{
			Field: "session",
			Type:  search.Keyword,
			Store: false,
			Index: true,
		},
		&search.Mapping{
			Field: "address",
			Type:  search.Ip,
			Store: false,
			Index: true,
		},
		&search.Mapping{
			Field: "timestamp",
			Type:  search.Date,
			Store: false,
			Index: true,
		},
		&search.Mapping{
			Field: "scheme",
			Type:  search.Keyword,
			Store: false,
			Index: true,
		},
		&search.Mapping{
			Field: "host",
			Type:  search.Keyword,
			Store: false,
			Index: true,
		},
		&search.Mapping{
			Field: "path",
			Type:  search.Keyword,
			Store: false,
			Index: true,
		},
		&search.Mapping{
			Field: "query",
			Type:  search.Object,
			Index: false,
		},
		&search.Mapping{
			Field: "body",
			Type:  search.Text,
			Store: false,
			Index: false,
		},
		&search.Mapping{
			Field: "header.User-Agent",
			Type:  search.Keyword,
			Store: false,
			Index: true,
		},
		&search.Mapping{
			Field: "header.X-Forwarded-User",
			Type:  search.Keyword,
			Store: false,
			Index: true,
		},
	}

	search.AddMappings("zero-requests", mappings)

	return
}
