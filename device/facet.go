package device

import (
	"fmt"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/settings"
)

type FacetVersion struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

type TrustedFacet struct {
	Ids     []string      `json:"ids"`
	Version *FacetVersion `json:"version"`
}

type Facets struct {
	TrustedFacets []*TrustedFacet `json:"trustedFacets"`
}

func GetFacets() (facets *Facets) {
	return &Facets{
		TrustedFacets: []*TrustedFacet{
			&TrustedFacet{
				Ids: settings.Local.Facets,
				Version: &FacetVersion{
					Major: 1,
					Minor: 0,
				},
			},
		},
	}
}

func GetAppId() string {
	return fmt.Sprintf("https://%s/auth/u2f/app.json", node.Self.UserDomain)
}
