package alertevent

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/device"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/twilio"
)

var (
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
)

type AlertParams struct {
	License string `json:"license"`
	Number  string `json:"number"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func SendTest(db *database.Database, devc *device.Device) (
	errData *errortypes.ErrorData, err error) {

	err = devc.SetActive(db)
	if err != nil {
		return
	}

	errData, err = Send(devc.Number, "Test alert message", devc.Type)
	if err != nil {
		return
	}

	return
}

func Send(number, message, alertType string) (
	errData *errortypes.ErrorData, err error) {

	if settings.System.TwilioAccount != "" {
		if alertType == device.Call {
			err = twilio.PhoneCall(number, message)
			if err != nil {
				return
			}
		} else if alertType == device.Message {
			err = twilio.TextMessage(number, message)
			if err != nil {
				return
			}
		} else {
			err = &errortypes.ParseError{
				errors.Wrap(
					err, "alert: Unknown alert type"),
			}
			return
		}
	} else {
		params := &AlertParams{
			License: settings.System.License,
			Number:  number,
			Type:    alertType,
			Message: message,
		}

		alertBody, err := json.Marshal(params)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(
					err, "alert: Failed to parse alert params"),
			}
			return
		}

		req, err := http.NewRequest(
			"POST",
			"https://app.pritunl.com/alert",
			bytes.NewBuffer(alertBody),
		)
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(err, "alert: Failed to create alert request"),
			}
			return
		}

		req.Header.Set("User-Agent", "pritunl-cloud")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(err, "alert: Alert request failed"),
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

			errData = &errortypes.ErrorData{}
			err = json.Unmarshal(data, errData)
			if err != nil || errData.Error == "" {
				errData = nil
			}

			err = &errortypes.RequestError{
				errors.Newf(
					"alert: Alert server error %d - %s",
					resp.StatusCode, body),
			}

			return
		}
	}

	return
}
