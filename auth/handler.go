package auth

import (
	"crypto/hmac"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type StateProvider struct {
	Id    interface{} `json:"id"`
	Type  string      `json:"type"`
	Label string      `json:"label"`
}

type StateProviders []*StateProvider

func (s StateProviders) Len() int {
	return len(s)
}

func (s StateProviders) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s StateProviders) Less(i, j int) bool {
	return s[i].Label < s[j].Label
}

type State struct {
	Providers StateProviders `json:"providers"`
}

func GetState() (state *State) {
	state = &State{
		Providers: StateProviders{},
	}

	google := false

	for _, provider := range settings.Auth.Providers {
		prv := &StateProvider{
			Type:  provider.Type,
			Label: provider.Label,
		}

		if provider.Type == Google {
			if google {
				continue
			}
			google = true
			prv.Id = Google
		} else {
			prv.Id = provider.Id
		}

		state.Providers = append(state.Providers, prv)
	}

	sort.Sort(state.Providers)

	return
}

func Local(db *database.Database, username, password string) (
	usr *user.User, errData *errortypes.ErrorData, err error) {

	usr, err = user.GetUsername(db, user.Local, username)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			usr = nil
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

	loc := utils.GetLocation(c.Request)

	id := c.Query("id")
	if id == Google {
		redirect, err := GoogleRequest(db, loc)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		c.Redirect(302, redirect)
		return
	} else {
		providerId := bson.ObjectIdHex(id)

		var provider *settings.Provider
		for _, prvidr := range settings.Auth.Providers {
			if prvidr.Id == providerId {
				provider = prvidr
				break
			}
		}

		if provider == nil {
			utils.AbortWithStatus(c, 404)
			return
		}

		switch provider.Type {
		case Azure:
			redirect, err := AzureRequest(db, loc, provider)
			if err != nil {
				utils.AbortWithError(c, 500, err)
				return
			}

			c.Redirect(302, redirect)
			return
		case OneLogin, Okta:
			body, err := SamlRequest(db, loc, provider)
			if err != nil {
				utils.AbortWithError(c, 500, err)
				return
			}

			c.Data(200, "text/html;charset=utf-8", body)
			return
		}
	}

	utils.AbortWithStatus(c, 404)
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
		switch err.(type) {
		case *database.NotFoundError:
			err = &InvalidState{
				errors.Wrap(err, "auth: Invalid state"),
			}
			break
		}
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

	username := params.Get("username")

	if username == "" {
		errData = &errortypes.ErrorData{
			Error:   "invalid_username",
			Message: "Invalid username",
		}
		return
	}

	var provider *settings.Provider
	if tokn.Type == Google {
		domainSpl := strings.SplitN(username, "@", 2)
		if len(domainSpl) == 2 {
			domain := domainSpl[1]
			if domain != "" {
				for _, prv := range settings.Auth.Providers {
					if prv.Type == Google && prv.Domain == domain {
						provider = prv
						break
					}
				}
			}
		}

		if provider == nil {
			errData = &errortypes.ErrorData{
				Error:   "unauthorized",
				Message: "Not authorized",
			}
			return
		}
	} else {
		provider = settings.Auth.GetProvider(tokn.Provider)
		if provider == nil {
			err = &errortypes.NotFoundError{
				errors.New("auth: Auth provider not found"),
			}
			return
		}
	}

	if provider.Type == Azure {
		usernameSpl := strings.SplitN(username, "/", 2)
		if len(usernameSpl) != 2 {
			errData = &errortypes.ErrorData{
				Error:   "invalid_username",
				Message: "Invalid username",
			}
			return
		}

		tenant := usernameSpl[0]
		username = usernameSpl[1]

		if tenant != provider.Tenant {
			errData = &errortypes.ErrorData{
				Error:   "invalid_tenant",
				Message: "Invalid tenant",
			}
			return
		}
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

	if provider.Type == Google {
		googleRoles, e := GoogleRoles(provider, username)
		if e != nil {
			err = e
			return
		}

		for _, role := range googleRoles {
			roles = append(roles, role)
		}
	}

	usr, err = user.GetUsername(db, provider.Type, username)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			usr = nil
			err = nil
			break
		default:
			return
		}
	}

	if usr == nil {
		if provider.AutoCreate {
			usr = &user.User{
				Type:     provider.Type,
				Username: username,
				Roles:    roles,
			}

			if errData != nil {
				return
			}

			err = usr.Upsert(db)
			if err != nil {
				return
			}

			event.PublishDispatch(db, "user.change")

			errData, err = usr.Validate(db)
			if err != nil {
				return
			}
		} else {
			errData = &errortypes.ErrorData{
				Error:   "unauthorized",
				Message: "Not authorized",
			}
			return
		}
	} else {
		switch provider.RoleManagement {
		case settings.Merge:
			changed := usr.RolesMerge(roles)
			if changed {
				errData, err = usr.Validate(db)
				if err != nil {
					return
				}

				if errData != nil {
					return
				}

				err = usr.CommitFields(db, set.NewSet("roles"))
				if err != nil {
					return
				}

				event.PublishDispatch(db, "user.change")
			}
			break
		case settings.Overwrite:
			changed := usr.RolesOverwrite(roles)
			if changed {
				errData, err = usr.Validate(db)
				if err != nil {
					return
				}

				if errData != nil {
					return
				}

				err = usr.CommitFields(db, set.NewSet("roles"))
				if err != nil {
					return
				}

				event.PublishDispatch(db, "user.change")
			}
			break
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

	return
}

func Validate(db *database.Database, usr *user.User,
	authr *authorizer.Authorizer, srvc *service.Service,
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

	if !authr.IsApi() {
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
	}

	return
}
