package auth

import (
	"bytes"
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"net/http"
	"time"
)

const (
	Google = "google"
)

func GoogleRequest(db *database.Database, location string,
	provider *settings.Provider) (redirect string, err error) {

	if provider.Type != Google {
		err = &errortypes.ParseError{
			errors.New("auth: Invalid provider type"),
		}
		return
	}

	coll := db.Tokens()

	state, err := utils.RandStr(64)
	if err != nil {
		return
	}

	secret, err := utils.RandStr(64)
	if err != nil {
		return
	}

	data, err := json.Marshal(struct {
		License  string `json:"license"`
		Callback string `json:"callback"`
		State    string `json:"state"`
		Secret   string `json:"secret"`
	}{
		License:  settings.System.License,
		Callback: location + "/auth/callback",
		State:    state,
		Secret:   secret,
	})

	req, err := http.NewRequest(
		"POST",
		settings.Auth.Server+"/v1/request/google",
		bytes.NewBuffer(data),
	)
	if err != nil {
		err = errortypes.RequestError{
			errors.Wrap(err, "auth: Auth request failed"),
		}
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		err = errortypes.RequestError{
			errors.Wrap(err, "auth: Auth request failed"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errortypes.RequestError{
			errors.Wrapf(err, "auth: Auth server error %d", resp.StatusCode),
		}
		return
	}

	authData := &authData{}
	err = json.NewDecoder(resp.Body).Decode(authData)
	if err != nil {
		err = errortypes.ParseError{
			errors.Wrap(
				err, "auth: Failed to parse auth response",
			),
		}
		return
	}

	tokn := &Token{
		Id:        state,
		Type:      Google,
		Secret:    secret,
		Timestamp: time.Now(),
		Provider:  provider.Id,
	}

	err = coll.Insert(tokn)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	redirect = authData.Url

	return
}
