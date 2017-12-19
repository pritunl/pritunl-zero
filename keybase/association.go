package keybase

import (
	"fmt"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/ssh"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Association struct {
	Id        string    `bson:"_id"`
	Type      string    `bson:"type"`
	Username  string    `bson:"username"`
	Timestamp time.Time `bson:"timestamp"`
	State     string    `bson:"state"`
}

func (a *Association) Message() string {
	return fmt.Sprintf(
		"%s&%s&%s",
		a.Id,
		a.Type,
		a.Username,
	)
}

func (a *Association) Validate(signature string) (
	err error, errData *errortypes.ErrorData) {

	valid, err := VerifySig(a.Message(), signature, a.Username)
	if err != nil {
		return
	}

	if !valid {
		errData = &errortypes.ErrorData{
			Error:   "invalid_signature",
			Message: "Keybase signature is invalid",
		}
		return
	}

	return
}

func (a *Association) Approve(db *database.Database,
	usr *user.User) (err error, errData *errortypes.ErrorData) {

	keybaseMode, err := policy.UserKeybaseMode(db, usr)
	if err != nil {
		return
	}

	if keybaseMode == policy.Disabled {
		err = a.Deny(db, usr)
		if err != nil {
			return
		}

		errData = &errortypes.ErrorData{
			Error:   "keybase_disabled",
			Message: "Keybase cannot be used with this user",
		}
		return
	}

	if a.State != "" {
		err = errortypes.WriteError{
			errors.New("keybase: Association has already been answered"),
		}
		return
	}

	if usr.Keybase != "" {
		err = a.Deny(db, usr)
		if err != nil {
			return
		}

		errData = &errortypes.ErrorData{
			Error:   "keybase_associated",
			Message: "Keybase already associated with this user",
		}
		return
	}

	data, err := getInfo(a.Username)
	if err != nil {
		return
	}
	keybaseId := data.Them.PublicKeys.Primary.UkbId

	coll := db.Users()

	err = coll.Update(&bson.M{
		"_id": usr.Id,
		"$or": []*bson.M{
			&bson.M{
				"keybase": &bson.M{
					"$exists": false,
				},
			},
			&bson.M{
				"keybase": "",
			},
		},
	}, &bson.M{
		"$set": &bson.M{
			"keybase":    a.Username,
			"keybase_id": keybaseId,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	usr.Keybase = a.Username
	usr.KeybaseId = a.Username

	a.State = ssh.Approved

	coll = db.KeybaseChallenges()

	err = coll.Update(&bson.M{
		"_id":   a.Id,
		"state": "",
	}, a)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (a *Association) Deny(db *database.Database, usr *user.User) (err error) {
	if a.State != "" {
		err = errortypes.WriteError{
			errors.New("keybase: Association has already been answered"),
		}
		return
	}

	a.State = ssh.Denied

	coll := db.KeybaseChallenges()

	err = coll.Update(&bson.M{
		"_id":   a.Id,
		"state": "",
	}, a)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (a *Association) GetInfo() (info *Info, err error) {
	data, err := getInfo(a.Username)
	if err != nil {
		return
	}

	info = &Info{
		Username: data.Them.Basics.Username,
		Picture:  data.Them.Pictures.Primary.Url,
	}

	if len(data.Them.Proofs.ByProofType.Twitter) > 0 {
		info.Twitter = data.Them.Proofs.ByProofType.Twitter[0].Name
	}

	if len(data.Them.Proofs.ByProofType.Github) > 0 {
		info.Github = data.Them.Proofs.ByProofType.Github[0].Name
	}

	return
}

func (a *Association) Commit(db *database.Database) (err error) {
	coll := db.KeybaseChallenges()

	err = coll.Commit(a.Id, a)
	if err != nil {
		return
	}

	return
}

func (a *Association) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.KeybaseChallenges()

	err = coll.CommitFields(a.Id, a, fields)
	if err != nil {
		return
	}

	return
}

func (a *Association) Insert(db *database.Database) (err error) {
	coll := db.KeybaseChallenges()

	err = coll.Insert(a)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func NewAssociation(db *database.Database, username string) (
	asc *Association, err error) {

	token, err := utils.RandStr(32)
	if err != nil {
		return
	}

	asc = &Association{
		Id:        token,
		Type:      AssociationChallenge,
		Timestamp: time.Now(),
		Username:  username,
	}

	err = asc.Insert(db)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAssociation(db *database.Database, ascId string) (
	asc *Association, err error) {

	coll := db.KeybaseChallenges()
	asc = &Association{}

	err = coll.FindOneId(ascId, asc)
	if err != nil {
		return
	}

	return
}
