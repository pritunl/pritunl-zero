package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"net/http"
	"net/url"
	"time"
)

const (
	AuthZero = "authzero"
)

func AuthZeroRequest(db *database.Database, location, query string,
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
		License   string `json:"license"`
		Callback  string `json:"callback"`
		State     string `json:"state"`
		Secret    string `json:"secret"`
		AppDomain string `json:"app_domain"`
		AppId     string `json:"app_id"`
		AppSecret string `json:"app_secret"`
	}{
		License:   settings.System.License,
		Callback:  location + "/auth/callback",
		State:     state,
		Secret:    secret,
		AppDomain: provider.Domain,
		AppId:     provider.ClientId,
		AppSecret: provider.ClientSecret,
	})
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		"POST",
		settings.Auth.Server+"/v1/request/authzero",
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
		Type:      AuthZero,
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

	redirect = authData.Url

	return
}

type authZeroJwks struct {
	Keys []json.RawMessage `json:"keys"`
}

type authZeroTokenReq struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
}

type authZeroTokenData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type authZeroToken struct {
	Iss   string `json:"iss"`
	Sub   string `json:"sub"`
	Aud   string `json:"aud"`
	Exp   int    `json:"exp"`
	Iat   int    `json:"iat"`
	Scope string `json:"scope"`
}

type authZeroAppAuthorization struct {
	Groups []string `json:"groups"`
	Roles  []string `json:"roles"`
}

type authZeroAppMetadata struct {
	Authorization authZeroAppAuthorization `json:"authorization"`
}

type authZeroUser struct {
	UserId      string              `json:"user_id"`
	Email       string              `json:"email"`
	AppMetadata authZeroAppMetadata `json:"app_metadata"`
}

//func authZeroGetJwk(provider *settings.Provider) (
//	jwk *jose.JSONWebKey, err error) {
//
//	req, err := http.NewRequest(
//		"GET",
//		fmt.Sprintf(
//			"https://%s.auth0.com/.well-known/jwks.json",
//			provider.Domain,
//		),
//		nil,
//	)
//	if err != nil {
//		err = &errortypes.RequestError{
//			errors.Wrap(err, "auth: Failed to create auth0 request"),
//		}
//		return
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		err = &errortypes.RequestError{
//			errors.Wrap(err, "auth: auth0 request failed"),
//		}
//		return
//	}
//	defer resp.Body.Close()
//
//	data := &authZeroJwks{}
//	err = json.NewDecoder(resp.Body).Decode(data)
//	if err != nil {
//		err = &errortypes.ParseError{
//			errors.Wrap(err, "auth: Failed to parse response"),
//		}
//		return
//	}
//
//	if len(data.Keys) < 1 {
//		err = &errortypes.ParseError{
//			errors.Wrap(err, "auth: No JWK keys available"),
//		}
//		return
//	}
//
//	jwk = &jose.JSONWebKey{}
//
//	err = jwk.UnmarshalJSON(data.Keys[0])
//	if err != nil {
//		err = &errortypes.ParseError{
//			errors.Wrap(err, "auth: Failed to parse jwt key"),
//		}
//		return
//	}
//
//	return
//}
//
//func authZeroGetJwkToken(provider *settings.Provider) (
//	accessToken string, token *authZeroToken, err error) {
//
//	reqData := &authZeroTokenReq{
//		GrantType:    "client_credentials",
//		ClientId:     provider.ClientId,
//		ClientSecret: provider.ClientSecret,
//		Audience: fmt.Sprintf(
//			"https://%s.auth0.com/api/v2/", provider.Domain),
//	}
//
//	reqDataBuf := &bytes.Buffer{}
//	err = json.NewEncoder(reqDataBuf).Encode(reqData)
//	if err != nil {
//		err = &errortypes.ParseError{
//			errors.Wrap(err, "auth: Failed to parse request data"),
//		}
//		return
//	}
//
//	req, err := http.NewRequest(
//		"POST",
//		fmt.Sprintf("https://%s.auth0.com/oauth/token", provider.Domain),
//		reqDataBuf,
//	)
//	if err != nil {
//		err = &errortypes.RequestError{
//			errors.Wrap(err, "auth: Failed to create auth0 request"),
//		}
//		return
//	}
//
//	req.Header.Add("Content-Type", "application/json")
//
//	resp, err := client.Do(req)
//	if err != nil {
//		err = &errortypes.RequestError{
//			errors.Wrap(err, "auth: auth0 request failed"),
//		}
//		return
//	}
//	defer resp.Body.Close()
//
//	tokenData := &authZeroTokenData{}
//	err = json.NewDecoder(resp.Body).Decode(tokenData)
//	if err != nil {
//		err = &errortypes.ParseError{
//			errors.Wrap(err, "auth: Failed to parse response"),
//		}
//		return
//	}
//
//	accessToken = tokenData.AccessToken
//
//	object, err := jose.ParseSigned(tokenData.AccessToken)
//	if err != nil {
//		err = &errortypes.ParseError{
//			errors.Wrap(err, "auth: Failed to parse jwt data"),
//		}
//		return
//	}
//
//	jwt, err := authZeroGetJwk(provider)
//	if err != nil {
//		return
//	}
//
//	data, err := object.Verify(jwt)
//	if err != nil {
//		err = &errortypes.ParseError{
//			errors.Wrap(err, "auth: Failed to verify jwt data"),
//		}
//		return
//	}
//
//	token = &authZeroToken{}
//	err = json.Unmarshal(data, token)
//	if err != nil {
//		err = &errortypes.ParseError{
//			errors.Wrap(err, "auth: Failed to parse jwt token"),
//		}
//		return
//	}
//
//	return
//}

func authZeroGetToken(provider *settings.Provider) (token string, err error) {
	reqData := &authZeroTokenReq{
		GrantType:    "client_credentials",
		ClientId:     provider.ClientId,
		ClientSecret: provider.ClientSecret,
		Audience: fmt.Sprintf(
			"https://%s.auth0.com/api/v2/", provider.Domain),
	}

	reqDataBuf := &bytes.Buffer{}
	err = json.NewEncoder(reqDataBuf).Encode(reqData)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse request data"),
		}
		return
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://%s.auth0.com/oauth/token", provider.Domain),
		reqDataBuf,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Failed to create auth0 request"),
		}
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: auth0 request failed"),
		}
		return
	}
	defer resp.Body.Close()

	tokenData := &authZeroTokenData{}
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

func AuthZeroRoles(provider *settings.Provider, username string) (
	roles []string, err error) {

	roles = []string{}

	token, err := authZeroGetToken(provider)
	if err != nil {
		return
	}

	reqUrl, err := url.Parse(fmt.Sprintf(
		"https://%s.auth0.com/api/v2/users",
		provider.Domain,
	))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse auth0 url"),
		}
		return
	}

	query := reqUrl.Query()
	query.Set("search_engine", "v3")
	query.Set("email", username)
	reqUrl.RawQuery = query.Encode()

	req, err := http.NewRequest(
		"GET",
		reqUrl.String(),
		nil,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: Failed to create auth0 request"),
		}
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "auth: auth0 request failed"),
		}
		return
	}
	defer resp.Body.Close()

	data := []*authZeroUser{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse response"),
		}
		return
	}

	userId := ""

	for _, usr := range data {
		if usr.Email != username {
			continue
		}

		userId = usr.UserId

		if usr.AppMetadata.Authorization.Roles != nil {
			roles = usr.AppMetadata.Authorization.Roles
		}

		break
	}

	if userId == "" {
		err = &errortypes.NotFoundError{
			errors.Wrap(err, "auth: Failed to find auth0 user"),
		}
		return
	}

	return
}
