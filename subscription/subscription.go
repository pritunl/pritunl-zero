package subscription

import (
	"bytes"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/settings"
	"net/http"
	"time"
)

var (
	Sub    = &Subscription{}
	client = &http.Client{
		Timeout: 30 * time.Second,
	}
)

type Subscription struct {
	Active            bool      `json:"active"`
	Status            string    `json:"status"`
	Plan              string    `json:"plan"`
	Quantity          int       `json:"quantity"`
	Amount            int       `json:"amount"`
	PeriodEnd         time.Time `json:"period_end"`
	TrialEnd          time.Time `json:"trial_end"`
	CancelAtPeriodEnd bool      `json:"cancel_at_period_end"`
	Balance           int64     `json:"balance"`
	UrlKey            string    `json:"url_key"`
}

type subscriptionData struct {
	Active            bool   `json:"active"`
	Status            string `json:"status"`
	Plan              string `json:"plan"`
	Quantity          int    `json:"quantity"`
	Amount            int    `json:"amount"`
	PeriodEnd         int64  `json:"period_end"`
	TrialEnd          int64  `json:"trial_end"`
	CancelAtPeriodEnd bool   `json:"cancel_at_period_end"`
	Balance           int64  `json:"balance"`
	UrlKey            string `json:"url_key"`
}

func Update() (errData *errortypes.ErrorData, err error) {
	sub := &Subscription{}

	if settings.System.License == "" {
		Sub = sub
		return
	}

	data, err := json.Marshal(struct {
		Id      string `json:"id"`
		License string `json:"license"`
	}{
		Id:      settings.System.Name,
		License: settings.System.License,
	})

	req, err := http.NewRequest(
		"GET",
		"https://app.pritunl.com/subscription",
		bytes.NewBuffer(data),
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "subscription: Subscription request failed"),
		}
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "subscription: Subscription request failed"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errData = &errortypes.ErrorData{}
		err = json.NewDecoder(resp.Body).Decode(errData)
		if err != nil {
			errData = nil
		} else {
			logrus.WithFields(logrus.Fields{
				"error":     errData.Error,
				"error_msg": errData.Message,
			}).Error("subscription: Subscription error")
		}

		err = &errortypes.RequestError{
			errors.Wrap(err, "subscription: Subscription server error"),
		}
		return
	}

	subData := &subscriptionData{}
	err = json.NewDecoder(resp.Body).Decode(subData)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(
				err,
				"subscription: Failed to parse subscription response",
			),
		}
		return
	}

	sub.Active = subData.Active
	sub.Status = subData.Status
	sub.Plan = subData.Plan
	sub.Quantity = subData.Quantity
	sub.Amount = subData.Amount
	sub.CancelAtPeriodEnd = subData.CancelAtPeriodEnd
	sub.Balance = subData.Balance
	sub.UrlKey = subData.UrlKey

	if subData.PeriodEnd != 0 {
		sub.PeriodEnd = time.Unix(subData.PeriodEnd, 0)
	}
	if subData.TrialEnd != 0 {
		sub.TrialEnd = time.Unix(subData.TrialEnd, 0)
	}

	Sub = sub

	return
}

func update() {
	for {
		time.Sleep(30 * time.Minute)
		err, _ := Update()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("subscription: Update error")
			return
		}
	}
}

func init() {
	module := requires.New("subscription")
	module.After("settings")

	module.Handler = func() (err error) {
		Update()
		go update()
		return
	}
}
