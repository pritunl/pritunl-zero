package signature

import (
	"crypto/hmac"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/nonce"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"strconv"
	"strings"
	"time"
)

type Signature struct {
	Token     string
	Nonce     string
	Timestamp time.Time
	Signature string
	Method    string
	Path      string
	user      *user.User
}

func (s *Signature) GetUser(db *database.Database) (
	usr *user.User, err error) {

	if s.user != nil || db == nil || s.Token == "" {
		usr = s.user
		return
	}

	usr, err = user.GetTokenUpdate(db, s.Token)
	if err != nil {
		return
	}

	s.user = usr

	return
}

func (s *Signature) Validate(db *database.Database) (err error) {
	if s.Token == "" {
		err = &errortypes.AuthenticationError{
			errors.New("signature: Invalid authentication token"),
		}
		return
	}

	if len(s.Nonce) < 16 || len(s.Nonce) > 128 {
		err = &errortypes.AuthenticationError{
			errors.New("signature: Invalid authentication nonce"),
		}
		return
	}

	if time.Since(s.Timestamp) > time.Duration(
		settings.Auth.Window)*time.Second {

		err = &errortypes.AuthenticationError{
			errors.New("signature: Authentication timestamp outside window"),
		}
		return
	}

	usr, err := s.GetUser(db)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			usr = nil
			err = nil
			break
		default:
			return
		}
	}

	if usr == nil || usr.Token == "" || usr.Secret == "" {
		err = &errortypes.AuthenticationError{
			errors.New("signature: User not found"),
		}
		return
	}

	authString := strings.Join([]string{
		usr.Token,
		strconv.FormatInt(s.Timestamp.Unix(), 10),
		s.Nonce,
		s.Method,
		s.Path,
	}, "&")

	err = nonce.Validate(db, s.Nonce)
	if err != nil {
		return
	}

	hashFunc := hmac.New(sha512.New, []byte(usr.Secret))
	hashFunc.Write([]byte(authString))
	rawSignature := hashFunc.Sum(nil)
	sig := base64.StdEncoding.EncodeToString(rawSignature)

	if subtle.ConstantTimeCompare([]byte(s.Signature), []byte(sig)) != 1 {
		err = &errortypes.AuthenticationError{
			errors.New("signature: Invalid signature"),
		}
		return
	}

	return
}
