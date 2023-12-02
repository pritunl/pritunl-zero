package auth

import (
	"fmt"
	"sort"

	"github.com/pritunl/pritunl-zero/settings"
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

	if !settings.Local.NoLocalAuth {
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

func GetFastAdminPath() (path string) {
	if !settings.Local.NoLocalAuth || !settings.Auth.FastLogin ||
		len(settings.Auth.Providers) == 0 {

		return
	}

	if len(settings.Auth.Providers) > 1 {
		googleOnly := true
		for _, provider := range settings.Auth.Providers {
			if provider.Type != Google {
				googleOnly = false
				break
			}
		}

		if !googleOnly {
			return
		}
	}

	for _, provider := range settings.Auth.Providers {
		if provider.Type == Google {
			path = fmt.Sprintf("/auth/request?id=%s", Google)
		} else {
			path = fmt.Sprintf("/auth/request?id=%s", provider.Id.Hex())
		}
		return
	}

	return
}

func GetFastUserPath() (path string) {
	if settings.Auth.FastLogin && settings.Auth.ForceFastUserLogin &&
		len(settings.Auth.Providers) != 0 {
	} else if !settings.Local.NoLocalAuth || !settings.Auth.FastLogin ||
		len(settings.Auth.Providers) == 0 {

		return
	}

	if len(settings.Auth.Providers) > 1 {
		googleOnly := true
		for _, provider := range settings.Auth.Providers {
			if provider.Type != Google {
				googleOnly = false
				break
			}
		}

		if !googleOnly {
			return
		}
	}

	for _, provider := range settings.Auth.Providers {
		if provider.Type == Google {
			path = fmt.Sprintf("/auth/request?id=%s", Google)
		} else {
			path = fmt.Sprintf("/auth/request?id=%s", provider.Id.Hex())
		}
		return
	}

	return
}

func GetFastServicePath() (path string) {
	if settings.Auth.FastLogin && settings.Auth.ForceFastServiceLogin &&
		len(settings.Auth.Providers) != 0 {
	} else if !settings.Local.NoLocalAuth || !settings.Auth.FastLogin ||
		len(settings.Auth.Providers) == 0 {

		return
	}

	if len(settings.Auth.Providers) > 1 {
		googleOnly := true
		for _, provider := range settings.Auth.Providers {
			if provider.Type != Google {
				googleOnly = false
				break
			}
		}

		if !googleOnly {
			return
		}
	}

	for _, provider := range settings.Auth.Providers {
		if provider.Type == Google {
			path = fmt.Sprintf("/auth/request?id=%s", Google)
		} else {
			path = fmt.Sprintf("/auth/request?id=%s", provider.Id.Hex())
		}
		return
	}

	return
}
