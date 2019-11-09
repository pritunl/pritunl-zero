package proxy

import (
	"net/http"

	"github.com/pritunl/pritunl-zero/node"
)

type TransportFix struct {
	transport *http.Transport
}

func (t *TransportFix) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("X-Forwarded-For", node.Self.GetRemoteAddr(r))
	return t.transport.RoundTrip(r)
}
