package dns

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/secret"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

type Google struct {
	service     *dns.Service
	project     string
	cacheZoneId map[string]string
}

type googleKey struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

type googleZoneInfo struct {
	name         string
	dnsName      string
	dnsNameClean string
}

func (g *Google) Connect(db *database.Database,
	secr *secret.Secret) (err error) {

	if secr.Type != secret.GoogleCloud {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Secret type not GCP"),
		}
		return
	}

	g.cacheZoneId = map[string]string{}

	googleKey := &googleKey{}
	if secr.Key != "" {
		err = json.Unmarshal([]byte(secr.Key), googleKey)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(
					err,
					"acme: Failed to parse Google Cloud credentials",
				),
			}
			return
		}
		g.project = googleKey.ProjectId
	} else {
		err = &errortypes.ParseError{
			errors.New("acme: GCP project ID not found"),
		}
		return
	}

	ctx := context.Background()
	opts := []option.ClientOption{
		option.WithCredentialsJSON([]byte(secr.Key)),
	}

	g.service, err = dns.NewService(ctx, opts...)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(
				err,
				"acme: Failed to create Google Cloud DNS service",
			),
		}
		return
	}

	return
}

func (g *Google) DnsZoneFind(db *database.Database, domain string) (
	zoneId string, err error) {

	domainClean := strings.Trim(domain, ".")

	zoneId = g.cacheZoneId[domainClean]
	if zoneId != "" {
		return
	}

	ctx := context.Background()
	zones, err := g.service.ManagedZones.List(g.project).Context(ctx).Do()
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Google Cloud zone list error"),
		}
		return
	}

	zoneList := []*googleZoneInfo{}
	for _, zone := range zones.ManagedZones {
		if zone.DnsName != "" {
			zoneList = append(zoneList, &googleZoneInfo{
				name:         zone.Name,
				dnsName:      zone.DnsName,
				dnsNameClean: strings.Trim(zone.DnsName, "."),
			})
		}
	}

	for i := 0; i < len(zoneList); i++ {
		for j := i + 1; j < len(zoneList); j++ {
			if len(zoneList[i].dnsNameClean) < len(zoneList[j].dnsNameClean) {
				zoneList[i], zoneList[j] = zoneList[j], zoneList[i]
			}
		}
	}

	for _, zone := range zoneList {
		if matchDomains(zone.dnsNameClean, domainClean) {
			zoneId = zone.name
			break
		}

		if strings.HasSuffix(domainClean, "."+zone.dnsNameClean) {
			zoneId = zone.name
			break
		}
	}

	if zoneId == "" {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Google Cloud DNS zone not found"),
		}
		return
	}

	g.cacheZoneId[domainClean] = zoneId

	return
}

func (g *Google) DnsCommit(db *database.Database,
	domain, recordType string, ops []*Operation) (err error) {

	domain = cleanDomain(domain)

	zoneId, err := g.DnsZoneFind(db, domain)
	if err != nil {
		return
	}

	ctx := context.Background()

	listCall := g.service.ResourceRecordSets.List(g.project, zoneId)
	listCall.Name(domain + ".")
	listCall.Type(recordType)

	existingRecords, err := listCall.Context(ctx).Do()
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: GCP record list error"),
		}
		return
	}

	change := &dns.Change{
		Additions: []*dns.ResourceRecordSet{},
		Deletions: []*dns.ResourceRecordSet{},
	}

	operations := []string{}

	var existingRecSet *dns.ResourceRecordSet
	for _, recSet := range existingRecords.Rrsets {
		if matchDomains(recSet.Name, domain+".") && recSet.Type == recordType {
			existingRecSet = recSet
			break
		}
	}

	addValues := []string{}
	removeValues := []string{}

	for _, op := range ops {
		if recordType == "AAAA" {
			val := normalizeIp(op.Value)
			if val == "" {
				err = &errortypes.ParseError{
					errors.Newf("dns: Invalid ipv6 address %s", op.Value),
				}
				return
			}
			op.Value = val
		}

		switch op.Operation {
		case UPSERT, RETAIN:
			operations = append(operations, "add:"+op.Value)
			addValues = append(addValues, op.Value)
		case DELETE:
			removeValues = append(removeValues, op.Value)
			operations = append(operations, "remove:"+op.Value)
		}
	}

	existingValues := set.NewSet()
	if existingRecSet != nil {
		for _, value := range existingRecSet.Rrdatas {
			existingValues.Add(value)
		}
	}

	newValues := existingValues.Copy()
	for _, value := range removeValues {
		newValues.Remove(value)
	}

	for _, value := range addValues {
		newValues.Add(value)
	}

	if existingValues.IsEqual(newValues) {
		return
	}

	if existingRecSet != nil {
		change.Deletions = append(change.Deletions, existingRecSet)
	}

	if newValues.Len() > 0 {
		ttl := int64(settings.Acme.DnsGoogleCloudTtl)
		if existingRecSet != nil && existingRecSet.Ttl > 0 {
			ttl = existingRecSet.Ttl
		}

		values := []string{}
		for valueInf := range newValues.Iter() {
			values = append(values, valueInf.(string))
		}
		sort.Strings(values)

		newRrSet := &dns.ResourceRecordSet{
			Name:    domain + ".",
			Type:    recordType,
			Ttl:     ttl,
			Rrdatas: values,
		}
		change.Additions = append(change.Additions, newRrSet)
	}

	logrus.WithFields(logrus.Fields{
		"domain":     domain,
		"operations": operations,
	}).Info("domain: Google Cloud dns batch operation")

	_, err = g.service.Changes.Create(g.project, zoneId, change).Context(ctx).Do()
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Google Cloud record change error"),
		}
		return
	}

	return
}

func (g *Google) DnsFind(db *database.Database, domain, recordType string) (
	vals []string, err error) {

	vals = []string{}
	domain = cleanDomain(domain)

	zoneId, err := g.DnsZoneFind(db, domain)
	if err != nil {
		return
	}

	ctx := context.Background()

	listCall := g.service.ResourceRecordSets.List(g.project, zoneId)
	listCall.Name(domain + ".")
	listCall.Type(recordType)

	records, err := listCall.Context(ctx).Do()
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: GCP record list error"),
		}
		return
	}

	for _, recSet := range records.Rrsets {
		if matchDomains(recSet.Name, domain+".") && recSet.Type == recordType {
			for _, rdata := range recSet.Rrdatas {
				val := rdata
				if recordType == "AAAA" {
					val = normalizeIp(val)
				}

				if val == "" {
					continue
				}

				vals = append(vals, val)
			}
		}
	}

	return
}
