package auth

import (
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"gopkg.in/mgo.v2/bson"
)

type StateProvider struct {
	Id    bson.ObjectId `json:"id"`
	Type  string        `json:"type"`
	Label string        `json:"label"`
}

type State struct {
	Providers []*StateProvider `json:"providers"`
}

func GetState() (state *State) {
	state = &State{
		Providers: []*StateProvider{},
	}

	for _, provider := range settings.Auth.Providers {
		provider := &StateProvider{
			Id:    provider.Id,
			Type:  provider.Type,
			Label: provider.Label,
		}
		state.Providers = append(state.Providers, provider)
	}

	return
}

func Local(db *database.Database, username, password string) (
	usr *user.User, errData *errortypes.ErrorData, err error) {

	usr, err = user.GetUsername(db, user.Local, username)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			err = nil
			errData = &errortypes.ErrorData{
				Error:   "auth_invalid",
				Message: "Authencation credentials are invalid",
			}
			break
		}
		return
	}

	valid := usr.CheckPassword(password)
	if !valid {
		errData = &errortypes.ErrorData{
			Error:   "auth_invalid",
			Message: "Authencation credentials are invalid",
		}
		return
	}

	return
}
