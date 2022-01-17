package endpoint

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/gorilla/websocket"
)

var (
	Upgrader = websocket.Upgrader{
		HandshakeTimeout: 30 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	WebSockets     = set.NewSet()
	WebSocketsLock = sync.Mutex{}
)

type WebSocket struct {
	Conn   *websocket.Conn
	Ticker *time.Ticker
	Cancel context.CancelFunc
	Closed bool
}

func (w *WebSocket) Close() {
	w.Closed = true
	func() {
		defer func() {
			recover()
		}()
		w.Cancel()
	}()
	func() {
		defer func() {
			recover()
		}()
		w.Ticker.Stop()
	}()
	func() {
		defer func() {
			recover()
		}()
		_ = w.Conn.Close()
	}()
}

func WebSocketsStop() {
	WebSocketsLock.Lock()
	for socketInf := range WebSockets.Iter() {
		func() {
			socket := socketInf.(*WebSocket)
			socket.Close()
		}()
	}
	WebSockets = set.NewSet()
	WebSocketsLock.Unlock()
}
