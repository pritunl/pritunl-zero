package event

import (
	"context"
	"github.com/dropbox/godropbox/container/set"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
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
	Conn     *websocket.Conn
	Ticker   *time.Ticker
	Listener *Listener
	Cancel   context.CancelFunc
}

func (w *WebSocket) Close() {
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
		w.Listener.Close()
	}()
	func() {
		defer func() {
			recover()
		}()
		w.Conn.Close()
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
