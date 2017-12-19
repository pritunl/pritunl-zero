package keybase

import (
	"encoding/json"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"net/http"
	"time"
)

var (
	client = &http.Client{
		Timeout: 20 * time.Second,
	}
)

type Info struct {
	Username string `json:"username"`
	Picture  string `json:"picture"`
	Twitter  string `json:"twitter"`
	Github   string `json:"github"`
}

type infoStatus struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

type infoBasics struct {
	Username string `json:"username"`
}

type infoProof struct {
	ProofType string `json:"proof_type"`
	Name      string `json:"nametag"`
}

type infoProofTypes struct {
	Github  []infoProof `json:"github"`
	Twitter []infoProof `json:"twitter"`
}

type infoProofs struct {
	ByProofType infoProofTypes `json:"by_proof_type"`
}

type infoPrimary struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Source string `json:"source"`
}

type infoPictures struct {
	Primary infoPrimary `json:"primary"`
}

type infoThem struct {
	Basics   infoBasics   `json:"basics"`
	Proofs   infoProofs   `json:"proofs_summary"`
	Pictures infoPictures `json:"pictures"`
}

type infoResp struct {
	Status infoStatus `json:"status"`
	Them   infoThem   `json:"them"`
}

func getInfo(username string) (data *infoResp, err error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://keybase.io/_/api/1.0/user/lookup.json?"+
			"username=%s&fields=basics,proofs_summary,pictures", username),
		nil,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "keybase: Info request failed"),
		}
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "keybase: Info request failed"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = &errortypes.RequestError{
			errors.Wrapf(err, "keybase: Info request bad status %d",
				resp.StatusCode),
		}
		return
	}

	data = &infoResp{}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(
				err, "keybase: Failed to parse info response",
			),
		}
		return
	}

	return
}
