package proxy

import (
	"net/http"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/node"
)

type TransportFix struct {
	transport *http.Transport
}

func (t *TransportFix) RoundTrip(r *http.Request) (
	res *http.Response, err error) {

	r.Header.Set("X-Forwarded-For", node.Self.GetRemoteAddr(r))

	res, err = t.transport.RoundTrip(r)
	if err != nil {
		return
	}

	if res.StatusCode == http.StatusSwitchingProtocols {
		err = &WebSocketBlock{
			errors.New("proxy: Blocking websocket connection"),
		}
		return
	}

	return
}
