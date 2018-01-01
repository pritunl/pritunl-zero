package secondary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	oneloginClient = &http.Client{
		Timeout: 20 * time.Second,
	}
)

type oneloginAuthParams struct {
	GrantType string `json:"grant_type"`
}

type oneloginAuth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type oneloginUsersData struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
}

type oneloginUsers struct {
	Data []oneloginUsersData `json:"data"`
}

type oneloginOtpDevicesDataDevices struct {
	Id              int    `json:"id"`
	TypeDisplayName string `json:"type_display_name"`
	UserDisplayName string `json:"user_display_name"`
	AuthFactorName  string `json:"auth_factor_name"`
	Active          bool   `json:"boolean"`
	Default         bool   `json:"default"`
	NeedsTrigger    bool   `json:"needs_trigger"`
}

type oneloginOtpDevicesData struct {
	OtpDevices []oneloginOtpDevicesDataDevices `json:"otp_devices"`
}

type oneloginOtpDevices struct {
	Data oneloginOtpDevicesData `json:"data"`
}

type oneloginActivateParams struct {
	IpAddr string `json:"ipaddr"`
}

type oneloginActivateData struct {
	Id         int    `json:"id"`
	DeviceId   int    `json:"device_id"`
	StateToken string `json:"state_token"`
}

type oneloginActivate struct {
	Data []oneloginActivateData `json:"data"`
}

type oneloginVerifyParams struct {
	OtpToken   string `json:"otp_token,omitempty"`
	StateToken string `json:"state_token,omitempty"`
}

type oneloginVerifyStatus struct {
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type oneloginVerify struct {
	Status oneloginVerifyStatus `json:"status"`
}

func onelogin(db *database.Database, provider *settings.SecondaryProvider,
	r *http.Request, usr *user.User, factor, passcode string) (
	result bool, err error) {

	if factor != Push && factor != Passcode {
		err = &errortypes.UnknownError{
			errors.New("secondary: OneLogin invalid factor"),
		}
		return
	}

	if factor == Passcode && passcode == "" {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: OneLogin passcode empty"),
		}
		return
	}

	apiUrl := fmt.Sprintf(
		"https://api.%s.onelogin.com",
		provider.OneLoginRegion,
	)

	authParams := &oneloginAuthParams{
		GrantType: "client_credentials",
	}
	authBody, err := json.Marshal(authParams)
	if err != nil {
		err = &errortypes.ParseError{
			errors.New("secondary: OneLogin failed to parse auth params"),
		}
		return
	}

	reqUrl, _ := url.Parse(apiUrl + "/auth/oauth2/v2/token")
	req, err := http.NewRequest(
		"POST",
		reqUrl.String(),
		bytes.NewBuffer(authBody),
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secondary: OneLogin auth request failed"),
		}
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(
		"Authorization",
		fmt.Sprintf(
			"client_id:%s, client_secret:%s",
			provider.OneLoginId,
			provider.OneLoginSecret,
		),
	)

	resp, err := oneloginClient.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secondary: OneLogin auth request failed"),
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
		}).Info("secondary: OneLogin auth request bad status")

		err = &errortypes.RequestError{
			errors.New("secondary: OneLogin auth request bad status"),
		}
		return
	}

	auth := &oneloginAuth{}
	err = json.NewDecoder(resp.Body).Decode(auth)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "secondary: OneLogin auth parse failed"),
		}
		return
	}

	apiHeader := fmt.Sprintf(
		"bearer:%s",
		auth.AccessToken,
	)

	reqVals := url.Values{}
	reqVals.Set("username", usr.Username)
	reqVals.Set("fields", "id,username,email,status")
	reqUrl, _ = url.Parse(apiUrl + "/api/1/users")
	reqUrl.RawQuery = reqVals.Encode()

	req, err = http.NewRequest(
		"GET",
		reqUrl.String(),
		nil,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secondary: OneLogin users request failed"),
		}
		return
	}

	req.Header.Set("Authorization", apiHeader)

	resp, err = oneloginClient.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secondary: OneLogin users request failed"),
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
		}).Info("secondary: OneLogin users request bad status")

		err = &errortypes.RequestError{
			errors.New("secondary: OneLogin users request bad status"),
		}
		return
	}

	users := &oneloginUsers{}
	err = json.NewDecoder(resp.Body).Decode(users)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "secondary: OneLogin users parse failed"),
		}
		return
	}

	if users.Data == nil || len(users.Data) == 0 {
		reqVals := url.Values{}
		reqVals.Set("email", usr.Username)
		reqVals.Set("fields", "id,username,email,status")
		reqUrl, _ = url.Parse(apiUrl + "/api/1/users")
		reqUrl.RawQuery = reqVals.Encode()

		req, err = http.NewRequest(
			"GET",
			reqUrl.String(),
			nil,
		)
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(
					err, "secondary: OneLogin users request failed"),
			}
			return
		}

		req.Header.Set("Authorization", apiHeader)

		resp, err = oneloginClient.Do(req)
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(
					err, "secondary: OneLogin users request failed"),
			}
			return
		}
		defer resp.Body.Close()

		users = &oneloginUsers{}
		err = json.NewDecoder(resp.Body).Decode(users)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "secondary: OneLogin users parse failed"),
			}
			return
		}
	}

	shortUsername := ""
	if (users.Data == nil || len(users.Data) == 0) &&
		strings.Contains(usr.Username, "@") {

		shortUsername = strings.SplitN(usr.Username, "@", 2)[0]

		reqVals := url.Values{}
		reqVals.Set("username", shortUsername)
		reqVals.Set("fields", "id,username,email,status")
		reqUrl, _ = url.Parse(apiUrl + "/api/1/users")
		reqUrl.RawQuery = reqVals.Encode()

		req, err = http.NewRequest(
			"GET",
			reqUrl.String(),
			nil,
		)
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(
					err, "secondary: OneLogin users request failed"),
			}
			return
		}

		req.Header.Set("Authorization", apiHeader)

		resp, err = oneloginClient.Do(req)
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(
					err, "secondary: OneLogin users request failed"),
			}
			return
		}
		defer resp.Body.Close()

		users = &oneloginUsers{}
		err = json.NewDecoder(resp.Body).Decode(users)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "secondary: OneLogin users parse failed"),
			}
			return
		}
	}

	if users.Data == nil || len(users.Data) == 0 {
		err = &errortypes.NotFoundError{
			errors.New("secondary: OneLogin user not found"),
		}
		return
	}

	if users.Data[0].Id == 0 {
		err = &errortypes.NotFoundError{
			errors.New("secondary: OneLogin unknown user ID"),
		}
		return
	}

	if usr.Username != users.Data[0].Username &&
		usr.Username != users.Data[0].Email &&
		(shortUsername != "" && shortUsername != users.Data[0].Username) {

		err = &errortypes.AuthenticationError{
			errors.New("secondary: OneLogin username mismatch"),
		}
		return
	}

	if users.Data[0].Status != 1 {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: OneLogin user is not active"),
		}
		return
	}

	userId := users.Data[0].Id

	reqUrl, _ = url.Parse(apiUrl + fmt.Sprintf(
		"/api/1/users/%d/otp_devices",
		userId,
	))

	req, err = http.NewRequest(
		"GET",
		reqUrl.String(),
		nil,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secondary: OneLogin devices request failed"),
		}
		return
	}

	req.Header.Set("Authorization", apiHeader)

	resp, err = oneloginClient.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secondary: OneLogin devices request failed"),
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
		}).Info("secondary: OneLogin devices request bad status")

		err = &errortypes.RequestError{
			errors.New("secondary: OneLogin devices request bad status"),
		}
		return
	}

	devices := &oneloginOtpDevices{}
	err = json.NewDecoder(resp.Body).Decode(devices)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "secondary: OneLogin users parse failed"),
		}
		return
	}

	if devices.Data.OtpDevices == nil || len(devices.Data.OtpDevices) == 0 {
		err = &errortypes.NotFoundError{
			errors.New("secondary: OneLogin user has no devices"),
		}
		return
	}

	deviceId := 0
	needsTrigger := false
	for _, device := range devices.Data.OtpDevices {
		if device.AuthFactorName != "OneLogin Protect" {
			continue
		}

		if device.Default {
			deviceId = device.Id
			needsTrigger = device.NeedsTrigger
			break
		} else if deviceId == 0 {
			deviceId = device.Id
			needsTrigger = device.NeedsTrigger
		}

	}

	if deviceId == 0 {
		err = &errortypes.NotFoundError{
			errors.New("secondary: OneLogin user device type not found"),
		}
		return
	}

	stateToken := ""
	if needsTrigger || factor == Push {
		reqUrl, _ = url.Parse(apiUrl + fmt.Sprintf(
			"/api/1/users/%d/otp_devices/%d/trigger",
			userId,
			deviceId,
		))

		var activateBuffer *bytes.Buffer
		if factor == Push {
			activateParams := &oneloginActivateParams{
				IpAddr: node.Self.GetRemoteAddr(r),
			}
			activateBody, e := json.Marshal(activateParams)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(
						e,
						"secondary: OneLogin failed to parse activate params",
					),
				}
				return
			}

			activateBuffer = bytes.NewBuffer(activateBody)
		}

		req, err = http.NewRequest(
			"POST",
			reqUrl.String(),
			activateBuffer,
		)
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(
					err, "secondary: OneLogin activate request failed"),
			}
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", apiHeader)
		req.Header.Set("User-Agent", r.UserAgent())
		req.Header.Set("X-Forwarded-For", node.Self.GetRemoteAddr(r))

		resp, err = oneloginClient.Do(req)
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(
					err, "secondary: OneLogin activate request failed"),
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
			}).Info("secondary: OneLogin activate request bad status")

			err = &errortypes.RequestError{
				errors.New("secondary: OneLogin activate request bad status"),
			}
			return
		}

		activate := &oneloginActivate{}
		err = json.NewDecoder(resp.Body).Decode(activate)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "secondary: OneLogin activate parse failed"),
			}
			return
		}

		if activate.Data == nil || len(activate.Data) == 0 {
			err = &errortypes.UnknownError{
				errors.New("secondary: OneLogin activate empty data"),
			}
			return
		}

		if activate.Data[0].Id != userId {
			err = &errortypes.AuthenticationError{
				errors.New("secondary: OneLogin activate user id mismatch"),
			}
			return
		}

		if activate.Data[0].DeviceId != deviceId {
			err = &errortypes.AuthenticationError{
				errors.New("secondary: OneLogin activate device id mismatch"),
			}
			return
		}

		if activate.Data[0].StateToken == "" {
			err = &errortypes.AuthenticationError{
				errors.New("secondary: OneLogin activate state token empty"),
			}
			return
		}

		stateToken = activate.Data[0].StateToken
	}

	verifyParams := &oneloginVerifyParams{
		OtpToken:   passcode,
		StateToken: stateToken,
	}
	verifyBody, err := json.Marshal(verifyParams)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(
				err, "secondary: OneLogin failed to parse verify params"),
		}
		return
	}

	start := time.Now()
	for {
		if time.Now().Sub(start) > 45*time.Second {
			err = audit.New(
				db,
				r,
				usr.Id,
				audit.OneLoginDeny,
				audit.Fields{
					"one_login_factor": factor,
					"one_login_error":  "timeout",
				},
			)
			if err != nil {
				return
			}

			result = false

			return
		}

		reqUrl, _ = url.Parse(apiUrl + fmt.Sprintf(
			"/api/1/users/%d/otp_devices/%d/verify",
			userId,
			deviceId,
		))

		req, err = http.NewRequest(
			"POST",
			reqUrl.String(),
			bytes.NewBuffer(verifyBody),
		)
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(err, "secondary: OneLogin verify request failed"),
			}
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", apiHeader)
		req.Header.Set("User-Agent", r.UserAgent())
		req.Header.Set("X-Forwarded-For", node.Self.GetRemoteAddr(r))

		resp, err = oneloginClient.Do(req)
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(err, "secondary: OneLogin verify request failed"),
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 401 {
			body := ""
			data, _ := ioutil.ReadAll(resp.Body)
			if data != nil {
				body = string(data)
			}

			logrus.WithFields(logrus.Fields{
				"username":    usr.Username,
				"status_code": resp.StatusCode,
				"body":        body,
			}).Info("secondary: OneLogin verify request bad status")

			err = &errortypes.RequestError{
				errors.New("secondary: OneLogin verify request bad status"),
			}
			return
		}

		verify := &oneloginVerify{}
		err = json.NewDecoder(resp.Body).Decode(verify)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "secondary: OneLogin verify parse failed"),
			}
			return
		}

		if resp.StatusCode == 401 {
			if strings.Contains(
				verify.Status.Message, "Authentication pending") {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			err = audit.New(
				db,
				r,
				usr.Id,
				audit.OneLoginDeny,
				audit.Fields{
					"one_login_factor": factor,
				},
			)
			if err != nil {
				return
			}

			result = false

			return
		}

		if verify.Status.Type != "success" || verify.Status.Code != 200 &&
			verify.Status.Error {

			err = &errortypes.UnknownError{
				errors.New("secondary: OneLogin verify request bad data"),
			}
			return
		}

		err = audit.New(
			db,
			r,
			usr.Id,
			audit.OneLoginApprove,
			audit.Fields{
				"one_login_factor": factor,
			},
		)
		if err != nil {
			return
		}

		result = true
		return
	}

	return
}
