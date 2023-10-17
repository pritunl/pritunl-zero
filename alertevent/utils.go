package alertevent

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/device"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/twilio"
)

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

	return
}
