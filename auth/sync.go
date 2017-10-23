package auth

import (
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"net/http"
	"net/url"
	"time"
)

func SyncUser(db *database.Database, usr *user.User) (
	active bool, err error) {

	if time.Since(usr.LastSync) < time.Duration(
		settings.Auth.Sync)*time.Second {

		active = true
		return
	}

	if usr.Type == user.Google {
		reqVals := url.Values{}
		reqVals.Set("user", usr.Username)
		reqVals.Set("license", settings.System.License)

		reqUrl, _ := url.Parse(settings.Auth.Server + "/update/google")
		reqUrl.RawQuery = reqVals.Encode()

		req, e := http.NewRequest(
			"GET",
			reqUrl.String(),
			nil,
		)
		if e != nil {
			err = &errortypes.RequestError{
				errors.Wrap(e, "auth: Auth request failed"),
			}
			return
		}

		resp, e := client.Do(req)
		if e != nil {
			err = &errortypes.RequestError{
				errors.Wrap(e, "auth: Auth request failed"),
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			active = true

			usr.LastSync = time.Now()
			err = usr.CommitFields(db, set.NewSet("last_sync"))
			if err != nil {
				return
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"username":    usr.Username,
				"status_code": resp.StatusCode,
			}).Info("session: User single sign-on sync failed")
		}
	} else {
		active = true
	}

	return
}
