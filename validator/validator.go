package validator

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/user"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

func ValidateAdmin(db *database.Database, usr *user.User,
	isApi bool, r *http.Request) (deviceAuth bool, secProvider bson.ObjectId,
	errAudit audit.Fields, errData *errortypes.ErrorData, err error) {

	if !usr.ActiveUntil.IsZero() && usr.ActiveUntil.Before(time.Now()) {
		usr.ActiveUntil = time.Time{}
		usr.Disabled = true
		err = usr.CommitFields(db, set.NewSet("active_until", "disabled"))
		if err != nil {
			return
		}

		event.PublishDispatch(db, "user.change")

		errAudit = audit.Fields{
			"error":   "user_disabled",
			"message": "User is disabled from expired active time",
		}
		errData = &errortypes.ErrorData{
			Error:   "unauthorized",
			Message: "Not authorized",
		}
		return
	}

	if usr.Disabled {
		errAudit = audit.Fields{
			"error":   "user_disabled",
			"message": "User is disabled",
		}
		errData = &errortypes.ErrorData{
			Error:   "unauthorized",
			Message: "Not authorized",
		}
		return
	}

	if usr.Administrator != "super" {
		errAudit = audit.Fields{
			"error":   "user_not_super",
			"message": "User is not super user",
		}
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
	errAudit audit.Fields, errData *errortypes.ErrorData, err error) {

	if !usr.ActiveUntil.IsZero() && usr.ActiveUntil.Before(time.Now()) {
		usr.ActiveUntil = time.Time{}
		usr.Disabled = true
		err = usr.CommitFields(db, set.NewSet("active_until", "disabled"))
		if err != nil {
			return
		}

		event.PublishDispatch(db, "user.change")

		errAudit = audit.Fields{
			"error":   "user_disabled",
			"message": "User is disabled from expired active time",
		}
		errData = &errortypes.ErrorData{
			Error:   "unauthorized",
			Message: "Not authorized",
		}
		return
	}

	if usr.Disabled {
		errAudit = audit.Fields{
			"error":   "user_disabled",
			"message": "User is disabled",
		}
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
	errAudit audit.Fields, errData *errortypes.ErrorData, err error) {

	if !usr.ActiveUntil.IsZero() && usr.ActiveUntil.Before(time.Now()) {
		usr.ActiveUntil = time.Time{}
		usr.Disabled = true
		err = usr.CommitFields(db, set.NewSet("active_until", "disabled"))
		if err != nil {
			return
		}

		event.PublishDispatch(db, "user.change")

		errAudit = audit.Fields{
			"error":   "user_disabled",
			"message": "User is disabled from expired active time",
		}
		errData = &errortypes.ErrorData{
			Error:   "unauthorized",
			Message: "Not authorized",
		}
		return
	}

	if usr.Disabled {
		errAudit = audit.Fields{
			"error":   "user_disabled",
			"message": "User is disabled",
		}
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
		errAudit = audit.Fields{
			"error":   "service_unauthorized",
			"message": "User does not have roles required to access service",
		}
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
