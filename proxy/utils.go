package proxy

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/sirupsen/logrus"
)

func WriteError(w http.ResponseWriter, r *http.Request, code int, err error) {
	var msg string
	if r != nil && r.Body != nil {
		defer r.Body.Close()

		body, readErr := io.ReadAll(io.LimitReader(r.Body, 4096))
		if readErr == nil && len(body) > 0 {
			msg = string(body)
		}
	}

	if msg == "" {
		if reqErr, ok := err.(*errortypes.RequestError); ok {
			msg = fmt.Sprintf("%d: %s", code, reqErr.GetMessage())
		} else {
			msg = utils.GetStatusMessage(code)
		}
	}

	http.Error(w, msg, code)
}

func WriteErrorLog(w http.ResponseWriter, r *http.Request, code int,
	err error) {

	WriteError(w, r, code, err)

	logrus.WithFields(logrus.Fields{
		"client": node.Self.GetRemoteAddr(r),
		"error":  err,
	}).Error("proxy: Serve error")
}

func stripCookieHeaders(r *http.Request) {
	r.Header.Del("Pritunl-Zero-Token")
	r.Header.Del("Pritunl-Zero-Signature")
	r.Header.Del("Pritunl-Zero-Timestamp")
	r.Header.Del("Pritunl-Zero-Nonce")

	cookie := r.Header.Get("Cookie")
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
		r.Header.Set("Cookie", cookie)
	} else {
		r.Header.Del("Cookie")
	}
}
