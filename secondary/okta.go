package secondary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
)

var (
	oktaClient = &http.Client{
		Timeout: 20 * time.Second,
	}
)

type oktaProfile struct {
	Email string `json:"email"`
	Login string `json:"login"`
}

type oktaUser struct {
	Id      string      `json:"id"`
	Status  string      `json:"status"`
	Profile oktaProfile `json:"profile"`
}

type oktaFactor struct {
	Id         string `json:"id"`
	FactorType string `json:"factorType"`
	Provider   string `json:"provider"`
	Status     string `json:"status"`
}

type oktaVerifyParams struct {
	Passcode string `json:"passCode,omitempty"`
}

type oktaLink struct {
	Href string `json:"href"`
}

type oktaLinks struct {
	Poll oktaLink `json:"poll"`
}

type oktaVerify struct {
	FactorResult string    `json:"factorResult"`
	Links        oktaLinks `json:"_links"`
}

func okta(db *database.Database, provider *settings.SecondaryProvider,
	r *http.Request, usr *user.User, factor, passcode string) (
	result bool, err error) {

	if factor != Push && factor != Passcode {
		err = &errortypes.UnknownError{
			errors.New("secondary: Okta invalid factor"),
		}
		return
	}

	if factor == Passcode && passcode == "" {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Okta passcode empty"),
		}
		return
	}

	apiUrl := fmt.Sprintf(
		"https://%s",
		provider.OktaDomain,
	)

	apiHeader := fmt.Sprintf(
		"SSWS %s",
		provider.OktaToken,
	)

	reqUrl, _ := url.Parse(apiUrl + "/api/v1/users/" + usr.Username)
	req, err := http.NewRequest(
		"GET",
		reqUrl.String(),
		nil,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secondary: Okta users request failed"),
		}
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", apiHeader)

	resp, err := oktaClient.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secondary: Okta users request failed"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body := ""
		data, _ := ioutil.ReadAll(resp.Body)
		if data != nil {
			body = string(data)
		}

		logrus.WithFields(logrus.Fields{
			"username":    usr.Username,
			"status_code": resp.StatusCode,
			"body":        body,
		}).Info("secondary: Okta users request bad status")

		err = &errortypes.RequestError{
			errors.New("secondary: Okta users request bad status"),
		}
		return
	}

	oktaUsr := &oktaUser{}
	err = json.NewDecoder(resp.Body).Decode(oktaUsr)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "secondary: Okta users parse failed"),
		}
		return
	}

	shortUsername := ""
	if oktaUsr.Id == "" && strings.Contains(usr.Username, "@") {
		shortUsername = strings.SplitN(usr.Username, "@", 2)[0]

		reqUrl, _ = url.Parse(apiUrl + "/api/v1/users/" + shortUsername)
		req, err = http.NewRequest(
			"GET",
			reqUrl.String(),
			nil,
		)
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(err, "secondary: Okta users request failed"),
			}
			return
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", apiHeader)

		resp, err = oktaClient.Do(req)
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(err, "secondary: Okta users request failed"),
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body := ""
			data, _ := ioutil.ReadAll(resp.Body)
			if data != nil {
				body = string(data)
			}

			logrus.WithFields(logrus.Fields{
				"username":    usr.Username,
				"status_code": resp.StatusCode,
				"body":        body,
			}).Info("secondary: Okta users request bad status")

			err = &errortypes.RequestError{
				errors.New("secondary: Okta users request bad status"),
			}
			return
		}

		oktaUsr = &oktaUser{}
		err = json.NewDecoder(resp.Body).Decode(oktaUsr)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "secondary: Okta users parse failed"),
			}
			return
		}
	}

	if oktaUsr.Id == "" {
		err = &errortypes.NotFoundError{
			errors.New("secondary: Okta users not found"),
		}
		return
	}

	if usr.Username != oktaUsr.Profile.Login &&
		usr.Username != oktaUsr.Profile.Email &&
		(shortUsername != "" && shortUsername != oktaUsr.Profile.Login) {

		err = &errortypes.AuthenticationError{
			errors.New("secondary: Okta username mismatch"),
		}
		return
	}

	if strings.ToLower(oktaUsr.Status) != "active" {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Okta user is not active"),
		}
		return
	}

	userId := oktaUsr.Id

	reqUrl, _ = url.Parse(apiUrl + fmt.Sprintf(
		"/api/v1/users/%s/factors", userId))
	req, err = http.NewRequest(
		"GET",
		reqUrl.String(),
		nil,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secondary: Okta factors request failed"),
		}
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", apiHeader)

	resp, err = oktaClient.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secondary: Okta factors request failed"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body := ""
		data, _ := ioutil.ReadAll(resp.Body)
		if data != nil {
			body = string(data)
		}

		logrus.WithFields(logrus.Fields{
			"username":    usr.Username,
			"status_code": resp.StatusCode,
			"body":        body,
		}).Info("secondary: Okta factors request bad status")

		err = &errortypes.RequestError{
			errors.New("secondary: Okta factors request bad status"),
		}
		return
	}

	factors := []*oktaFactor{}
	err = json.NewDecoder(resp.Body).Decode(&factors)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "secondary: Okta factors parse failed"),
		}
		return
	}

	if factors == nil || len(factors) == 0 {
		err = &errortypes.NotFoundError{
			errors.New("secondary: Okta user has no factors"),
		}
		return
	}

	factorId := ""
	for _, fctr := range factors {
		if fctr.Id == "" {
			continue
		}

		if strings.ToLower(fctr.Status) != "active" ||
			strings.ToLower(fctr.Provider) != "okta" {

			continue
		}

		switch factor {
		case Push:
			if strings.ToLower(fctr.FactorType) != "push" {
				continue
			}
			break
		case Passcode:
			if strings.ToLower(fctr.FactorType) != "token:software:totp" {
				continue
			}
			break
		default:
			continue
		}

		factorId = fctr.Id
	}

	verifyParams := &oktaVerifyParams{
		Passcode: passcode,
	}
	verifyBody, err := json.Marshal(verifyParams)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(
				err, "secondary: Okta failed to parse verify params"),
		}
		return
	}

	reqUrl, _ = url.Parse(apiUrl + fmt.Sprintf(
		"/api/v1/users/%s/factors/%s/verify", userId, factorId))
	req, err = http.NewRequest(
		"POST",
		reqUrl.String(),
		bytes.NewBuffer(verifyBody),
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secondary: Okta verify request failed"),
		}
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", apiHeader)
	req.Header.Set("User-Agent", r.UserAgent())
	req.Header.Set("X-Forwarded-For", node.Self.GetRemoteAddr(r))

	resp, err = oktaClient.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secondary: Okta verify request failed"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 &&
		resp.StatusCode != 403 {

		body := ""
		data, _ := ioutil.ReadAll(resp.Body)
		if data != nil {
			body = string(data)
		}

		logrus.WithFields(logrus.Fields{
			"username":    usr.Username,
			"status_code": resp.StatusCode,
			"body":        body,
		}).Info("secondary: Okta verify request bad status")

		err = &errortypes.RequestError{
			errors.New("secondary: Okta verify request bad status"),
		}
		return
	}

	verify := &oktaVerify{}
	err = json.NewDecoder(resp.Body).Decode(verify)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "secondary: Okta verify parse failed"),
		}
		return
	}

	if strings.ToLower(verify.FactorResult) == "waiting" &&
		verify.Links.Poll.Href != "" {

		start := time.Now()
		for {
			if time.Now().Sub(start) > 45*time.Second {
				err = audit.New(
					db,
					r,
					usr.Id,
					audit.OktaDeny,
					audit.Fields{
						"okta_factor": factor,
						"okta_error":  "timeout",
					},
				)
				if err != nil {
					return
				}

				result = false

				return
			}

			reqUrl, _ = url.Parse(verify.Links.Poll.Href)
			req, err = http.NewRequest(
				"GET",
				reqUrl.String(),
				nil,
			)
			if err != nil {
				err = &errortypes.RequestError{
					errors.Wrap(err, "secondary: Okta verify request failed"),
				}
				return
			}

			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", apiHeader)
			req.Header.Set("User-Agent", r.UserAgent())
			req.Header.Set("X-Forwarded-For", node.Self.GetRemoteAddr(r))

			resp, err = oktaClient.Do(req)
			if err != nil {
				err = &errortypes.RequestError{
					errors.Wrap(err, "secondary: Okta verify request failed"),
				}
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 && resp.StatusCode != 201 {
				body := ""
				data, _ := ioutil.ReadAll(resp.Body)
				if data != nil {
					body = string(data)
				}

				logrus.WithFields(logrus.Fields{
					"username":    usr.Username,
					"status_code": resp.StatusCode,
					"body":        body,
				}).Info("secondary: Okta verify request bad status")

				err = &errortypes.RequestError{
					errors.New("secondary: Okta verify request bad status"),
				}
				return
			}

			verify = &oktaVerify{}
			err = json.NewDecoder(resp.Body).Decode(verify)
			if err != nil {
				err = &errortypes.ParseError{
					errors.Wrap(err, "secondary: Okta verify parse failed"),
				}
				return
			}

			if strings.ToLower(verify.FactorResult) == "waiting" &&
				verify.Links.Poll.Href != "" {

				continue
			}

			break
		}
	}

	if strings.ToLower(verify.FactorResult) == "success" {
		err = audit.New(
			db,
			r,
			usr.Id,
			audit.OktaApprove,
			audit.Fields{
				"okta_factor": factor,
			},
		)
		if err != nil {
			return
		}

		result = true
	} else {
		err = audit.New(
			db,
			r,
			usr.Id,
			audit.OktaDeny,
			audit.Fields{
				"okta_factor": factor,
			},
		)
		if err != nil {
			return
		}

		result = false
	}

	return
}
