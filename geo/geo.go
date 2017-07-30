package geo

import (
	"bytes"
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"net/http"
	"time"
)

var (
	client = &http.Client{
		Timeout: 3 * time.Second,
	}
)

type Geo struct {
	Address   string    `bson:"_id" json:"address"`
	Isp       string    `bson:"i" json:"isp"`
	Country   string    `bson:"c" json:"country"`
	Region    string    `bson:"r" json:"region"`
	City      string    `bson:"a" json:"city"`
	Longitude float64   `bson:"x" json:"longitude"`
	Latitude  float64   `bson:"y" json:"latitude"`
	Timestamp time.Time `bson:"t" json:"-"`
}

type geoData struct {
	License string `json:"license"`
	Address string `json:"address"`
}

func get(addr string) (ge *Geo, err error) {
	reqGeoData := &geoData{
		License: settings.System.License,
		Address: addr,
	}

	reqData := &bytes.Buffer{}
	err = json.NewEncoder(reqData).Encode(reqGeoData)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "geo: Failed to parse request data"),
		}
		return
	}

	req, err := http.NewRequest(
		"GET",
		settings.Auth.Server+"/geo",
		reqData,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "geo: Failed to create request"),
		}
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "geo: Failed to send request"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = &errortypes.ParseError{
			errors.Newf("geo: Request failed with status %d", resp.StatusCode),
		}
		return
	}

	ge = &Geo{}
	err = json.NewDecoder(resp.Body).Decode(ge)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "geo: Failed to parse response"),
		}
		return
	}

	return
}

func Get(db *database.Database, addr string) (ge *Geo, err error) {
	ge = &Geo{}
	coll := db.Geo()

	err = coll.FindOneId(addr, ge)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			ge = nil
			err = nil
		default:
			return
		}
	}

	if ge == nil {
		ge, err = get(addr)
		if err != nil {
			return
		}

		ge.Timestamp = time.Now()

		coll.Insert(ge)
	}

	return
}
