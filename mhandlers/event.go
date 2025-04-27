package mhandlers

import (
	"context"
	"fmt"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/sirupsen/logrus"
)

const (
	writeTimeout = 10 * time.Second
	pingInterval = 30 * time.Second
	pingWait     = 40 * time.Second
)

func eventGet(c *gin.Context) {
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
			errors.Wrap(err, "mhandlers: Failed to upgrade request"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}
	socket.Conn = conn

	err = conn.SetReadDeadline(time.Now().Add(pingWait))
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "mhandlers: Failed to set read deadline"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	conn.SetPongHandler(func(x string) (err error) {
		err = conn.SetReadDeadline(time.Now().Add(pingWait))
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(err, "mhandlers: Failed to set read deadline"),
			}
			utils.AbortWithError(c, 500, err)
			return
		}

		return
	})

	lst, err := event.SubscribeListener(db, []string{"dispatch"})
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}
	socket.Listener = lst

	ticker := time.NewTicker(pingInterval)
	socket.Ticker = ticker
	sub := lst.Listen()
	defer lst.Close()

	go func() {
		defer func() {
			r := recover()
			if r != nil && !socket.Closed {
				logrus.WithFields(logrus.Fields{
					"error": errors.New(fmt.Sprintf("%s", r)),
				}).Error("mhandlers: Event panic")
			}
		}()
		for {
			_, _, err := conn.NextReader()
			if err != nil {
				_ = conn.Close()
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-sub:
			if !ok {
				err = conn.WriteControl(websocket.CloseMessage, []byte{},
					time.Now().Add(writeTimeout))
				if err != nil {
					err = &errortypes.RequestError{
						errors.Wrap(err,
							"mhandlers: Failed to set write control"),
					}
					return
				}

				return
			}

			err = conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err != nil {
				err = &errortypes.RequestError{
					errors.Wrap(err,
						"mhandlers: Failed to set write deadline"),
				}
				return
			}

			err = conn.WriteJSON(msg)
			if err != nil {
				err = &errortypes.RequestError{
					errors.Wrap(err,
						"mhandlers: Failed to set write json"),
				}
				return
			}
		case <-ticker.C:
			err = conn.WriteControl(websocket.PingMessage, []byte{},
				time.Now().Add(writeTimeout))
			if err != nil {
				err = &errortypes.RequestError{
					errors.Wrap(err,
						"mhandlers: Failed to set write control"),
				}
				return
			}
		}
	}
}
