package uhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pritunl/pritunl-zero/authority"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
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

func hsmGet(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	authr := c.MustGet("authority").(*authority.Authority)

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

	err = conn.SetReadDeadline(time.Now().Add(pingWait))
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "uhandlers: Failed to set read deadline"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	conn.SetPongHandler(func(x string) (err error) {
		err = conn.SetReadDeadline(time.Now().Add(pingWait))
		if err != nil {
			return
		}

		return
	})

	lst, err := event.SubscribeListener(db, []string{"pritunl_hsm_send"})
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}
	socket.Listener = lst

	ticker := time.NewTicker(pingInterval)
	socket.Ticker = ticker
	sub := lst.Listen()
	defer lst.Close()

	authr.HsmStatus = authority.Connected
	authr.HsmTimestamp = time.Now()
	err = authr.CommitFields(db, set.NewSet("hsm_status", "hsm_timestamp"))
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "authority.change")

	go func() {
		defer func() {
			r := recover()
			if r != nil && !socket.Closed {
				logrus.WithFields(logrus.Fields{
					"error": errors.New(fmt.Sprintf("%s", r)),
				}).Error("uhandlers: Socket hsm panic")
			}
		}()

		lstDb := database.GetDatabase()
		defer lstDb.Close()

		for {
			_, message, e := conn.ReadMessage()
			if e != nil {
				logrus.WithFields(logrus.Fields{
					"error": e,
				}).Error("uhandlers: Socket hsm listen error")

				authr.HsmStatus = authority.Disconnected
				_ = authr.CommitFields(db, set.NewSet("hsm_status"))
				_ = event.PublishDispatch(db, "authority.change")

				_ = conn.Close()
				break
			}

			payload := &authority.HsmPayload{}
			e = json.Unmarshal(message, payload)
			if e != nil {
				logrus.WithFields(logrus.Fields{
					"error": e,
				}).Error("uhandlers: Failed to unmarshal hsm payload")
				continue
			}

			if payload.Type == "status" {
				e = authr.HandleHsmStatus(lstDb, payload)
				if e != nil {
					logrus.WithFields(logrus.Fields{
						"error": e,
					}).Error("uhandlers: Failed to handle hsm status")
					continue
				}
			} else {
				e = event.Publish(lstDb, "pritunl_hsm_recv", payload)
				if e != nil {
					logrus.WithFields(logrus.Fields{
						"error": e,
					}).Error("uhandlers: Socket hsm publish event error")
					continue
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-sub:
			if !ok {
				_ = conn.WriteControl(websocket.CloseMessage, []byte{},
					time.Now().Add(writeTimeout))
				return
			}

			err = conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err != nil {
				return
			}

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
