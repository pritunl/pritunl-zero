package secondary

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
)

type duoApiResp struct {
	Result    string `json:"result"`
	Status    string `json:"status"`
	StatusMsg string `json:"status_msg"`
}

type duoApi struct {
	Stat     string     `json:"stat"`
	Code     int        `json:"code"`
	Message  string     `json:"message"`
	Response duoApiResp `json:"response"`
}

func duo(db *database.Database, provider *settings.SecondaryProvider,
	r *http.Request, usr *user.User, factor, passcode string) (
	result bool, err error) {

	if factor == Passcode && passcode == "" {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Duo passcode empty"),
		}
		return
	}

	api := duoapi.NewDuoApi(
		provider.DuoKey,
		provider.DuoSecret,
		provider.DuoHostname,
		"pritunl-zero",
	)

	query := url.Values{}
	query.Set("username", usr.Username)
	query.Set("ipaddr", node.Self.GetRemoteAddr(r))

	switch factor {
	case Push:
		query.Set("factor", "push")
		query.Set("device", "auto")
		break
	case Phone:
		query.Set("factor", "phone")
		query.Set("device", "auto")
		break
	case Passcode:
		query.Set("factor", "passcode")
		query.Set("passcode", passcode)
		break
	case Sms:
		query.Set("factor", "sms")
		query.Set("device", "auto")
		break
	}

	resp, data, err := api.SignedCall(
		"POST",
		"/auth/v2/auth",
		query,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secondary: Duo auth request failed"),
		}
		return
	}

	if data == nil {
		err = &errortypes.RequestError{
			errors.Newf(
				"secondary: Duo auth request failed %d",
				resp.StatusCode,
			),
		}
		return
	}

	duoData := &duoApi{}
	err = json.Unmarshal(data, duoData)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrapf(
				err,
				"secondary: Failed to parse Duo response %d",
				resp.StatusCode,
			),
		}
		return
	}

	if resp.StatusCode != 200 {
		logrus.WithFields(logrus.Fields{
			"username":       usr.Username,
			"status_code":    resp.StatusCode,
			"duo_factor":     factor,
			"duo_stat":       duoData.Stat,
			"duo_code":       duoData.Code,
			"duo_msg":        duoData.Message,
			"duo_result":     duoData.Response.Result,
			"duo_status":     duoData.Response.Status,
			"duo_status_msg": duoData.Response.StatusMsg,
		}).Error("secondary: Duo auth request failed")

		err = &errortypes.RequestError{
			errors.New("secondary: Duo auth request failed"),
		}
	}

	switch duoData.Response.Result {
	case "allow":
		err = audit.New(
			db,
			r,
			usr.Id,
			audit.DuoApprove,
			audit.Fields{
				"duo_factor": factor,
			},
		)
		if err != nil {
			return
		}

		result = true

		break
	case "deny":
		if factor != Sms {
			err = audit.New(
				db,
				r,
				usr.Id,
				audit.DuoDeny,
				audit.Fields{
					"duo_factor":     factor,
					"duo_status":     duoData.Response.Status,
					"duo_status_msg": duoData.Response.StatusMsg,
				},
			)
			if err != nil {
				return
			}
		}

		break
	default:
		logrus.WithFields(logrus.Fields{
			"username":       usr.Username,
			"status_code":    resp.StatusCode,
			"duo_factor":     factor,
			"duo_stat":       duoData.Stat,
			"duo_code":       duoData.Code,
			"duo_msg":        duoData.Message,
			"duo_result":     duoData.Response.Result,
			"duo_status":     duoData.Response.Status,
			"duo_status_msg": duoData.Response.StatusMsg,
		}).Error("secondary: Duo auth request unknown")

		err = &errortypes.RequestError{
			errors.New("secondary: Duo auth request unknown"),
		}
		return
	}

	return
}
