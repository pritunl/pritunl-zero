package proxy

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gorilla/websocket"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/searches"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/pritunl/pritunl-zero/validator"
	"github.com/sirupsen/logrus"
)

var (
	webSocketConns     = set.NewSet()
	webSocketConnsLock = sync.Mutex{}
)

type webSocket struct {
	reqHost     string
	serverHost  string
	serverProto string
	proxyProto  string
	proxyPort   int
	tlsConfig   *tls.Config
	upgrader    *websocket.Upgrader
}

type webSocketConn struct {
	authr *authorizer.Authorizer
	r     *http.Request
	back  *websocket.Conn
	front *websocket.Conn
}

func (w *webSocketConn) Run(db *database.Database) {
	webSocketConnsLock.Lock()
	webSocketConns.Add(w)
	webSocketConnsLock.Unlock()

	defer func() {
		webSocketConnsLock.Lock()
		webSocketConns.Remove(w)
		webSocketConnsLock.Unlock()
	}()

	ticker := time.NewTicker(30 * time.Second)
	closer := make(chan bool, 1)
	waiter := sync.WaitGroup{}
	waiter.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithFields(logrus.Fields{
					"error": errors.New(fmt.Sprintf("%s", r)),
				}).Error("proxy: WebSocket update panic")
				w.Close()
			}
		}()
		defer func() {
			waiter.Done()
		}()

		for {
			select {
			case <-ticker.C:
				if w.authr.IsValid() {
					usr, err := w.authr.GetUser(db)
					if err != nil {
						switch err.(type) {
						case *database.NotFoundError:
							break
						default:
							logrus.WithFields(logrus.Fields{
								"error": err,
							}).Error("proxy: WebSocket user error")
						}
						w.Close()
						return
					}

					sess := w.authr.GetSession()
					if sess != nil {
						err = sess.Update(db)
						if err != nil {
							switch err.(type) {
							case *database.NotFoundError:
								break
							default:
								logrus.WithFields(logrus.Fields{
									"error": err,
								}).Error("proxy: WebSocket session error")
							}
							w.Close()
							return
						}

						if !sess.Active() {
							w.Close()
							return
						}
					}

					srvcId := w.authr.ServiceId()
					if !srvcId.IsZero() {
						srvc, err := service.Get(db, srvcId)
						if err != nil {
							switch err.(type) {
							case *database.NotFoundError:
								break
							default:
								logrus.WithFields(logrus.Fields{
									"error": err,
								}).Error("proxy: WebSocket service error")
							}
							w.Close()
							return
						}

						_, _, _, errData, err := validator.ValidateProxy(
							db, usr, w.authr.IsApi(), srvc, w.r)
						if err != nil {
							logrus.WithFields(logrus.Fields{
								"error": err,
							}).Error("proxy: WebSocket validate error")
							w.Close()
							return
						}

						if errData != nil {
							w.Close()
							return
						}
					}
				}

				break
			case <-closer:
				return
			}
		}
	}()

	wait := make(chan bool, 4)
	go func() {
		defer func() {
			rec := recover()
			if rec != nil {
				logrus.WithFields(logrus.Fields{
					"panic": rec,
				}).Error("proxy: WebSocket back panic")
				wait <- true
			}
		}()

		for {
			msgType, msg, err := w.back.ReadMessage()
			if err != nil {
				closeMsg := websocket.FormatCloseMessage(
					websocket.CloseNormalClosure, fmt.Sprintf("%v", err))
				if e, ok := err.(*websocket.CloseError); ok {
					if e.Code != websocket.CloseNoStatusReceived {
						closeMsg = websocket.FormatCloseMessage(e.Code, e.Text)
					}
				}
				_ = w.front.WriteMessage(websocket.CloseMessage, closeMsg)
				break
			}

			_ = w.front.WriteMessage(msgType, msg)
		}

		wait <- true
	}()
	go func() {
		defer func() {
			rec := recover()
			if rec != nil {
				logrus.WithFields(logrus.Fields{
					"panic": rec,
				}).Error("proxy: WebSocket front panic")
				wait <- true
			}
		}()

		for {
			msgType, msg, err := w.front.ReadMessage()
			if err != nil {
				closeMsg := websocket.FormatCloseMessage(
					websocket.CloseNormalClosure, fmt.Sprintf("%v", err))
				if e, ok := err.(*websocket.CloseError); ok {
					if e.Code != websocket.CloseNoStatusReceived {
						closeMsg = websocket.FormatCloseMessage(e.Code, e.Text)
					}
				}
				_ = w.back.WriteMessage(websocket.CloseMessage, closeMsg)
				break
			}

			_ = w.back.WriteMessage(msgType, msg)
		}

		wait <- true
	}()
	<-wait

	ticker.Stop()
	closer <- true
	w.Close()
	waiter.Wait()
}

func (w *webSocketConn) Close() {
	defer func() {
		recover()
	}()
	if w.back != nil {
		_ = w.back.Close()
	}
	if w.front != nil {
		_ = w.front.Close()
	}
}

func (w *webSocket) Director(req *http.Request, authr *authorizer.Authorizer) (
	u *url.URL, header http.Header) {

	header = utils.CloneHeader(req.Header)
	u = &url.URL{}
	*u = *req.URL

	u.Scheme = w.serverProto
	u.Host = w.serverHost

	header.Set("X-Forwarded-For",
		node.Self.GetRemoteAddr(req))
	header.Set("X-Forwarded-Host", req.Host)
	header.Set("X-Forwarded-Proto", w.proxyProto)
	header.Set("X-Forwarded-Port", strconv.Itoa(w.proxyPort))

	if authr != nil {
		usr, _ := authr.GetUser(nil)
		if usr != nil {
			req.Header.Set("X-Forwarded-User", usr.Username)
		}
	}

	header.Del("Upgrade")
	header.Del("Connection")
	header.Del("Sec-Websocket-Key")
	header.Del("Sec-Websocket-Version")
	header.Del("Sec-Websocket-Extensions")

	stripCookieHeaders(req)

	return
}

func (w *webSocket) ServeHTTP(rw http.ResponseWriter, r *http.Request,
	db *database.Database, authr *authorizer.Authorizer) {

	u, header := w.Director(r, authr)

	scheme := ""
	if u.Scheme == "https" {
		scheme = "wss"
	} else {
		scheme = "ws"
	}

	if settings.Elastic.ProxyRequests {
		index := searches.Request{
			Address:   node.Self.GetRemoteAddr(r),
			Timestamp: time.Now(),
			Scheme:    scheme,
			Host:      u.Host,
			Path:      r.URL.Path,
			Query:     r.URL.Query(),
			Header:    r.Header.Clone(),
		}

		if authr.IsValid() {
			usr, _ := authr.GetUser(nil)

			if usr != nil {
				index.User = usr.Id.Hex()
				index.Username = usr.Username
				index.Session = authr.SessionId()
			}
		}

		index.Index()
	}

	var backConn *websocket.Conn
	var backResp *http.Response
	var err error

	dialer := &websocket.Dialer{
		Proxy: func(req *http.Request) (url *url.URL, err error) {
			req.Header.Set("X-Forwarded-For",
				node.Self.GetRemoteAddr(r))
			req.Header.Set("X-Forwarded-Host", r.Host)
			req.Header.Set("X-Forwarded-Proto", w.proxyProto)
			req.Header.Set("X-Forwarded-Port", strconv.Itoa(w.proxyPort))

			if authr != nil {
				usr, _ := authr.GetUser(nil)
				if usr != nil {
					req.Header.Set("X-Forwarded-User", usr.Username)
				}
			}

			if w.reqHost != "" {
				req.Host = w.reqHost
			} else {
				req.Host = r.Host
			}
			return
		},
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  w.tlsConfig,
	}

	backConn, backResp, err = dialer.Dial(u.String(), header)
	if err != nil {
		if backResp != nil {
			err = &errortypes.RequestError{
				errors.Wrapf(err, "proxy: WebSocket dial error %d",
					backResp.StatusCode),
			}
		} else {
			err = &errortypes.RequestError{
				errors.Wrap(err, "proxy: WebSocket dial error"),
			}
		}
		WriteErrorLog(rw, r, 500, err)
		return
	}
	defer func() {
		_ = backConn.Close()
	}()

	upgradeHeaders := getUpgradeHeaders(backResp)
	frontConn, err := w.upgrader.Upgrade(rw, r, upgradeHeaders)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "proxy: WebSocket upgrade error"),
		}
		WriteErrorLog(rw, r, 500, err)
		return
	}
	defer func() {
		_ = frontConn.Close()
	}()

	conn := &webSocketConn{
		front: frontConn,
		back:  backConn,
		authr: authr,
		r:     r,
	}

	conn.Run(db)
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

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
	}
	if settings.Router.SkipVerify || net.ParseIP(server.Hostname) != nil {
		tlsConfig.InsecureSkipVerify = true
	}

	if host.ClientCertificate != nil {
		tlsConfig.Certificates = []tls.Certificate{
			*host.ClientCertificate,
		}
	}

	ws = &webSocket{
		reqHost:    host.Domain.Host,
		serverHost: utils.FormatHostPort(server.Hostname, server.Port),
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
		tlsConfig: tlsConfig,
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
