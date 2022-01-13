package auth

import (
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
)

func SyncUser(db *database.Database, usr *user.User) (
	active bool, err error) {

	if time.Since(usr.LastSync) < time.Duration(
		settings.Auth.Sync)*time.Second {

		active = true
		return
	}

	provider := settings.Auth.GetProvider(usr.Provider)

	if usr.Type == user.AuthZero && provider != nil &&
		provider.Type == user.AuthZero {

		active, err = AuthZeroSync(db, usr, provider)
		if err != nil {
			return
		}
	} else if usr.Type == user.Azure && provider != nil &&
		provider.Type == user.Azure {

		active, err = AzureSync(db, usr, provider)
		if err != nil {
			return
		}
	} else if usr.Type == user.Google {
		active, err = GoogleSync(db, usr)
		if err != nil {
			return
		}
	} else if usr.Type == user.JumpCloud {
		active, err = JumpcloudSync(db, usr, provider)
		if err != nil {
			return
		}
	} else {
		active = true
	}

	if active {
		usr.LastSync = time.Now()
		err = usr.CommitFields(db, set.NewSet("last_sync"))
		if err != nil {
			return
		}
	}

	return
}
