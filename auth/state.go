package auth

import (
	"github.com/pritunl/pritunl-zero/settings"
	"sort"
)

type StateProvider struct {
	Id    interface{} `json:"id"`
	Type  string      `json:"type"`
	Label string      `json:"label"`
}

type StateProviders []*StateProvider

func (s StateProviders) Len() int {
	return len(s)
}

func (s StateProviders) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s StateProviders) Less(i, j int) bool {
	return s[i].Label < s[j].Label
}

type State struct {
	Providers StateProviders `json:"providers"`
}

func GetState() (state *State) {
	state = &State{
		Providers: StateProviders{},
	}

	if settings.Local.HasLocalAuth {
		prv := &StateProvider{
			Type: "local",
		}

		state.Providers = append(state.Providers, prv)
	}

	google := false

	for _, provider := range settings.Auth.Providers {
		prv := &StateProvider{
			Type:  provider.Type,
			Label: provider.Label,
		}

		if provider.Type == Google {
			if google {
				continue
			}
			google = true
			prv.Id = Google
		} else {
			prv.Id = provider.Id
		}

		state.Providers = append(state.Providers, prv)
	}

	sort.Sort(state.Providers)

	return
}
