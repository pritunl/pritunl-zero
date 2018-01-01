package secondary

import (
	"fmt"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"time"
)

type SecondaryData struct {
	Token    string `json:"token"`
	Label    string `json:"label"`
	Push     bool   `json:"push"`
	Phone    bool   `json:"phone"`
	Passcode bool   `json:"passcode"`
	Sms      bool   `json:"sms"`
}

type Secondary struct {
	usr         *user.User                  `bson:"-"`
	provider    *settings.SecondaryProvider `bson:"-"`
	Id          string                      `bson:"_id"`
	ProviderId  bson.ObjectId               `bson:"provider_id"`
	UserId      bson.ObjectId               `bson:"user_id"`
	ChallengeId string                      `bson:"challenge_id"`
	Timestamp   time.Time                   `bson:"timestamp"`
	PushSent    bool                        `bson:"push_sent"`
	PhoneSent   bool                        `bson:"phone_sent"`
	SmsSent     bool                        `bson:"sms_sent"`
	Disabled    bool                        `bson:"disabled"`
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

func (s *Secondary) GetData() (data *SecondaryData, err error) {
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
		provider.Label,
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
	err = s.CommitFields(db, set.NewSet("disabled"))
	if err != nil {
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

	err = coll.Insert(s)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
