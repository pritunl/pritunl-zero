package auth

import (
	"bytes"
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	OneLogin = "onelogin"
	Okta     = "okta"
)

func SamlRequest(db *database.Database, location, query string,
	provider *settings.Provider) (body []byte, err error) {

	if provider.Type != OneLogin && provider.Type != Okta {
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
		License   string `json:"license"`
		Callback  string `json:"callback"`
		State     string `json:"state"`
		Secret    string `json:"secret"`
		SsoUrl    string `json:"sso_url"`
		IssuerUrl string `json:"issuer_url"`
		Cert      string `json:"cert"`
	}{
		License:   settings.System.License,
		Callback:  location + "/auth/callback",
		State:     state,
		Secret:    secret,
		SsoUrl:    provider.SamlUrl,
		IssuerUrl: provider.IssuerUrl,
		Cert:      provider.SamlCert,
	})

	req, err := http.NewRequest(
		"POST",
		settings.Auth.Server+"/v1/request/saml",
		bytes.NewBuffer(data),
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Auth request failed"),
		}
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Auth request failed"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = &errortypes.RequestError{
			errors.Wrapf(err, "auth: Auth server error %d", resp.StatusCode),
		}
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(
				err, "auth: Failed to parse auth response",
			),
		}
		return
	}

	tokn := &Token{
		Id:        state,
		Type:      provider.Type,
		Secret:    secret,
		Timestamp: time.Now(),
		Provider:  provider.Id,
		Query:     query,
	}

	err = coll.Insert(tokn)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
