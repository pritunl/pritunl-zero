package notification

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/sirupsen/logrus"
)

var (
	clientTransport = &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
		},
	}
	client = &http.Client{
		Transport: clientTransport,
		Timeout:   15 * time.Second,
	}
)

type notificationResp struct {
	Web     bool   `json:"web"`
	Message string `json:"message"`
}

func Check() (err error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "app.pritunl.com",
		Path: fmt.Sprintf(
			"/notification/zero/%d",
			utils.GetIntVer(constants.Version),
		),
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"GET",
		u.String(),
		nil,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "notification: Request init error"),
		}
		return
	}

	req.Header.Set("User-Agent", "pritunl-zero")

	res, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "notification: Request get error"),
		}
		return
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != 200 {
		err = &errortypes.RequestError{
			errors.Newf("notification: Bad status %d", res.StatusCode),
		}
		return
	}

	data := &notificationResp{}
	err = json.NewDecoder(res.Body).Decode(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "notification: Failed to parse response body"),
		}
		return
	}

	if data.Web {
		settings.Local.DisableMsg = utils.FilterStr(data.Message, 256)
		logrus.WithFields(logrus.Fields{
			"message": settings.Local.DisableMsg,
		}).Error("notification: Disabling web server from vulnerability report")
		settings.Local.DisableWeb = true
	} else {
		settings.Local.DisableWeb = false
		settings.Local.DisableMsg = ""
	}

	return
}
