package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
)

const (
	Azure = "azure"
)

func AzureRequest(db *database.Database, location, query string,
	provider *settings.Provider) (redirect string, err error) {

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
		License     string `json:"license"`
		Callback    string `json:"callback"`
		State       string `json:"state"`
		Secret      string `json:"secret"`
		DirectoryId string `json:"directory_id"`
		AppId       string `json:"app_id"`
		AppSecret   string `json:"app_secret"`
	}{
		License:     settings.System.License,
		Callback:    location + "/auth/callback",
		State:       state,
		Secret:      secret,
		DirectoryId: provider.Tenant,
		AppId:       provider.ClientId,
		AppSecret:   provider.ClientSecret,
	})
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		"POST",
		settings.Auth.Server+"/v1/request/azure",
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
		Type:      Azure,
		Secret:    secret,
		Timestamp: time.Now(),
		Provider:  provider.Id,
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

type azureTokenData struct {
	AccessToken string `json:"access_token"`
	Resource    string `json:"resource"`
	TokenType   string `json:"token_type"`
}

type azureMemberData struct {
	Value []string `json:"value"`
}

type azureUserData struct {
	AccountEnabled bool `json:"accountEnabled"`
}

type azureGroupData struct {
	DisplayName string `json:"displayName"`
}

func azureGetToken(provider *settings.Provider) (token string, err error) {
	reqForm := url.Values{}
	reqForm.Add("grant_type", "client_credentials")
	reqForm.Add("client_id", provider.ClientId)
	reqForm.Add("client_secret", provider.ClientSecret)
	reqForm.Add("resource", "https://graph.microsoft.com")

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"https://login.microsoftonline.com/%s/oauth2/token",
			provider.Tenant,
		),
		strings.NewReader(reqForm.Encode()),
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Failed to create azure request"),
		}
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Azure request failed"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = &errortypes.RequestError{
			errors.Wrapf(err, "auth: Azure server error %d", resp.StatusCode),
		}
		return
	}

	tokenData := &azureTokenData{}
	err = json.NewDecoder(resp.Body).Decode(tokenData)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse response"),
		}
		return
	}

	token = tokenData.AccessToken

	return
}

func azureGetGroupName(provider *settings.Provider, token, groupId string) (
	name string, err error) {

	reqUrl, err := url.Parse(fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/%s/groups/%s",
		provider.Tenant,
		groupId,
	))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse azure url"),
		}
		return
	}

	query := reqUrl.Query()
	query.Set("$select", "displayName")
	reqUrl.RawQuery = query.Encode()

	req, err := http.NewRequest(
		"GET",
		reqUrl.String(),
		nil,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Failed to create azure request"),
		}
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Azure request failed"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = &errortypes.RequestError{
			errors.Wrapf(err, "auth: Azure server error %d", resp.StatusCode),
		}
		return
	}

	data := &azureGroupData{}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse response"),
		}
		return
	}

	name = data.DisplayName

	return
}

func AzureRoles(provider *settings.Provider, username string) (
	roles []string, err error) {

	roles = []string{}

	token, err := azureGetToken(provider)
	if err != nil {
		return
	}

	reqUrl, err := url.Parse(fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/%s/users/%s/getMemberGroups",
		provider.Tenant,
		username,
	))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse azure url"),
		}
		return
	}

	reqData, err := json.Marshal(struct {
		SecurityEnabledOnly string `json:"securityEnabledOnly"`
	}{
		SecurityEnabledOnly: "false",
	})
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		"POST",
		reqUrl.String(),
		bytes.NewBuffer(reqData),
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Failed to create azure request"),
		}
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Azure request failed"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = &errortypes.RequestError{
			errors.Wrapf(err, "auth: Azure server error %d", resp.StatusCode),
		}
		return
	}

	data := &azureMemberData{}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse response"),
		}
		return
	}

	for _, groupId := range data.Value {
		groupName, e := azureGetGroupName(provider, token, groupId)
		if e != nil {
			err = e
			return
		}

		roles = append(roles, groupName)
	}

	return
}

func AzureSync(db *database.Database, usr *user.User,
	provider *settings.Provider) (active bool, err error) {

	token, err := azureGetToken(provider)
	if err != nil {
		return
	}

	reqUrl, err := url.Parse(fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/%s/users/%s",
		provider.Tenant,
		usr.Username,
	))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse azure url"),
		}
		return
	}

	query := reqUrl.Query()
	query.Set("$select", "accountEnabled")
	reqUrl.RawQuery = query.Encode()

	req, err := http.NewRequest(
		"GET",
		reqUrl.String(),
		nil,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Failed to create azure request"),
		}
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Azure request failed"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = &errortypes.RequestError{
			errors.Wrapf(err, "auth: Azure server error %d", resp.StatusCode),
		}
		return
	}

	data := &azureUserData{}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse response"),
		}
		return
	}

	active = data.AccountEnabled

	return
}
