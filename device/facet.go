package device

import (
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
