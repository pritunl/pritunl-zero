package mhandlers

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/demo"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
)

type settingsData struct {
	AuthProviders             []*settings.Provider          `json:"auth_providers"`
	AuthSecondaryProviders    []*settings.SecondaryProvider `json:"auth_secondary_providers"`
	AuthAdminExpire           int                           `json:"auth_admin_expire"`
	AuthAdminMaxDuration      int                           `json:"auth_admin_max_duration"`
	AuthProxyExpire           int                           `json:"auth_proxy_expire"`
	AuthProxyMaxDuration      int                           `json:"auth_proxy_max_duration"`
	AuthUserExpire            int                           `json:"auth_user_expire"`
	AuthUserMaxDuration       int                           `json:"auth_user_max_duration"`
	AuthFastLogin             bool                          `json:"auth_fast_login"`
	AuthForceFastUserLogin    bool                          `json:"auth_force_fast_user_login"`
	AuthForceFastServiceLogin bool                          `json:"auth_force_fast_service_login"`
	TwilioAccount             string                        `json:"twilio_account"`
	TwilioSecret              string                        `json:"twilio_secret"`
	TwilioNumber              string                        `json:"twilio_number"`
	ElasticAddress            string                        `json:"elastic_address"`
	ElasticUsername           string                        `json:"elastic_username"`
	ElasticPassword           string                        `json:"elastic_password"`
	ElasticProxyRequests      bool                          `json:"elastic_proxy_requests"`
}

func getSettingsData() *settingsData {
	data := &settingsData{
		AuthProviders:             settings.Auth.Providers,
		AuthSecondaryProviders:    settings.Auth.SecondaryProviders,
		AuthAdminExpire:           settings.Auth.AdminExpire,
		AuthAdminMaxDuration:      settings.Auth.AdminMaxDuration,
		AuthProxyExpire:           settings.Auth.ProxyExpire,
		AuthProxyMaxDuration:      settings.Auth.ProxyMaxDuration,
		AuthUserExpire:            settings.Auth.UserExpire,
		AuthUserMaxDuration:       settings.Auth.UserMaxDuration,
		AuthFastLogin:             settings.Auth.FastLogin,
		AuthForceFastUserLogin:    settings.Auth.ForceFastUserLogin,
		AuthForceFastServiceLogin: settings.Auth.ForceFastServiceLogin,
		ElasticUsername:           settings.Elastic.Username,
		ElasticPassword:           settings.Elastic.Password,
		ElasticProxyRequests:      settings.Elastic.ProxyRequests,
		TwilioAccount:             settings.System.TwilioAccount,
		TwilioSecret:              settings.System.TwilioSecret,
		TwilioNumber:              settings.System.TwilioNumber,
	}

	if len(settings.Elastic.Addresses) != 0 {
		data.ElasticAddress = settings.Elastic.Addresses[0]
	}

	return data
}

func settingsGet(c *gin.Context) {
	data := getSettingsData()
	c.JSON(200, data)
}

func settingsPut(c *gin.Context) {
	if demo.Blocked(c) {
		return
	}

	db := c.MustGet("db").(*database.Database)
	data := &settingsData{}

	err := c.Bind(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Bind error"),
		}
		utils.AbortWithError(c, 500, err)
		return
	}

	fields := set.NewSet()

	elasticAddr := ""
	if len(settings.Elastic.Addresses) != 0 {
		elasticAddr = settings.Elastic.Addresses[0]
	}

	if elasticAddr != data.ElasticAddress {
		if data.ElasticAddress == "" {
			settings.Elastic.Addresses = []string{}
		} else {
			settings.Elastic.Addresses = []string{
				data.ElasticAddress,
			}
		}
		fields.Add("addresses")
	}

	if settings.Elastic.Username != data.ElasticUsername {
		settings.Elastic.Username = data.ElasticUsername
		fields.Add("username")
	}

	if settings.Elastic.Password != data.ElasticPassword {
		settings.Elastic.Password = data.ElasticPassword
		fields.Add("password")
	}

	if settings.Elastic.ProxyRequests != data.ElasticProxyRequests {
		settings.Elastic.ProxyRequests = data.ElasticProxyRequests
		fields.Add("proxy_requests")
	}

	if fields.Len() != 0 {
		err = settings.Commit(db, settings.Elastic, fields)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}
	}

	fields = set.NewSet()

	if settings.System.TwilioAccount != data.TwilioAccount {
		settings.System.TwilioAccount = data.TwilioAccount
		fields.Add("twilio_account")
	}

	if settings.System.TwilioSecret != data.TwilioSecret {
		settings.System.TwilioSecret = data.TwilioSecret
		fields.Add("twilio_secret")
	}

	if settings.System.TwilioNumber != data.TwilioNumber {
		settings.System.TwilioNumber = data.TwilioNumber
		fields.Add("twilio_number")
	}

	if fields.Len() != 0 {
		err = settings.Commit(db, settings.System, fields)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}
	}

	fields = set.NewSet(
		"providers",
		"secondary_providers",
	)

	if settings.Auth.AdminExpire != data.AuthAdminExpire {
		settings.Auth.AdminExpire = data.AuthAdminExpire
		fields.Add("admin_expire")
	}
	if settings.Auth.AdminMaxDuration != data.AuthAdminMaxDuration {
		settings.Auth.AdminMaxDuration = data.AuthAdminMaxDuration
		fields.Add("admin_max_duration")
	}
	if settings.Auth.ProxyExpire != data.AuthProxyExpire {
		settings.Auth.ProxyExpire = data.AuthProxyExpire
		fields.Add("proxy_expire")
	}
	if settings.Auth.ProxyMaxDuration != data.AuthProxyMaxDuration {
		settings.Auth.ProxyMaxDuration = data.AuthProxyMaxDuration
		fields.Add("proxy_max_duration")
	}
	if settings.Auth.UserExpire != data.AuthUserExpire {
		settings.Auth.UserExpire = data.AuthUserExpire
		fields.Add("user_expire")
	}
	if settings.Auth.UserMaxDuration != data.AuthUserMaxDuration {
		settings.Auth.UserMaxDuration = data.AuthUserMaxDuration
		fields.Add("user_max_duration")
	}
	if settings.Auth.FastLogin != data.AuthFastLogin {
		settings.Auth.FastLogin = data.AuthFastLogin
		fields.Add("fast_login")
	}
	if settings.Auth.ForceFastUserLogin != data.AuthForceFastUserLogin {
		settings.Auth.ForceFastUserLogin = data.AuthForceFastUserLogin
		fields.Add("force_fast_user_login")
	}
	if settings.Auth.ForceFastServiceLogin != data.AuthForceFastServiceLogin {
		settings.Auth.ForceFastServiceLogin = data.AuthForceFastServiceLogin
		fields.Add("force_fast_service_login")
	}

	for _, provider := range data.AuthProviders {
		errData, err := provider.Validate(db)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		if errData != nil {
			c.JSON(400, errData)
			return
		}
	}
	settings.Auth.Providers = data.AuthProviders

	for _, provider := range data.AuthSecondaryProviders {
		errData, err := provider.Validate(db)
		if err != nil {
			utils.AbortWithError(c, 500, err)
			return
		}

		if errData != nil {
			c.JSON(400, errData)
			return
		}
	}
	settings.Auth.SecondaryProviders = data.AuthSecondaryProviders

	err = settings.Commit(db, settings.Auth, fields)
	if err != nil {
		utils.AbortWithError(c, 500, err)
		return
	}

	_ = event.PublishDispatch(db, "settings.change")

	data = getSettingsData()
	c.JSON(200, data)
}
