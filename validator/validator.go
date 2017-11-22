package validator

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/authorizer"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/user"
	"net/http"
)

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

func ValidateUser(db *database.Database, usr *user.User,
	authr *authorizer.Authorizer, r *http.Request) (
	errData *errortypes.ErrorData, err error) {

	if usr.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "unauthorized",
			Message: "Not authorized",
		}
		return
	}

	if !authr.IsApi() {
		polices, e := policy.GetRoles(db, usr.Roles)
		if e != nil {
			err = e
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

func ValidateProxy(db *database.Database, usr *user.User,
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
		polices, e := policy.GetService(db, srvc.Id)
		if e != nil {
			err = e
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
