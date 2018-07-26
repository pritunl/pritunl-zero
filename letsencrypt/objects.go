package letsencrypt

import (
	"time"

	"github.com/square/go-jose"
)

const (
	StatusPending = "pending"
	StatusInvalid = "invalid"
	StatusValid   = "valid"
)

// Challenge represents a server challenge for a given domain name.
type Challenge struct {
	ID        int64     `json:"id,omitempty"`
	Type      string    `json:"type"`
	URI       string    `json:"uri"`
	Status    string    `json:"status,omitempty"`
	Validated time.Time `json:"validated,omitempty"`
	Error     *Error    `json:"error,omitempty"`

	// Data used by various challenges
	Token            string           `json:"token,omitempty"`
	KeyAuthorization string           `json:"keyAuthorization,omitempty"`
	N                int              `json:"n,omitempty"`
	Certs            []string         `json:"certs,omitempty"`
	AccountKey       *jose.JsonWebKey `json:"accountKey,omitempty"`
	Authorization    *JWSValidation   `json:"authorization,omitempty"`
}

type JWSValidation struct {
	Header    *jose.JsonWebKey `json:"header,omitempty"`
	Payload   string           `json:"payload,omitempty"`
	Signature string           `json:"signature,omitempty"`
}

// Authorization represents a set of challenges issued by the server
// for the given identifier.
type Authorization struct {
	Identifier struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"identifier"`

	Status     string      `json:"status,omitempty"`
	Expires    time.Time   `json:"expires,omitempty"`
	Challenges []Challenge `json:"challenges,omitempty"`
	Combs      [][]int     `json:"combinations,omitempty"`
}

// Combinations returns the set of challenges which the client supports.
// Completing one of these sets is enough to prove ownership of an identifier.
func (a Authorization) Combinations(supportedChallenges ...string) [][]Challenge {
	supported := func(chal Challenge) bool {
		for _, c := range supportedChallenges {
			if c == chal.Type {
				return true
			}
		}
		return false
	}

	chals := [][]Challenge{}
	for _, comb := range a.Combs {
		chalList := make([]Challenge, len(comb))
		sup := true
		for i, idx := range comb {
			if idx >= 0 && len(a.Challenges) > idx && supported(a.Challenges[idx]) {
				chalList[i] = a.Challenges[idx]
			} else {
				sup = false
				break
			}
		}
		if sup {
			chals = append(chals, chalList)
		}
	}
	return chals
}

// Registration holds account information for a given key pair.
type Registration struct {
	PublicKey      *jose.JsonWebKey `json:"key,omitempty"`
	Contact        []string         `json:"contact,omitempty"`
	Agreement      string           `json:"agreement,omitempty"`
	Authorizations string           `json:"authorizations,omitempty"`
	Certificates   string           `json:"certificates,omitempty"`

	Id        int       `json:"id,omitempty"`
	InitialIp string    `json:"initialIp,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`

	Resource string `json:"resource,omitempty"`
}
