package twilio

import (
	"encoding/xml"

	"github.com/sirupsen/logrus"
	"github.com/twilio/twilio-go"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwimlSay struct {
	XMLName xml.Name `xml:"Say"`
	Voice   string   `xml:"voice,attr"`
	Loop    string   `xml:"loop,attr"`
	Message string   `xml:",chardata"`
}

type TwimlResponse struct {
	XMLName xml.Name  `xml:"Response"`
	Say     *TwimlSay `xml:"Say"`
}

func PhoneCall(number, message string) (err error) {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: settings.System.TwilioAccount,
		Password: settings.System.TwilioSecret,
	})

	params := &openapi.CreateCallParams{}
	params.SetFrom(settings.System.TwilioNumber)
	params.SetTo(number)

	twiml := &TwimlResponse{
		Say: &TwimlSay{
			Voice:   "alice",
			Loop:    "3",
			Message: FilterStrPhone(message, 160),
		},
	}

	twimlData, err := xml.Marshal(twiml)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "twilio: Failed to marshal twiml message"),
		}
		return
	}

	params.SetTwiml(string(twimlData))

	resp, err := client.Api.CreateCall(params)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "twilio: Twilio call error"),
		}
		return
	}

	respSid := *resp.Sid
	if respSid == "" {
		err = &errortypes.RequestError{
			errors.Wrap(err, "twilio: Invalid call sid"),
		}
		return
	}

	return
}

func TextMessage(number, message string) (err error) {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: settings.System.TwilioAccount,
		Password: settings.System.TwilioSecret,
	})

	params := &openapi.CreateMessageParams{}
	params.SetFrom(settings.System.TwilioNumber)
	params.SetTo(number)
	params.SetBody("Pritunl Alert: " + FilterStrMessage(message, 800))

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "twilio: Twilio message error"),
		}
		return
	}

	respSid := *resp.Sid
	if respSid == "" {
		err = &errortypes.RequestError{
			errors.Wrap(err, "twilio: Invalid message sid"),
		}
		return
	}

	if resp.ErrorCode != nil && resp.ErrorMessage != nil &&
		*resp.ErrorMessage != "" {

		logrus.WithFields(logrus.Fields{
			"number":        number,
			"message":       message,
			"source_number": settings.System.TwilioNumber,
			"error_code":    resp.ErrorCode,
			"error_message": resp.ErrorMessage,
		}).Error("twilio: Text message error")

		err = &errortypes.RequestError{
			errors.Wrap(err, "twilio: Twilio message error"),
		}
		return
	}

	return
}
