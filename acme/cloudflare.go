package acme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/secret"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	cloudflareClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

type Cloudflare struct {
	token       string
	cacheZoneId map[string]string
}

type CloudflareError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CloudflareResponse struct {
	Errors []CloudflareError `json:"errors"`
}

type CloudflareZones struct {
	Result []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
	Errors []CloudflareError `json:"errors"`
}

type CloudflareRecord struct {
	Id      string `json:"id,omitempty"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
}

func (c *Cloudflare) Connect(secr *secret.Secret) (err error) {
	if secr.Type != secret.Cloudflare {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Secret type not Cloudflare"),
		}
		return
	}

	c.token = utils.FilterStr(secr.Key, 256)
	c.cacheZoneId = map[string]string{}

	return
}

func (c *Cloudflare) DnsZoneFind(domain string) (zoneId string, err error) {
	domain = extractDomain(domain)

	zoneId = c.cacheZoneId[domain]
	if zoneId != "" {
		return
	}

	u, err := url.Parse("https://api.cloudflare.com/client/v4/zones")
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "acme: Failed to parse Cloudflare URL"),
		}
		return
	}

	params := url.Values{}
	params.Add("name", domain)
	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Failed to create cloudflare request"),
		}
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := cloudflareClient.Do(req)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Cloudflare request failed"),
		}
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err = &errortypes.ApiError{
			errors.Wrapf(err, "acme: Cloudflare request bad status %s",
				cloudflareGetError(resp.StatusCode, body)),
		}
		return
	}

	zones := &CloudflareZones{}

	err = json.Unmarshal(body, &zones)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "acme: Failed to parse Cloudflare response"),
		}
		return
	}

	for _, zone := range zones.Result {
		if strings.Trim(zone.Name, ".") == strings.Trim(domain, ".") {
			zoneId = zone.ID
			break
		}
	}

	if zoneId == "" {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Cloudflare zone not found"),
		}
		return
	}

	c.cacheZoneId[domain] = zoneId

	return
}

type CloudflareRecords struct {
	Result []CloudflareRecord `json:"result"`
}

func (c *Cloudflare) DnsRecordFind(domain string) (recordId string,
	err error) {

	domain = cleanDomain(domain)

	zoneId, err := c.DnsZoneFind(domain)
	if err != nil {
		return
	}

	u, err := url.Parse(fmt.Sprintf(
		"https://api.cloudflare.com/client/v4/zones/%s/dns_records",
		zoneId,
	))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "acme: Failed to parse Cloudflare URL"),
		}
		return
	}

	params := url.Values{}
	params.Add("name", domain)
	params.Add("type", "TXT")
	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Failed to create cloudflare request"),
		}
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := cloudflareClient.Do(req)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Cloudflare request failed"),
		}
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err = &errortypes.ApiError{
			errors.Wrapf(err, "acme: Cloudflare request bad status %s",
				cloudflareGetError(resp.StatusCode, body)),
		}
		return
	}

	records := &CloudflareRecords{}

	err = json.Unmarshal(body, &records)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "acme: Failed to parse Cloudflare response"),
		}
		return
	}

	for _, record := range records.Result {
		if cleanDomain(record.Name) == domain && record.Type == "TXT" {
			recordId = record.Id
			break
		}
	}

	if recordId == "" {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Cloudflare record not found"),
		}
		return
	}

	return
}

func (c *Cloudflare) DnsTxtUpsert(domain, val string) (err error) {
	domain = cleanDomain(domain)

	zoneId, err := c.DnsZoneFind(domain)
	if err != nil {
		return
	}

	u, err := url.Parse(fmt.Sprintf(
		"https://api.cloudflare.com/client/v4/zones/%s/dns_records",
		zoneId,
	))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "acme: Failed to parse Cloudflare URL"),
		}
		return
	}

	record := CloudflareRecord{
		Type:    "TXT",
		Name:    domain,
		Content: "\"" + val + "\"",
		Ttl:     settings.Acme.DnsCloudflareTtl,
	}

	data, err := json.Marshal(record)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "acme: Failed to marshal cloudflare request"),
		}
		return
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(data))
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Failed to create cloudflare request"),
		}
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := cloudflareClient.Do(req)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Cloudflare request failed"),
		}
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err = &errortypes.ApiError{
			errors.Wrapf(err, "acme: Cloudflare request bad status %s",
				cloudflareGetError(resp.StatusCode, body)),
		}
		return
	}

	return
}

func (c *Cloudflare) DnsTxtDelete(domain, val string) (err error) {
	domain = cleanDomain(domain)

	zoneId, err := c.DnsZoneFind(domain)
	if err != nil {
		return
	}

	recordId, err := c.DnsRecordFind(domain)
	if err != nil {
		return
	}

	u, err := url.Parse(fmt.Sprintf(
		"https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s",
		zoneId,
		recordId,
	))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "acme: Failed to parse Cloudflare URL"),
		}
		return
	}

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Failed to create cloudflare request"),
		}
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := cloudflareClient.Do(req)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Cloudflare request failed"),
		}
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err = &errortypes.ApiError{
			errors.Wrapf(err, "acme: Cloudflare request bad status %s",
				cloudflareGetError(resp.StatusCode, body)),
		}
		return
	}

	return
}

func cloudflareGetError(status int, body []byte) (msg string) {
	msg = fmt.Sprintf("[%d]", status)

	cfResp := &CloudflareResponse{}

	_ = json.Unmarshal(body, &cfResp)

	if cfResp.Errors == nil {
		return
	}

	for _, cfErr := range cfResp.Errors {
		msg += fmt.Sprintf(" %d: %s", cfErr.Code, cfErr.Message)
		return
	}

	return
}
