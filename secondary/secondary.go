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
	"strings"
)

type SecondaryData struct {
	Token    string `json:"token"`
	Push     bool   `json:"push"`
	Phone    bool   `json:"phone"`
	Passcode bool   `json:"passcode"`
	Sms      bool   `json:"sms"`
}

type Secondary struct {
	usr        *user.User                  `bson:"-"`
	provider   *settings.SecondaryProvider `bson:"-"`
	Id         string                      `bson:"_id"`
	ProviderId bson.ObjectId               `bson:"provider_id"`
	UserId     bson.ObjectId               `bson:"user_id"`
	SmsSent    bool                        `bson:"sms_sent"`
}

func (s *Secondary) Push(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

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

	// TODO

	return
}

func (s *Secondary) Phone(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	provider, err := s.GetProvider()
	if err != nil {
		return
	}

	if !provider.PushFactor {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Phone factor not available"),
		}
		return
	}

	// TODO

	return
}

func (s *Secondary) Passcode(db *database.Database, passcode string) (
	errData *errortypes.ErrorData, err error) {

	provider, err := s.GetProvider()
	if err != nil {
		return
	}

	if !provider.PushFactor {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Passcode factor not available"),
		}
		return
	}

	// TODO

	return
}

func (s *Secondary) Sms(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	provider, err := s.GetProvider()
	if err != nil {
		return
	}

	if !provider.PushFactor {
		err = &errortypes.AuthenticationError{
			errors.New("secondary: Sms factor not available"),
		}
		return
	}

	// TODO

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
		Push:     provider.PushFactor,
		Phone:    provider.PhoneFactor,
		Passcode: provider.PasscodeFactor,
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
	if provider.PasscodeFactor {
		factors = append(factors, "passcode")
	}
	if provider.SmsFactor {
		factors = append(factors, "sms")
	}

	query = fmt.Sprintf(
		"secondary=%s&factors=%s",
		s.Id,
		strings.Join(factors, ","),
	)

	return
}

func (s *Secondary) Handle(db *database.Database, factor, passcode string) (
	errData *errortypes.ErrorData, err error) {

	switch factor {
	case Push:
		errData, err = s.Push(db)
		break
	case Phone:
		errData, err = s.Phone(db)
		break
	case Passcode:
		errData, err = s.Passcode(db, passcode)
		break
	case Sms:
		errData, err = s.Sms(db)
		break
	default:
		err = &errortypes.UnknownError{
			errors.New("secondary: Unknown secondary factor"),
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
