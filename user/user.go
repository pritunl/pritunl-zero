package user

import (
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/requires"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id            bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Type          string        `bson:"type" json:"type"`
	Username      string        `bson:"username" json:"username"`
	Password      string        `bson:"password" json:"-"`
	Roles         []string      `bson:"roles" json:"roles"`
	Administrator string        `bson:"administrator" json:"administrator"`
	Permissions   []string      `bson:"permissions" json:"permissions"`
}

func (u *User) Commit(db *database.Database) (err error) {
	coll := db.Users()

	err = coll.Commit(u.Id, u)
	if err != nil {
		return
	}

	return
}

func (u *User) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Users()

	err = coll.CommitFields(u.Id, u, fields)
	if err != nil {
		return
	}

	return
}

func (u *User) Insert(db *database.Database) (
	err error) {

	coll := db.Users()

	if u.Id != "" {
		err = &errortypes.DatabaseError{
			errors.New("user: User already exists"),
		}
		return
	}

	err = coll.Insert(u)
	if err != nil {
		return
	}

	return
}

func (u *User) SetPassword(password string) (err error) {
	if u.Type != "local" {
		err = &errortypes.UnknownError{
			errors.New("user: User type cannot store password"),
		}
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "user: Failed to hash password"),
		}
		return
	}

	u.Password = string(hash)

	return
}

func (u *User) CheckPassword(password string) bool {
	if u.Type != "local" || u.Password == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func init() {
	module := requires.New("user")
	module.After("settings")

	module.Handler = func() (err error) {
		db := database.GetDatabase()
		defer db.Close()

		exists, err := HasSuper(db)
		if err != nil {
			return
		}

		if !exists {
			logrus.Info("setup: Creating default super user")

			usr := User{
				Type:          "local",
				Username:      "pritunl",
				Administrator: "super",
			}

			err = usr.SetPassword("pritunl")
			if err != nil {
				return
			}

			err = usr.Insert(db)
			if err != nil {
				return
			}
		}

		return
	}
}
