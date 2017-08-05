package proxy

import (
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/gorilla/websocket"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"io"
	"net/http"
	"net/url"
	"time"
)

type WebSocket struct {
	Director func(*http.Request) (*url.URL, http.Header)
	upgrader *websocket.Upgrader
}

func (w *WebSocket) Init() {
	w.upgrader = &websocket.Upgrader{
		HandshakeTimeout: time.Duration(
			settings.Router.HandshakeTimeout) * time.Second,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func (w *WebSocket) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	u, header := w.Director(r)

	header.Del("Upgrade")
	header.Del("Connection")
	header.Del("Sec-Websocket-Key")
	header.Del("Sec-Websocket-Version")
	header.Del("Sec-Websocket-Extensions")

	backConn, backResp, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "proxy: WebSocket dial error"),
		}
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("proxy: WebSocket dial error")
	}
	defer backConn.Close()

	upgradeHeaders := getUpgradeHeaders(backResp)
	frontConn, err := w.upgrader.Upgrade(rw, r, upgradeHeaders)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "proxy: WebSocket upgrade error"),
		}
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("proxy: WebSocket upgrade error")
	}
	defer frontConn.Close()

	wait := make(chan bool, 2)
	go func() {
		io.Copy(backConn.UnderlyingConn(), frontConn.UnderlyingConn())
		wait <- true
	}()
	go func() {
		io.Copy(frontConn.UnderlyingConn(), backConn.UnderlyingConn())
		wait <- true
	}()
	<-wait
}

func getUpgradeHeaders(resp *http.Response) (header http.Header) {
	header = http.Header{}

	val := resp.Header.Get("Sec-Websocket-Protocol")
	if val != "" {
		header.Set("Sec-Websocket-Protocol", val)
	}

	val = resp.Header.Get("Set-Cookie")
	if val != "" {
		header.Set("Set-Cookie", val)
	}

	return
}
