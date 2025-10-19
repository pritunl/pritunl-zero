package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/sirupsen/logrus"
)

const (
	JumpCloud = "jumpcloud"
)

type jumpcloudResponse struct {
	Results    []*jumpcloudUser `json:"results"`
	TotalCount int              `json:"totalCount"`
}

type jumpcloudUser struct {
	Id            string `json:"id"`
	IdAlt         string `json:"_id"`
	Email         string `json:"email"`
	AccountLocked bool   `json:"account_locked"`
	Suspended     bool   `json:"suspended"`
	Activated     bool   `json:"activated"`
}

type jumpcloudApp struct {
	Id string `json:"id"`
}

func jumpcloudCheckApp(provider *settings.Provider, userId string) (
	attached bool, err error) {

	reqUrl := &url.URL{
		Scheme: "https",
		Host:   "console.jumpcloud.com",
		Path:   fmt.Sprintf("/api/v2/users/%s/applications", userId),
	}

	req, err := http.NewRequest(
		"GET",
		reqUrl.String(),
		nil,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Failed to create jumpcloud request"),
		}
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Api-Key", provider.JumpCloudSecret)

	resp, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Jumpcloud request failed"),
		}
		return
	}
	defer resp.Body.Close()

	err = utils.CheckRequest(resp, "auth: Auth server error")
	if err != nil {
		return
	}

	data := []*jumpcloudApp{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse jumpcloud response"),
		}
		return
	}

	for _, app := range data {
		if app.Id == provider.JumpCloudAppId {
			attached = true
			return
		}
	}

	return
}

func JumpcloudSync(db *database.Database, usr *user.User,
	provider *settings.Provider) (active bool, err error) {

	reqUrl := &url.URL{
		Scheme: "https",
		Host:   "console.jumpcloud.com",
		Path:   "/api/systemusers",
	}

	query := reqUrl.Query()
	query.Set("filter", fmt.Sprintf("email:$eq:%s", usr.Username))

	reqUrl.RawQuery = query.Encode()

	req, err := http.NewRequest(
		"GET",
		reqUrl.String(),
		nil,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Failed to create jumpcloud request"),
		}
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Api-Key", provider.JumpCloudSecret)

	resp, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Jumpcloud request failed"),
		}
		return
	}
	defer resp.Body.Close()

	err = utils.CheckRequest(resp, "auth: Auth server error")
	if err != nil {
		return
	}

	data := &jumpcloudResponse{}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse jumpcloud response"),
		}
		return
	}

	userId := ""
	if data.TotalCount > 0 && data.Results != nil {
		for _, authUser := range data.Results {
			if authUser.Email != usr.Username {
				continue
			}

			if authUser.AccountLocked || authUser.Suspended ||
				!authUser.Activated {

				logrus.WithFields(logrus.Fields{
					"user_id":  usr.Id.Hex(),
					"username": usr.Username,
				}).Info("auth: Jumpcloud user disabled")

				return
			} else {
				if authUser.Id != "" {
					userId = authUser.Id
				} else {
					userId = authUser.IdAlt
				}
				break
			}
		}
	}

	if userId == "" {
		err = &errortypes.NotFoundError{
			errors.Wrap(err, "auth: Jumpcloud user not found"),
		}
		return
	}

	if provider.JumpCloudAppId != "" {
		attached, e := jumpcloudCheckApp(provider, userId)
		if e != nil {
			err = e
			return
		}

		if attached {
			active = true
		} else {
			logrus.WithFields(logrus.Fields{
				"user_id":  usr.Id.Hex(),
				"username": usr.Username,
			}).Info("auth: Jumpcloud user not bound to application")

			return
		}
	} else {
		active = true
	}

	return
}
