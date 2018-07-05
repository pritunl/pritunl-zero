package device

import (
	"crypto/ecdsa"
	"crypto/x509"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/u2flib"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type Device struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id"`
	User         bson.ObjectId `bson:"user" json:"user"`
	Name         string        `bson:"name" json:"name"`
	Type         string        `bson:"type" json:"type"`
	Mode         string        `bson:"mode" json:"mode"`
	Timestamp    time.Time     `bson:"timestamp" json:"timestamp"`
	Disabled     bool          `bson:"disabled" json:"disabled"`
	ActiveUntil  time.Time     `bson:"activeactive_until_until" json:"active_until"`
	LastActive   time.Time     `bson:"last_active" json:"last_active"`
	SshPublicKey string        `bson:"ssh_public_key" json:"ssh_public_key"`
	U2fRaw       []byte        `bson:"u2f_raw" json:"-"`
	U2fCounter   uint32        `bson:"u2f_counter" json:"-"`
	U2fKeyHandle []byte        `bson:"u2f_key_handle" json:"-"`
	U2fPublicKey []byte        `bson:"u2f_public_key" json:"-"`
}

func (d *Device) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

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

	if d.Type != U2f && d.Type != SmartCard {
		errData = &errortypes.ErrorData{
			Error:   "device_type_invalid",
			Message: "Device type is invalid",
		}
		return
	}

	if d.Mode != Ssh && d.Mode != Secondary {
		errData = &errortypes.ErrorData{
			Error:   "device_mode_invalid",
			Message: "Device mode is invalid",
		}
		return
	}

	if d.Mode == Ssh {
		if d.Type != SmartCard {
			errData = &errortypes.ErrorData{
				Error:   "device_mode_type_invalid",
				Message: "Device mode and type is invalid",
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

	err = coll.Insert(d)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
