package validator

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/user"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func ValidateAdmin(db *database.Database, usr *user.User,
	isApi bool, r *http.Request) (deviceAuth bool, secProvider bson.ObjectId,
	errData *errortypes.ErrorData, err error) {

	if usr.Disabled || usr.Administrator != "super" {
		errData = &errortypes.ErrorData{
			Error:   "unauthorized",
			Message: "Not authorized",
		}
		return
	}

	if !isApi {
		policies, e := policy.GetRoles(db, usr.Roles)
		if e != nil {
			err = e
			return
		}

		for _, polcy := range policies {
			if polcy.AdminDeviceSecondary {
				deviceAuth = true
			}

			if polcy.AdminSecondary != "" && secProvider == "" {
				secProvider = polcy.AdminSecondary
			}
		}
	}

	return
}

func ValidateUser(db *database.Database, usr *user.User,
	isApi bool, r *http.Request) (deviceAuth bool, secProvider bson.ObjectId,
	errData *errortypes.ErrorData, err error) {

	if usr.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "unauthorized",
			Message: "Not authorized",
		}
		return
	}

	if !isApi {
		policies, e := policy.GetRoles(db, usr.Roles)
		if e != nil {
			err = e
			return
		}

		for _, polcy := range policies {
			errData, err = polcy.ValidateUser(db, usr, r)
			if err != nil || errData != nil {
				return
			}
		}

		for _, polcy := range policies {
			if polcy.UserDeviceSecondary {
				deviceAuth = true
			}

			if polcy.UserSecondary != "" && secProvider == "" {
				secProvider = polcy.UserSecondary
			}
		}
	}

	return
}

func ValidateProxy(db *database.Database, usr *user.User,
	isApi bool, srvc *service.Service, r *http.Request) (
	deviceAuth bool, secProvider bson.ObjectId,
	errData *errortypes.ErrorData, err error) {

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

	if !isApi {
		policies, e := policy.GetService(db, srvc.Id)
		if e != nil {
			err = e
			return
		}

		for _, polcy := range policies {
			errData, err = polcy.ValidateUser(db, usr, r)
			if err != nil || errData != nil {
				return
			}
		}

		for _, polcy := range policies {
			if polcy.ProxyDeviceSecondary {
				deviceAuth = true
			}

			if polcy.ProxySecondary != "" && secProvider == "" {
				secProvider = polcy.ProxySecondary
			}
		}

		policies, err = policy.GetRoles(db, usr.Roles)
		if err != nil {
			return
		}

		for _, polcy := range policies {
			errData, err = polcy.ValidateUser(db, usr, r)
			if err != nil || errData != nil {
				return
			}
		}

		for _, polcy := range policies {
			if polcy.ProxyDeviceSecondary {
				deviceAuth = true
			}

			if polcy.ProxySecondary != "" && secProvider == "" {
				secProvider = polcy.ProxySecondary
			}
		}
	}

	return
}
