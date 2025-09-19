package secondary

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/device"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/sirupsen/logrus"
)

type SecondaryData struct {
	Token          string `json:"token"`
	Label          string `json:"label"`
	Push           bool   `json:"push"`
	Phone          bool   `json:"phone"`
	Passcode       bool   `json:"passcode"`
	Sms            bool   `json:"sms"`
	Device         bool   `json:"device"`
	DeviceRegister bool   `json:"device_register"`
}

type Secondary struct {
	usr         *user.User                  `bson:"-"`
	provider    *settings.SecondaryProvider `bson:"-"`
	Id          string                      `bson:"_id"`
	ProviderId  bson.ObjectID               `bson:"provider_id,omitempty"`
	UserId      bson.ObjectID               `bson:"user_id"`
	Type        string                      `bson:"type"`
	ChallengeId string                      `bson:"challenge_id"`
	Timestamp   time.Time                   `bson:"timestamp"`
	PushSent    bool                        `bson:"push_sent"`
	PhoneSent   bool                        `bson:"phone_sent"`
	SmsSent     bool                        `bson:"sms_sent"`
	Disabled    bool                        `bson:"disabled"`
	WanSession  *webauthn.SessionData       `bson:"wan_session"`
}

func (s *Secondary) Push(db *database.Database, r *http.Request) (
	errData *errortypes.ErrorData, err error) {

	if s.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "secondary_disabled",
			Message: "Secondary authentication has already been completed",
		}
		return
	}

	if s.PushSent {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Push already sent"),
		}
		return
	}
	s.PushSent = true
	err = s.CommitFields(db, set.NewSet("push_sent"))
	if err != nil {
		return
	}

	provider, err := s.GetProvider()
	if err != nil {
		return
	}

	if !provider.PushFactor {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Push factor not available"),
		}
		return
	}

	usr, err := s.GetUser(db)
	if err != nil {
		return
	}

	result := false
	switch provider.Type {
	case Duo:
		result, err = duo(db, provider, r, usr, Push, "")
		if err != nil {
			return
		}
		break
	case OneLogin:
		result, err = onelogin(db, provider, r, usr, Push, "")
		if err != nil {
			return
		}
		break
	case Okta:
		result, err = okta(db, provider, r, usr, Push, "")
		if err != nil {
			return
		}
		break
	default:
		err = &errortypes.UnknownError{
			errors.New("secondary: Unknown secondary provider type"),
		}
		return
	}

	if !result {
		errData = &errortypes.ErrorData{
			Error:   "secondary_denied",
			Message: "Secondary authentication was denied",
		}
		return
	}

	return
}

func (s *Secondary) Phone(db *database.Database, r *http.Request) (
	errData *errortypes.ErrorData, err error) {

	if s.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "secondary_disabled",
			Message: "Secondary authentication has already been completed",
		}
		return
	}

	if s.PhoneSent {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Phone already sent"),
		}
		return
	}
	s.PhoneSent = true
	err = s.CommitFields(db, set.NewSet("phone_sent"))
	if err != nil {
		return
	}

	provider, err := s.GetProvider()
	if err != nil {
		return
	}

	if !provider.PhoneFactor {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Phone factor not available"),
		}
		return
	}

	usr, err := s.GetUser(db)
	if err != nil {
		return
	}

	result := false
	switch provider.Type {
	case Duo:
		result, err = duo(db, provider, r, usr, Phone, "")
		if err != nil {
			return
		}
		break
	default:
		err = &errortypes.UnknownError{
			errors.New("secondary: Unknown secondary provider type"),
		}
		return
	}

	if !result {
		errData = &errortypes.ErrorData{
			Error:   "secondary_denied",
			Message: "Secondary authentication was denied",
		}
		return
	}

	return
}

func (s *Secondary) Passcode(db *database.Database, r *http.Request,
	passcode string) (errData *errortypes.ErrorData, err error) {

	if s.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "secondary_disabled",
			Message: "Secondary authentication has already been completed",
		}
		return
	}

	provider, err := s.GetProvider()
	if err != nil {
		return
	}

	if !provider.PasscodeFactor {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Passcode factor not available"),
		}
		return
	}

	usr, err := s.GetUser(db)
	if err != nil {
		return
	}

	result := false
	switch provider.Type {
	case Duo:
		result, err = duo(db, provider, r, usr, Passcode, passcode)
		if err != nil {
			return
		}
		break
	case OneLogin:
		result, err = onelogin(db, provider, r, usr, Passcode, passcode)
		if err != nil {
			return
		}
		break
	case Okta:
		result, err = okta(db, provider, r, usr, Passcode, passcode)
		if err != nil {
			return
		}
		break
	default:
		err = &errortypes.UnknownError{
			errors.New("secondary: Unknown secondary provider type"),
		}
		return
	}

	if !result {
		errData = &errortypes.ErrorData{
			Error:   "secondary_denied",
			Message: "Secondary authentication was denied",
		}
		return
	}

	return
}

func (s *Secondary) Sms(db *database.Database, r *http.Request) (
	errData *errortypes.ErrorData, err error) {

	if s.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "secondary_disabled",
			Message: "Secondary authentication has already been completed",
		}
		return
	}

	if s.SmsSent {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Sms already sent"),
		}
		return
	}

	provider, err := s.GetProvider()
	if err != nil {
		return
	}

	if !provider.SmsFactor {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Sms factor not available"),
		}
		return
	}

	usr, err := s.GetUser(db)
	if err != nil {
		return
	}

	switch provider.Type {
	case Duo:
		_, err = duo(db, provider, r, usr, Sms, "")
		if err != nil {
			return
		}
		break
	default:
		err = &errortypes.UnknownError{
			errors.New("secondary: Unknown secondary provider type"),
		}
		return
	}

	s.SmsSent = true
	err = s.CommitFields(db, set.NewSet("sms_sent"))
	if err != nil {
		return
	}

	err = &IncompleteError{
		errors.New("secondary: Secondary auth is incomplete"),
	}

	return
}

func (s *Secondary) DeviceRegisterRequest(db *database.Database,
	origin string) (jsonResp interface{}, errData *errortypes.ErrorData,
	err error) {

	if s.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "secondary_disabled",
			Message: "Secondary registration has already been completed",
		}
		return
	}

	if s.ProviderId != DeviceProvider {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Device register not available"),
		}
		return
	}

	if s.WanSession != nil {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Device registration already requested"),
		}
		return
	}

	usr, err := s.GetUser(db)
	if err != nil {
		return
	}

	web, err := node.Self.GetWebauthn(origin, true)
	if err != nil {
		return
	}

	options, sessionData, err := web.BeginRegistration(usr)
	if err != nil {
		err = utils.ParseWebauthnError(err)
		return
	}

	s.WanSession = sessionData
	err = s.CommitFields(db, set.NewSet("wan_session"))
	if err != nil {
		return
	}

	jsonResp = options

	return
}

func (s *Secondary) DeviceRegisterResponse(db *database.Database,
	origin string, body io.Reader, name string) (
	devc *device.Device, errData *errortypes.ErrorData, err error) {

	if s.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "secondary_disabled",
			Message: "Secondary registration has already been completed",
		}
		return
	}

	if s.ProviderId != DeviceProvider {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Device register not available"),
		}
		return
	}

	if s.WanSession == nil {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Device registration not requested"),
		}
		return
	}

	usr, err := s.GetUser(db)
	if err != nil {
		return
	}

	data, err := protocol.ParseCredentialCreationResponseBody(body)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Webauthn parse error"),
		}
		return
	}

	web, err := node.Self.GetWebauthn(origin, true)
	if err != nil {
		return
	}

	credential, err := web.CreateCredential(usr, *s.WanSession, data)
	if err != nil {
		err = utils.ParseWebauthnError(err)
		return
	}

	devc = device.New(usr.Id, device.WebAuthn, device.Secondary)
	devc.User = usr.Id
	devc.Name = name
	devc.WanRpId = web.Config.RPID

	devc.MarshalWebauthn(credential)

	errData, err = devc.Validate(db)
	if err != nil || errData != nil {
		return
	}

	err = devc.Insert(db)
	if err != nil {
		return
	}

	return
}

func (s *Secondary) DeviceRequest(db *database.Database, origin string) (
	jsonResp interface{}, errData *errortypes.ErrorData, err error) {

	if s.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "secondary_disabled",
			Message: "Secondary authentication has already been completed",
		}
		return
	}

	if s.ProviderId != DeviceProvider {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Device sign not available"),
		}
		return
	}

	if s.WanSession != nil {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Device sign already requested"),
		}
		return
	}

	usr, err := s.GetUser(db)
	if err != nil {
		return
	}

	web, err := node.Self.GetWebauthn(origin, false)
	if err != nil {
		return
	}

	_, hasU2f, err := usr.LoadWebAuthnDevices(db)
	if err != nil {
		return
	}

	loginOpts := []webauthn.LoginOption{
		webauthn.WithUserVerification(protocol.VerificationPreferred),
	}
	if hasU2f {
		loginOpts = append(
			loginOpts,
			webauthn.WithAssertionExtensions(
				protocol.AuthenticationExtensions{
					"appid": settings.Local.AppId,
				},
			),
		)
	}

	options, sessionData, err := web.BeginLogin(usr, loginOpts...)
	if err != nil {
		err = utils.ParseWebauthnError(err)
		return
	}

	s.WanSession = sessionData
	err = s.CommitFields(db, set.NewSet("wan_session"))
	if err != nil {
		return
	}

	jsonResp = options

	return
}

func (s *Secondary) DeviceRespond(db *database.Database, origin string,
	body io.Reader) (errData *errortypes.ErrorData, err error) {

	if s.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "secondary_disabled",
			Message: "Secondary authentication has already been completed",
		}
		return
	}

	if s.ProviderId != DeviceProvider {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Device sign not available"),
		}
		return
	}

	if s.WanSession == nil {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Device sign not requested"),
		}
		return
	}

	usr, err := s.GetUser(db)
	if err != nil {
		return
	}

	data, err := protocol.ParseCredentialRequestResponseBody(body)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "handler: Webauthn parse error"),
		}
		return
	}

	web, err := node.Self.GetWebauthn(origin, false)
	if err != nil {
		return
	}

	devices, _, err := usr.LoadWebAuthnDevices(db)
	if err != nil {
		return
	}

	credential, err := web.ValidateLogin(
		usr, *s.WanSession, data)
	if err != nil {
		err = utils.ParseWebauthnError(err)
		logrus.WithFields(logrus.Fields{
			"user_id": s.UserId.Hex(),
			"error":   err,
		}).Error("secondary: Secondary authentication was denied")

		errData = &errortypes.ErrorData{
			Error:   "secondary_denied",
			Message: "Secondary authentication was denied",
		}
		return
	}

	for _, devc := range devices {
		if devc.Type == device.U2f {
			if !bytes.Equal(devc.U2fKeyHandle, credential.ID) {
				continue
			}
		} else if devc.Type == device.WebAuthn {
			if !bytes.Equal(devc.WanId, credential.ID) ||
				!bytes.Equal(devc.WanPublicKey, credential.PublicKey) {

				continue
			}
		} else {
			continue
		}

		devc.LastActive = time.Now()
		devc.MarshalWebauthn(credential)

		err = devc.CommitFields(db, set.NewSet(
			"last_active", "u2f_counter", "wan_authenticator"))
		if err != nil {
			return
		}

		return
	}

	errData = &errortypes.ErrorData{
		Error:   "secondary_denied",
		Message: "Secondary authentication was denied",
	}

	return
}

func (s *Secondary) DeviceRegisterSmartCard(db *database.Database,
	pubKey string, name string) (
	devc *device.Device, errData *errortypes.ErrorData, err error) {

	if s.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "secondary_disabled",
			Message: "Secondary registration has already been completed",
		}
		return
	}

	if s.ProviderId != DeviceProvider {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Device register not available"),
		}
		return
	}

	pubKey = strings.TrimSpace(pubKey)

	usr, err := s.GetUser(db)
	if err != nil {
		return
	}

	devcs, err := device.GetAllMode(db, usr.Id, device.Ssh)
	if err != nil {
		return
	}

	for _, dvc := range devcs {
		if dvc.SshPublicKey == pubKey {
			errData = &errortypes.ErrorData{
				Error:   "device_already_registered",
				Message: "Smart Card has already been registered",
			}
			return
		}
	}

	devc = device.New(usr.Id, device.SmartCard, device.Ssh)
	devc.User = usr.Id
	devc.Name = name
	devc.SshPublicKey = pubKey

	errData, err = devc.Validate(db)
	if err != nil || errData != nil {
		return
	}

	err = devc.Insert(db)
	if err != nil {
		return
	}

	return
}

func (s *Secondary) GetData() (data *SecondaryData, err error) {
	if s.ProviderId == DeviceProvider {
		label := ""
		register := false

		if strings.Contains(s.Type, "register") {
			label = "Register Device"
			register = true
		} else {
			label = "Device Authentication"
			register = false
		}

		data = &SecondaryData{
			Token:          s.Id,
			Label:          label,
			Push:           false,
			Phone:          false,
			Passcode:       false,
			Sms:            false,
			Device:         !register,
			DeviceRegister: register,
		}
		return
	}

	provider, err := s.GetProvider()
	if err != nil {
		return
	}

	data = &SecondaryData{
		Token:    s.Id,
		Label:    provider.Label,
		Push:     provider.PushFactor,
		Phone:    provider.PhoneFactor,
		Passcode: provider.PasscodeFactor || provider.SmsFactor,
		Sms:      provider.SmsFactor,
	}
	return
}

func (s *Secondary) GetQuery() (query string, err error) {
	if s.ProviderId == DeviceProvider {
		label := ""
		factor := ""

		if strings.Contains(s.Type, "register") {
			label = "Register Device"
			factor = "device_register"
		} else {
			label = "Device Authentication"
			factor = "device"
		}

		query = fmt.Sprintf(
			"secondary=%s&label=%s&factors=%s",
			s.Id,
			url.PathEscape(label),
			factor,
		)
		return
	}

	provider, err := s.GetProvider()
	if err != nil {
		return
	}

	factors := []string{}
	if provider.PushFactor {
		factors = append(factors, "push")
	}
	if provider.PhoneFactor {
		factors = append(factors, "phone")
	}
	if provider.PasscodeFactor || provider.SmsFactor {
		factors = append(factors, "passcode")
	}
	if provider.SmsFactor {
		factors = append(factors, "sms")
	}

	query = fmt.Sprintf(
		"secondary=%s&label=%s&factors=%s",
		s.Id,
		url.PathEscape(provider.Label),
		strings.Join(factors, ","),
	)

	return
}

func (s *Secondary) Complete(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if s.Disabled {
		errData = &errortypes.ErrorData{
			Error:   "secondary_disabled",
			Message: "Secondary authentication is already completed",
		}
		return
	}
	s.Disabled = true

	coll := db.SecondaryTokens()
	resp, err := coll.UpdateOne(db, &bson.M{
		"_id":      s.Id,
		"disabled": false,
	}, &bson.M{
		"$set": &bson.M{
			"disabled": true,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	if resp.ModifiedCount == 0 {
		errData = &errortypes.ErrorData{
			Error:   "secondary_update_disabled",
			Message: "Secondary authentication update is already completed",
		}
		return
	}

	return
}

func (s *Secondary) Handle(db *database.Database, r *http.Request,
	factor, passcode string) (errData *errortypes.ErrorData, err error) {

	switch factor {
	case Push:
		errData, err = s.Push(db, r)
		break
	case Phone:
		errData, err = s.Phone(db, r)
		break
	case Passcode:
		errData, err = s.Passcode(db, r, passcode)
		break
	case Sms:
		errData, err = s.Sms(db, r)
		break
	default:
		err = &errortypes.UnknownError{
			errors.New("secondary: Unknown secondary factor"),
		}
	}

	if err == nil && errData == nil && factor != Sms {
		errData, err = s.Complete(db)
		if err != nil || errData != nil {
			return
		}
	}

	return
}

func (s *Secondary) GetUser(db *database.Database) (
	usr *user.User, err error) {

	if s.usr != nil {
		usr = s.usr
		return
	}

	usr, err = user.Get(db, s.UserId)
	if err != nil {
		return
	}

	s.usr = usr

	return
}

func (s *Secondary) GetProvider() (provider *settings.SecondaryProvider,
	err error) {

	provider = settings.Auth.GetSecondaryProvider(s.ProviderId)
	if provider == nil {
		err = &errortypes.NotFoundError{
			errors.New("secondary: Secondary provider not found"),
		}
		return
	}

	return
}

func (s *Secondary) Commit(db *database.Database) (err error) {
	coll := db.SecondaryTokens()

	err = coll.Commit(s.Id, s)
	if err != nil {
		return
	}

	return
}

func (s *Secondary) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.SecondaryTokens()

	err = coll.CommitFields(s.Id, s, fields)
	if err != nil {
		return
	}

	return
}

func (s *Secondary) Insert(db *database.Database) (err error) {
	coll := db.SecondaryTokens()

	_, err = coll.InsertOne(db, s)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
