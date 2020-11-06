package user

import (
	"sort"
	"strings"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type            string             `bson:"type" json:"type"`
	Username        string             `bson:"username" json:"username"`
	Password        string             `bson:"password" json:"-"`
	DefaultPassword string             `bson:"default_password" json:"-"`
	Token           string             `bson:"token" json:"token"`
	Secret          string             `bson:"secret" json:"secret"`
	Theme           string             `bson:"theme" json:"-"`
	LastActive      time.Time          `bson:"last_active" json:"last_active"`
	LastSync        time.Time          `bson:"last_sync" json:"last_sync"`
	Roles           []string           `bson:"roles" json:"roles"`
	Administrator   string             `bson:"administrator" json:"administrator"`
	Disabled        bool               `bson:"disabled" json:"disabled"`
	ActiveUntil     time.Time          `bson:"active_until" json:"active_until"`
	Permissions     []string           `bson:"permissions" json:"permissions"`
}

func (u *User) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if u.Roles == nil {
		u.Roles = []string{}
	}

	if u.Permissions == nil {
		u.Permissions = []string{}
	}

	if !types.Contains(u.Type) {
		errData = &errortypes.ErrorData{
			Error:   "user_type_invalid",
			Message: "User type is not valid",
		}
		return
	}

	if u.Username == "" {
		errData = &errortypes.ErrorData{
			Error:   "user_username_invalid",
			Message: "User username is not valid",
		}
		return
	}

	if u.Type == Local && u.Password == "" {
		errData = &errortypes.ErrorData{
			Error:   "user_password_missing",
			Message: "User password is not set",
		}
		return
	}

	u.Format()

	return
}

func (u *User) SuperExists(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if u.Administrator != "super" && !u.Id.IsZero() {
		exists, e := hasSuperSkip(db, u.Id)
		if e != nil {
			err = e
			return
		}

		if !exists {
			errData = &errortypes.ErrorData{
				Error:   "user_missing_super",
				Message: "Missing super administrator",
			}
			return
		}
	}

	return
}

func (u *User) Format() {
	if u.Type == Local {
		u.Username = strings.ToLower(u.Username)
	}

	roles := []string{}
	rolesSet := set.NewSet()

	for _, role := range u.Roles {
		rolesSet.Add(role)
	}

	for role := range rolesSet.Iter() {
		roles = append(roles, role.(string))
	}

	sort.Strings(roles)

	u.Roles = roles
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

func (u *User) Insert(db *database.Database) (err error) {
	coll := db.Users()

	if !u.Id.IsZero() {
		err = &errortypes.DatabaseError{
			errors.New("user: User already exists"),
		}
		return
	}

	_, err = coll.InsertOne(db, u)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (u *User) Upsert(db *database.Database) (err error) {
	coll := db.Users()

	opts := &options.FindOneAndUpdateOptions{}
	opts.SetUpsert(true)
	opts.SetReturnDocument(options.After)

	err = coll.FindOneAndUpdate(
		db,
		&bson.M{
			"type":     u.Type,
			"username": u.Username,
		},
		&bson.M{
			"$setOnInsert": u,
		},
		opts,
	).Decode(u)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (u *User) RolesMatch(roles []string) bool {
	usrRoles := set.NewSet()
	for _, role := range u.Roles {
		usrRoles.Add(role)
	}

	for _, role := range roles {
		if usrRoles.Contains(role) {
			return true
		}
	}

	return false
}

func (u *User) RolesMerge(roles []string) bool {
	newRoles := set.NewSet()
	curRoles := set.NewSet()

	for _, role := range roles {
		newRoles.Add(role)
	}

	for _, role := range u.Roles {
		newRoles.Add(role)
		curRoles.Add(role)
	}

	if !curRoles.IsEqual(newRoles) {
		rls := []string{}

		for role := range newRoles.Iter() {
			rls = append(rls, role.(string))
		}

		u.Roles = rls
		return true
	}

	return false
}

func (u *User) RolesOverwrite(roles []string) bool {
	newRoles := set.NewSet()
	curRoles := set.NewSet()

	for _, role := range roles {
		newRoles.Add(role)
	}

	for _, role := range u.Roles {
		curRoles.Add(role)
	}

	if !curRoles.IsEqual(newRoles) {
		u.Roles = roles
		return true
	}

	return false
}

func (u *User) SetPassword(password string) (err error) {
	if u.Type != Local {
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
	u.DefaultPassword = ""

	return
}

func (u *User) GenerateDefaultPassword() (err error) {
	passwd, err := utils.RandStr(12)
	if err != nil {
		return
	}

	err = u.SetPassword(passwd)
	if err != nil {
		return
	}

	u.DefaultPassword = passwd

	return
}

func (u *User) CheckPassword(password string) bool {
	if u.Type != Local || u.Password == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func (u *User) GenerateToken() (err error) {
	u.Token, err = utils.RandStr(48)
	if err != nil {
		return
	}

	u.Secret, err = utils.RandStr(48)
	if err != nil {
		return
	}

	return
}

func init() {
	module := requires.New("user")
	module.After("settings")

	module.Handler = func() (err error) {
		db := database.GetDatabase()
		defer db.Close()

		coll := db.Users()

		cursor, err := coll.Find(db, &bson.M{})
		if err != nil {
			err = database.ParseError(err)
			return
		}
		defer cursor.Close(db)

		for cursor.Next(db) {
			usr := &User{}
			err = cursor.Decode(usr)
			if err != nil {
				err = database.ParseError(err)
				return
			}

			newUsername := strings.ToLower(usr.Username)
			if usr.Username != newUsername {
				err = coll.UpdateId(usr.Id, &bson.M{
					"$set": &bson.M{
						"username": newUsername,
					},
				})
				if err != nil {
					return
				}
			}
		}

		count, err := Count(db)
		if err != nil {
			return
		}

		if count == 0 {
			logrus.Info("user: Creating default super user")

			usr := User{
				Type:          Local,
				Username:      "pritunl",
				Administrator: "super",
			}

			err = usr.GenerateDefaultPassword()
			if err != nil {
				return
			}

			_, err = usr.Validate(db)
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
