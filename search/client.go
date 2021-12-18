package search

import (
	"crypto/tls"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/opensearch-project/opensearch-go"
	"github.com/pritunl/pritunl-zero/errortypes"
)

var (
	ip4reg = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	ip6reg = regexp.MustCompile("/\\[[a-fA-F0-9:]*\\]/")
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

	skipVerify := false
	for _, addr := range addrs {
		if len(ip4reg.FindAllString(addr, -1)) > 0 ||
			len(ip6reg.FindAllString(addr, -1)) > 0 {

			skipVerify = true
			break
		}
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
				InsecureSkipVerify: skipVerify,
				MinVersion:         tls.VersionTLS12,
				MaxVersion:         tls.VersionTLS13,
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
