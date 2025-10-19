package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
)

const (
	Google = "google"
)

func GoogleRequest(db *database.Database, location, query string) (
	redirect string, err error) {

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
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		"POST",
		settings.Auth.Server+"/v1/request/google",
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

	err = utils.CheckRequest(resp, "auth: Auth server error")
	if err != nil {
		return
	}

	authData := &authData{}
	err = json.NewDecoder(resp.Body).Decode(authData)
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
		Type:      Google,
		Secret:    secret,
		Timestamp: time.Now(),
		Query:     query,
	}

	_, err = coll.InsertOne(db, tokn)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	redirect = authData.Url

	return
}

func GoogleRoles(provider *settings.Provider, username string) (
	roles []string, err error) {

	roles = []string{}

	if provider.GoogleKey == "" && provider.GoogleEmail == "" {
		return
	}

	conf, err := google.JWTConfigFromJSON(
		[]byte(provider.GoogleKey),
		"https://www.googleapis.com/auth/admin.directory.group.readonly",
	)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(
				err, "auth: Failed to parse google key",
			),
		}
		return
	}

	conf.Subject = provider.GoogleEmail

	client := conf.Client(oauth2.NoContext)

	service, err := admin.New(client)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(
				err, "auth: Failed to parse google client",
			),
		}
		return
	}

	results, err := service.Groups.List().UserKey(username).Do()
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(
				err, "auth: Google api error getting user groups",
			),
		}
		return
	}

	for _, group := range results.Groups {
		roles = append(roles, group.Name)
	}

	return
}

func GoogleSync(db *database.Database, usr *user.User) (
	active bool, err error) {

	reqVals := url.Values{}
	reqVals.Set("user", usr.Username)
	reqVals.Set("license", settings.System.License)

	reqUrl, _ := url.Parse(settings.Auth.Server + "/update/google")
	reqUrl.RawQuery = reqVals.Encode()

	req, err := http.NewRequest(
		"GET",
		reqUrl.String(),
		nil,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Google request failed"),
		}
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Google request failed"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		active = true
	} else {
		err = &errortypes.RequestError{
			errors.Newf("auth: Google request bad status %d",
				resp.StatusCode),
		}
		return
	}

	return
}
