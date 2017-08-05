package proxy

import (
	"github.com/Sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/utils"
	"net/http"
)

func WriteError(w http.ResponseWriter, r *http.Request, code int, err error) {
	http.Error(w, utils.GetStatusMessage(code), code)

	logrus.WithFields(logrus.Fields{
		"client": node.Self.GetRemoteAddr(r),
		"error":  err,
	}).Error("proxy: Serve error")
}
