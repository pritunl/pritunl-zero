package uhandlers

import (
	"context"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/utils"
	"time"
)

const (
	writeTimeout = 10 * time.Second
	pingInterval = 30 * time.Second
	pingWait     = 40 * time.Second
)

func hsmGet(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)
	socket := &event.WebSocket{}

	defer func() {
		socket.Close()
		event.WebSocketsLock.Lock()
		event.WebSockets.Remove(socket)
		event.WebSocketsLock.Unlock()
	}()

	event.WebSocketsLock.Lock()
	event.WebSockets.Add(socket)
	event.WebSocketsLock.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	socket.Cancel = cancel

	conn, err := event.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "uhandlers: Failed to upgrade hsm request"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}
	socket.Conn = conn

	conn.SetReadDeadline(time.Now().Add(pingWait))
	conn.SetPongHandler(func(x string) (err error) {
		conn.SetReadDeadline(time.Now().Add(pingWait))
		return
	})

	lst, err := event.SubscribeListener([]string{"pritunl_hsm_send"})
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}
	socket.Listener = lst

	ticker := time.NewTicker(pingInterval)
	socket.Ticker = ticker
	sub := lst.Listen()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithFields(logrus.Fields{
					"error": errors.New(fmt.Sprintf("%s", r)),
				}).Error("uhandlers: Socket hsm panic")
			}
		}()
		for {
			_, message, e := conn.ReadMessage()
			if e != nil {
				logrus.WithFields(logrus.Fields{
					"error": e,
				}).Error("uhandlers: Socket hsm listen error")
				conn.Close()
				break
			}

			e = event.Publish(db, "pritunl_hsm_recv", message)
			if e != nil {
				logrus.WithFields(logrus.Fields{
					"error": e,
				}).Error("uhandlers: Socket hsm publish event error")
				continue
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-sub:
			if !ok {
				conn.WriteControl(websocket.CloseMessage, []byte{},
					time.Now().Add(writeTimeout))
				return
			}

			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			err = conn.WriteJSON(msg.Data)
			if err != nil {
				return
			}
		case <-ticker.C:
			err = conn.WriteControl(websocket.PingMessage, []byte{},
				time.Now().Add(writeTimeout))
			if err != nil {
				return
			}
		}
	}
}
