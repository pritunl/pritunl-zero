package auth

import (
	"crypto/hmac"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"strings"
)

type StateProvider struct {
	Id    bson.ObjectId `json:"id"`
	Type  string        `json:"type"`
	Label string        `json:"label"`
}

type State struct {
	Providers []*StateProvider `json:"providers"`
}

func GetState() (state *State) {
	state = &State{
		Providers: []*StateProvider{},
	}

	for _, provider := range settings.Auth.Providers {
		provider := &StateProvider{
			Id:    provider.Id,
			Type:  provider.Type,
			Label: provider.Label,
		}
		state.Providers = append(state.Providers, provider)
	}

	return
}

func Local(db *database.Database, username, password string) (
	usr *user.User, errData *errortypes.ErrorData, err error) {

	usr, err = user.GetUsername(db, user.Local, username)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			err = nil
			errData = &errortypes.ErrorData{
				Error:   "auth_invalid",
				Message: "Authencation credentials are invalid",
			}
			break
		}
		return
	}

	valid := usr.CheckPassword(password)
	if !valid {
		errData = &errortypes.ErrorData{
			Error:   "auth_invalid",
			Message: "Authencation credentials are invalid",
		}
		return
	}

	return
}

func Request(c *gin.Context) {
	db := c.MustGet("db").(*database.Database)

	providerId := bson.ObjectIdHex(c.Query("id"))

	var provider *settings.Provider
	for _, prvidr := range settings.Auth.Providers {
		if prvidr.Id == providerId {
			provider = prvidr
			break
		}
	}

	if provider == nil {
		c.AbortWithStatus(404)
		return
	}

	loc := location.Get(c).String()

	switch provider.Type {
	case Google:
		redirect, err := GoogleRequest(db, loc, provider)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Redirect(302, redirect)
		return
	case OneLogin, Okta:
		body, err := SamlRequest(db, loc, provider)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Data(200, "text/html;charset=utf-8", body)
		return
	}

	c.AbortWithStatus(404)
}

func Callback(db *database.Database, sig, query string) (
	usr *user.User, errData *errortypes.ErrorData, err error) {

	params, err := url.ParseQuery(query)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "auth: Failed to parse query"),
		}
		return
	}

	state := params.Get("state")

	tokn, err := Get(db, state)
	if err != nil {
		return
	}

	hashFunc := hmac.New(sha512.New, []byte(tokn.Secret))
	hashFunc.Write([]byte(query))
	rawSignature := hashFunc.Sum(nil)
	testSig := base64.URLEncoding.EncodeToString(rawSignature)

	if subtle.ConstantTimeCompare([]byte(sig), []byte(testSig)) != 1 {
		errData = &errortypes.ErrorData{
			Error:   "authentication_error",
			Message: "Authentication error occurred",
		}
		return
	}

	provider := settings.Auth.GetProvider(tokn.Provider)
	if provider == nil {
		err = &errortypes.NotFoundError{
			errors.New("auth: Auth provider not found"),
		}
		return
	}

	err = tokn.Remove(db)
	if err != nil {
		return
	}

	roles := []string{}
	roles = append(roles, provider.DefaultRoles...)

	for _, role := range strings.Split(params.Get("roles"), ",") {
		if role != "" {
			roles = append(roles, role)
		}
	}

	username := params.Get("username")

	if provider.AutoCreate {
		usr = &user.User{
			Type:     provider.Type,
			Username: username,
			Roles:    roles,
		}

		errData, err = usr.Validate(db)
		if err != nil {
			return
		}

		if errData != nil {
			return
		}

		err = usr.Upsert(db)
		if err != nil {
			return
		}
	} else {
		usr, err = user.GetUsername(db, provider.Type, username)
		if err != nil {
			switch err.(type) {
			case *database.NotFoundError:
				err = nil
				errData = &errortypes.ErrorData{
					Error:   "unauthorized",
					Message: "Not authorized",
				}
				break
			}
			return
		}
	}

	return
}

func ValidateAdmin(db *database.Database, usr *user.User) (
	errData *errortypes.ErrorData, err error) {

	if usr.Disabled || usr.Administrator != "super" {
		errData = &errortypes.ErrorData{
			Error:   "unauthorized",
			Message: "Not authorized",
		}
		return
	}

	err = usr.SetActive(db)
	if err != nil {
		return
	}

	return
}

func Validate(db *database.Database, usr *user.User, srvc *service.Service,
	r *http.Request) (errData *errortypes.ErrorData, err error) {

	if usr.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "unauthorized",
			Message: "Not authorized",
		}
		return
	}

	usrRoles := set.NewSet()
	for _, role := range usr.Roles {
		usrRoles.Add(role)
	}

	roleMatch := false
	for _, role := range srvc.Roles {
		if usrRoles.Contains(role) {
			roleMatch = true
			break
		}
	}

	if !roleMatch {
		errData = &errortypes.ErrorData{
			Error:   "service_unauthorized",
			Message: "Not authorized for service",
		}
		return
	}

	polices, err := policy.GetService(db, srvc.Id)
	if err != nil {
		return
	}

	for _, polcy := range polices {
		errData, err = polcy.ValidateUser(db, usr, r)
		if err != nil || errData != nil {
			return
		}
	}

	polices, err = policy.GetRoles(db, usr.Roles)
	if err != nil {
		return
	}

	for _, polcy := range polices {
		errData, err = polcy.ValidateUser(db, usr, r)
		if err != nil || errData != nil {
			return
		}
	}

	err = usr.SetActive(db)
	if err != nil {
		return
	}

	return
}
