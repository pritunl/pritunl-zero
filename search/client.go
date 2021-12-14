package search

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/opensearch-project/opensearch-go"
	"github.com/pritunl/pritunl-zero/errortypes"
)

type Client struct {
	clnt    *opensearch.Client
	indexes set.Set
}

func NewClient(username, password string, addrs []string) (
	c *Client, err error) {

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

	clnt, err := opensearch.NewClient(cfg)
	if err != nil {
		err = &errortypes.DatabaseError{
			errors.Wrap(err, "search: Failed to create elastic client"),
		}
		return
	}

	c = &Client{
		clnt:    clnt,
		indexes: set.NewSet(),
	}

	return
}
