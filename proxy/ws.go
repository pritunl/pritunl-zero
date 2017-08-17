package proxy

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gorilla/websocket"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	webSocketConns     = set.NewSet()
	webSocketConnsLock = sync.Mutex{}
)

type webSocket struct {
	serverHost  string
	serverProto string
	proxyProto  string
	proxyPort   int
	upgrader    *websocket.Upgrader
}

type webSocketConn struct {
	back  *websocket.Conn
	front *websocket.Conn
}

func (w *webSocketConn) Run() {
	webSocketConnsLock.Lock()
	webSocketConns.Add(w)
	webSocketConnsLock.Unlock()

	defer func() {
		webSocketConnsLock.Lock()
		webSocketConns.Remove(w)
		webSocketConnsLock.Unlock()
	}()

	wait := make(chan bool, 2)
	go func() {
		io.Copy(w.back.UnderlyingConn(), w.front.UnderlyingConn())
		wait <- true
	}()
	go func() {
		io.Copy(w.front.UnderlyingConn(), w.back.UnderlyingConn())
		wait <- true
	}()
	<-wait
}

func (w *webSocketConn) Close() {
	if w.back != nil {
		w.back.Close()
	}
	if w.front != nil {
		w.front.Close()
	}
}

func (w *webSocket) Director(req *http.Request) (
	u *url.URL, header http.Header) {

	header = utils.CloneHeader(req.Header)
	u = &url.URL{}
	*u = *req.URL

	u.Scheme = w.serverProto
	u.Host = w.serverHost

	header.Set("X-Forwarded-For",
		strings.Split(req.RemoteAddr, ":")[0])
	header.Set("X-Forwarded-Proto", w.proxyProto)
	header.Set("X-Forwarded-Port", strconv.Itoa(w.proxyPort))

	cookie := header.Get("Cookie")
	start := strings.Index(cookie, "pritunl-zero=")
	if start != -1 {
		str := cookie[start:]
		end := strings.Index(str, ";")
		if end != -1 {
			if len(str) > end+1 && string(str[end+1]) == " " {
				end += 1
			}
			cookie = cookie[:start] + cookie[start+end+1:]
		} else {
			cookie = cookie[:start]
		}
	}

	cookie = strings.TrimSpace(cookie)

	if len(cookie) > 0 {
		header.Set("Cookie", cookie)
	} else {
		header.Del("Cookie")
	}

	return
}

func (w *webSocket) ServeHTTP(rw http.ResponseWriter, r *http.Request,
	sess *session.Session) {

	stripCookie(r)

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

	conn := &webSocketConn{
		front: frontConn,
		back:  backConn,
	}

	conn.Run()
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

func newWebSocket(proxyProto string, proxyPort int, host *Host,
	server *service.Server) (ws *webSocket) {

	ws = &webSocket{
		serverHost: fmt.Sprintf("%s:%d", server.Hostname, server.Port),
		proxyProto: proxyProto,
		proxyPort:  proxyPort,
		upgrader: &websocket.Upgrader{
			HandshakeTimeout: time.Duration(
				settings.Router.HandshakeTimeout) * time.Second,
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	if server.Protocol == "http" {
		ws.serverProto = "ws"
	} else {
		ws.serverProto = "wss"
	}

	return
}

func WebSocketsStop() {
	webSocketConnsLock.Lock()
	for socketInf := range webSocketConns.Iter() {
		func() {
			socket := socketInf.(*webSocketConn)
			socket.Close()
		}()
	}
	webSocketConns = set.NewSet()
	webSocketConnsLock.Unlock()
}
