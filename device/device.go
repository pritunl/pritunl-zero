package device

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"strings"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/u2flib"
	"github.com/pritunl/pritunl-zero/utils"
)

type Device struct {
	Id                 primitive.ObjectID      `bson:"_id,omitempty" json:"id"`
	User               primitive.ObjectID      `bson:"user" json:"user"`
	Name               string                  `bson:"name" json:"name"`
	Type               string                  `bson:"type" json:"type"`
	Mode               string                  `bson:"mode" json:"mode"`
	Timestamp          time.Time               `bson:"timestamp" json:"timestamp"`
	Disabled           bool                    `bson:"disabled" json:"disabled"`
	ActiveUntil        time.Time               `bson:"activeactive_until_until" json:"active_until"`
	LastActive         time.Time               `bson:"last_active" json:"last_active"`
	AlertLevels        []int                   `bson:"alert_levels" json:"alert_levels"`
	Number             string                  `bson:"number" json:"number"`
	SshPublicKey       string                  `bson:"ssh_public_key" json:"ssh_public_key"`
	U2fRaw             []byte                  `bson:"u2f_raw" json:"-"`
	U2fCounter         uint32                  `bson:"u2f_counter" json:"-"`
	U2fKeyHandle       []byte                  `bson:"u2f_key_handle" json:"-"`
	U2fPublicKey       []byte                  `bson:"u2f_public_key" json:"-"`
	WanId              []byte                  `bson:"wan_id" json:"-"`
	WanPublicKey       []byte                  `bson:"wan_public_key" json:"-"`
	WanAttestationType string                  `bson:"wan_attestation_type" json:"-"`
	WanAuthenticator   *webauthn.Authenticator `bson:"wan_authenticator" json:"-"`
	WanRpId            string                  `bson:"wan_rp_id" json:"wan_rp_id"`
}

func (d *Device) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	d.Name = utils.FilterName(d.Name)

	if len(d.Name) == 0 {
		errData = &errortypes.ErrorData{
			Error:   "device_name_missing",
			Message: "Device name is required",
		}
		return
	}

	if len(d.Name) > 22 {
		errData = &errortypes.ErrorData{
			Error:   "device_name_invalid",
			Message: "Device name is too long",
		}
		return
	}

	switch d.Mode {
	case Ssh:
		if d.Type != SmartCard {
			errData = &errortypes.ErrorData{
				Error:   "device_type_invalid",
				Message: "Device type is invalid",
			}
			return
		}

		if d.SshPublicKey == "" {
			errData = &errortypes.ErrorData{
				Error:   "device_ssh_key_missing",
				Message: "Device SSH public key is required",
			}
			return
		}

		if !strings.Contains(d.SshPublicKey, "cardno") {
			errData = &errortypes.ErrorData{
				Error:   "device_ssh_key_invalid",
				Message: "Device SSH public key is not from a Smart Card",
			}
			return
		}
		break
	case Secondary:
		if d.Type != U2f && d.Type != WebAuthn {
			errData = &errortypes.ErrorData{
				Error:   "device_type_invalid",
				Message: "Device type is invalid",
			}
			return
		}
		break
	case Phone:
		if d.Type != Call && d.Type != Message {
			errData = &errortypes.ErrorData{
				Error:   "device_type_invalid",
				Message: "Device type is invalid",
			}
			return
		}

		if len(d.Number) == 10 {
			d.Number = "+1" + d.Number
		}

		if len(d.Number) < 10 {
			errData = &errortypes.ErrorData{
				Error:   "device_number_invalid",
				Message: "Device phone number invalid",
			}
			return
		}

		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "device_mode_invalid",
			Message: "Device mode is invalid",
		}
		return
	}

	if d.AlertLevels == nil {
		d.AlertLevels = []int{}
	}
	for _, level := range d.AlertLevels {
		switch level {
		case Low, Medium, High:
			break
		default:
			errData = &errortypes.ErrorData{
				Error:   "device_alert_level_invalid",
				Message: "Device alert level is invalid",
			}
			return
		}
	}

	return
}

func (d *Device) SetActive(db *database.Database) (err error) {
	d.LastActive = time.Now()
	err = d.CommitFields(db, set.NewSet("last_active"))
	if err != nil {
		return
	}

	return
}

func (d *Device) MarshalWebauthn(cred *webauthn.Credential) {
	if d.Type == U2f {
		d.U2fCounter = cred.Authenticator.SignCount
	} else {
		d.WanId = cred.ID
		d.WanPublicKey = cred.PublicKey
		d.WanAttestationType = cred.AttestationType
		d.WanAuthenticator = &cred.Authenticator
	}

	return
}

func (d *Device) UnmarshalWebauthn() (cred webauthn.Credential, err error) {
	if d.Type == U2f {
		pubKeyItf, e := x509.ParsePKIXPublicKey(d.U2fPublicKey)
		if e != nil {
			err = &errortypes.ParseError{
				errors.Wrap(e, "device: Failed to parse device public key"),
			}
			return
		}

		pubKey, ok := pubKeyItf.(*ecdsa.PublicKey)
		if !ok {
			err = &errortypes.ParseError{
				errors.New("device: Device public key invalid type"),
			}
			return
		}

		pubKeyByte := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)

		cred = webauthn.Credential{
			ID:              d.U2fKeyHandle,
			PublicKey:       pubKeyByte,
			AttestationType: "fido-u2f",
			Authenticator: webauthn.Authenticator{
				AAGUID:    d.U2fRaw,
				SignCount: d.U2fCounter,
			},
		}
		return
	}

	cred = webauthn.Credential{
		ID:              d.WanId,
		PublicKey:       d.WanPublicKey,
		AttestationType: d.WanAttestationType,
		Authenticator:   *d.WanAuthenticator,
	}
	return
}

func (d *Device) MarshalRegistration(reg *u2flib.Registration) (err error) {
	pubPkix, err := x509.MarshalPKIXPublicKey(&reg.PubKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "device: Failed to marshal device public key"),
		}
		return
	}

	d.U2fRaw = reg.Raw
	d.U2fKeyHandle = reg.KeyHandle
	d.U2fPublicKey = pubPkix

	return
}

func (d *Device) UnmarshalRegistration() (
	reg u2flib.Registration, err error) {

	pubKeyItf, err := x509.ParsePKIXPublicKey(d.U2fPublicKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "device: Failed to parse device public key"),
		}
		return
	}

	pubKey, ok := pubKeyItf.(*ecdsa.PublicKey)
	if !ok {
		err = &errortypes.ParseError{
			errors.Wrap(err, "device: Device public key invalid type"),
		}
		return
	}

	reg = u2flib.Registration{
		KeyHandle: d.U2fKeyHandle,
		PubKey:    *pubKey,
	}

	return
}

func (d *Device) CheckLevel(level int) bool {
	if d.AlertLevels == nil {
		return false
	}

	for _, lvl := range d.AlertLevels {
		if level == lvl {
			return true
		}
	}

	return false
}

func (d *Device) Commit(db *database.Database) (err error) {
	coll := db.Devices()

	err = coll.Commit(d.Id, d)
	if err != nil {
		return
	}

	return
}

func (d *Device) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Devices()

	err = coll.CommitFields(d.Id, d, fields)
	if err != nil {
		return
	}

	return
}

func (d *Device) Insert(db *database.Database) (err error) {
	coll := db.Devices()

	_, err = coll.InsertOne(db, d)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
