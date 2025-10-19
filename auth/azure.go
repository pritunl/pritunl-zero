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
		Region      string `json:"region"`
		DirectoryId string `json:"directory_id"`
		AppId       string `json:"app_id"`
		AppSecret   string `json:"app_secret"`
	}{
		License:     settings.System.License,
		Callback:    location + "/auth/callback",
		State:       state,
		Secret:      secret,
		Region:      provider.Region,
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

	err = utils.CheckRequest(resp, "auth: Azure server error")
	if err != nil {
		return
	}

	athData := &authData{}
	err = json.NewDecoder(resp.Body).Decode(athData)
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

	redirect = athData.Url

	return
}

type azureTokenData struct {
	AccessToken string `json:"access_token"`
	Resource    string `json:"resource"`
	TokenType   string `json:"token_type"`
}

type azureMemberData struct {
	NextLink string           `json:"@odata.nextLink"`
	Value    []azureGroupData `json:"value"`
}

type azureUserData struct {
	Id                string `json:"id"`
	UserPrincipalName string `json:"userPrincipalName"`
	AccountEnabled    bool   `json:"accountEnabled"`
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

	err = utils.CheckRequest(resp, "auth: Azure server error")
	if err != nil {
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

func AzureRoles(provider *settings.Provider, username string) (
	roles []string, err error) {

	userId, active, err := AzureSync(provider, username)
	if err != nil {
		return
	}

	if !active {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Azure sync user disabled"),
		}
		return
	}

	if userId == "" {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Azure sync missing user ID"),
		}
		return
	}

	roles = []string{}

	token, err := azureGetToken(provider)
	if err != nil {
		return
	}

	reqUrlStr := fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/users/%s/memberOf",
		userId,
	)
	start := time.Now()

	for {
		reqUrl, e := url.Parse(reqUrlStr)
		if e != nil {
			err = &errortypes.ParseError{
				errors.Wrap(e, "auth: Failed to parse azure url"),
			}
			return
		}

		reqData, e := json.Marshal(struct {
			SecurityEnabledOnly string `json:"securityEnabledOnly"`
		}{
			SecurityEnabledOnly: "false",
		})
		if e != nil {
			err = e
			return
		}

		req, e := http.NewRequest(
			"GET",
			reqUrl.String(),
			bytes.NewBuffer(reqData),
		)
		if e != nil {
			err = &errortypes.RequestError{
				errors.Wrap(e, "auth: Failed to create azure request"),
			}
			return
		}

		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, e := client.Do(req)
		if e != nil {
			err = &errortypes.RequestError{
				errors.Wrap(e, "auth: Azure request failed"),
			}
			return
		}
		defer resp.Body.Close()

		err = utils.CheckRequest(resp, "auth: Azure server error")
		if err != nil {
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

		for _, groupData := range data.Value {
			groupName := groupData.DisplayName
			if groupName == "" {
				continue
			}

			roles = append(roles, groupName)
		}

		if data.NextLink != "" {
			reqUrlStr = data.NextLink
		} else {
			break
		}

		if time.Since(start) > 45*time.Second {
			err = &errortypes.RequestError{
				errors.New("auth: Azure group paging timeout"),
			}
			return
		}
	}

	return
}

func AzureSync(provider *settings.Provider, username string) (
	userId string, active bool, err error) {

	token, err := azureGetToken(provider)
	if err != nil {
		return
	}

	reqUrl, err := url.Parse(fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/%s/users/%s",
		provider.Tenant,
		url.QueryEscape(username),
	))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse azure url"),
		}
		return
	}

	query := reqUrl.Query()
	query.Set("$select", "id,userPrincipalName,accountEnabled")
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

	err = utils.CheckRequest(resp, "auth: Azure server error")
	if err != nil {
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

	if strings.ToLower(username) != strings.ToLower(
		data.UserPrincipalName) {

		err = &errortypes.ApiError{
			errors.Wrapf(
				err,
				"auth: Azure principal name '%s' does not match user '%s'",
				data.UserPrincipalName, username,
			),
		}
		return
	}

	userId = data.Id
	active = data.AccountEnabled

	return
}
