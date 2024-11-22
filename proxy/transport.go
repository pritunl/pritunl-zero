package proxy

import (
	"net/http"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/node"
	"golang.org/x/net/http2"
)

type TransportFix struct {
	transport  *http.Transport
	transport2 *http2.Transport
}

func (t *TransportFix) RoundTrip(r *http.Request) (
	res *http.Response, err error) {

	r.Header.Set("X-Forwarded-For", node.Self.GetRemoteAddr(r))

	if t.transport2 != nil {
		res, err = t.transport2.RoundTrip(r)
		if err != nil {
			return
		}
	} else {
		res, err = t.transport.RoundTrip(r)
		if err != nil {
			return
		}
	}

	if res.StatusCode == http.StatusSwitchingProtocols {
		err = &WebSocketBlock{
			errors.New("proxy: Blocking websocket connection"),
		}
		return
	}

	return
}
