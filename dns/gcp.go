package dns

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/secret"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

type Gcp struct {
	service     *dns.Service
	project     string
	cacheZoneId map[string]string
}

type gcpCredentials struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

func (g *Gcp) Connect(db *database.Database,
	secr *secret.Secret) (err error) {

	if secr.Type != secret.GCP {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Secret type not GCP"),
		}
		return
	}

	g.cacheZoneId = map[string]string{}

	// Parse service account JSON from Key field
	var creds gcpCredentials
	if secr.Key != "" {
		err = json.Unmarshal([]byte(secr.Key), &creds)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "acme: Failed to parse GCP credentials"),
			}
			return
		}
		g.project = creds.ProjectID
	} else {
		err = &errortypes.ParseError{
			errors.New("acme: GCP project ID not found"),
		}
		return
	}

	ctx := context.Background()

	// Create DNS service with service account credentials
	opts := []option.ClientOption{
		option.WithCredentialsJSON([]byte(secr.Key)),
	}

	g.service, err = dns.NewService(ctx, opts...)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Failed to create GCP DNS service"),
		}
		return
	}

	return
}

func (g *Gcp) DnsZoneFind(db *database.Database, domain string) (
	zoneId string, err error) {

	domainClean := strings.Trim(domain, ".")

	// Check cache using original domain
	zoneId = g.cacheZoneId[domainClean]
	if zoneId != "" {
		return
	}

	ctx := context.Background()

	zones, err := g.service.ManagedZones.List(g.project).Context(ctx).Do()
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: GCP zone list error"),
		}
		return
	}

	// Try to match domain against zones, starting with most specific (longest) zones first
	// Sort zones by DNS name length (longest first) to match most specific zones first
	type zoneInfo struct {
		name         string
		dnsName      string
		dnsNameClean string
	}
	zoneList := []zoneInfo{}
	for _, zone := range zones.ManagedZones {
		if zone.DnsName != "" {
			zoneList = append(zoneList, zoneInfo{
				name:         zone.Name,
				dnsName:      zone.DnsName,
				dnsNameClean: strings.Trim(zone.DnsName, "."),
			})
		}
	}

	// Sort by DNS name length (longest first)
	for i := 0; i < len(zoneList); i++ {
		for j := i + 1; j < len(zoneList); j++ {
			if len(zoneList[i].dnsNameClean) < len(zoneList[j].dnsNameClean) {
				zoneList[i], zoneList[j] = zoneList[j], zoneList[i]
			}
		}
	}

	for _, zone := range zoneList {
		// Check if domain matches zone DNS name exactly
		if matchDomains(zone.dnsNameClean, domainClean) {
			zoneId = zone.name
			break
		}

		// Check if domain is a subdomain of the zone (e.g., sub.example.com matches example.com zone)
		if strings.HasSuffix(domainClean, "."+zone.dnsNameClean) {
			zoneId = zone.name
			break
		}
	}

	if zoneId == "" {
		zoneNames := []string{}
		zoneDnsNames := []string{}
		for _, zone := range zones.ManagedZones {
			if zone.DnsName != "" {
				zoneNames = append(zoneNames, zone.Name)
				zoneDnsNames = append(zoneDnsNames, zone.DnsName)
			}
		}
		logrus.WithFields(logrus.Fields{
			"domain":                   domainClean,
			"original_domain":          domain,
			"available_zone_names":     zoneNames,
			"available_zone_dns_names": zoneDnsNames,
		}).Error("acme: GCP zone not found")

		err = &errortypes.ApiError{
			errors.Wrapf(err,
				"acme: GCP zone not found for domain '%s'. Available zones: %v",
				domain, zoneDnsNames),
		}
		return
	}

	g.cacheZoneId[domainClean] = zoneId

	return
}

func (g *Gcp) DnsCommit(db *database.Database,
	domain, recordType string, ops []*Operation) (err error) {

	domain = cleanDomain(domain)

	zoneId, err := g.DnsZoneFind(db, domain)
	if err != nil {
		return
	}

	ctx := context.Background()

	// Get existing records
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

	// Build change request
	change := &dns.Change{
		Additions: []*dns.ResourceRecordSet{},
		Deletions: []*dns.ResourceRecordSet{},
	}

	operations := []string{}

	// Find existing record set for this domain and type
	var existingRrSet *dns.ResourceRecordSet
	for _, rrSet := range existingRecords.Rrsets {
		if matchDomains(rrSet.Name, domain+".") && rrSet.Type == recordType {
			existingRrSet = rrSet
			break
		}
	}

	// Process operations
	valuesToAdd := []string{}
	valuesToDelete := []string{}

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
			valuesToAdd = append(valuesToAdd, op.Value)
			operations = append(operations, "add:"+op.Value)
		case DELETE:
			valuesToDelete = append(valuesToDelete, op.Value)
			operations = append(operations, "remove:"+op.Value)
		}
	}

	// If there's an existing record set, we need to handle it
	if existingRrSet != nil {
		// Delete existing record set if we're modifying it
		if len(valuesToAdd) > 0 || len(valuesToDelete) > 0 {
			change.Deletions = append(change.Deletions, existingRrSet)
		}

		// Add back values that should be retained
		if existingRrSet.Rrdatas != nil {
			for _, existingVal := range existingRrSet.Rrdatas {
				shouldRetain := false
				for _, op := range ops {
					if op.Operation == RETAIN && op.Value == existingVal {
						shouldRetain = true
						break
					}
				}
				if shouldRetain {
					valuesToAdd = append(valuesToAdd, existingVal)
				}
			}
		}
	}

	// Create new record set with values to add
	if len(valuesToAdd) > 0 {
		ttl := int64(settings.Acme.DnsGcpTtl)
		if existingRrSet != nil && existingRrSet.Ttl > 0 {
			ttl = existingRrSet.Ttl
		}

		newRrSet := &dns.ResourceRecordSet{
			Name:    domain + ".",
			Type:    recordType,
			Ttl:     ttl,
			Rrdatas: valuesToAdd,
		}
		change.Additions = append(change.Additions, newRrSet)
	}

	if len(change.Additions) == 0 && len(change.Deletions) == 0 {
		return
	}

	logrus.WithFields(logrus.Fields{
		"domain":     domain,
		"operations": operations,
	}).Info("domain: GCP dns batch operation")

	_, err = g.service.Changes.Create(g.project, zoneId, change).Context(ctx).Do()
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: GCP record change error"),
		}
		return
	}

	return
}

func (g *Gcp) DnsFind(db *database.Database, domain, recordType string) (
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

	for _, rrSet := range records.Rrsets {
		if matchDomains(rrSet.Name, domain+".") && rrSet.Type == recordType {
			for _, rdata := range rrSet.Rrdatas {
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

func (g *Gcp) DnsTxtGet(db *database.Database, domain string) (
	vals []string, err error) {

	vals = []string{}

	zoneId, err := g.DnsZoneFind(db, domain)
	if err != nil {
		return
	}

	ctx := context.Background()

	listCall := g.service.ResourceRecordSets.List(g.project, zoneId)
	listCall.Name(cleanDomain(domain) + ".")
	listCall.Type("TXT")

	records, err := listCall.Context(ctx).Do()
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: GCP TXT record list error"),
		}
		return
	}

	for _, rrSet := range records.Rrsets {
		if matchDomains(rrSet.Name, cleanDomain(domain)+".") && rrSet.Type == "TXT" {
			for _, rdata := range rrSet.Rrdatas {
				vals = append(vals, rdata)
			}
		}
	}

	return
}

func (g *Gcp) DnsTxtUpsert(db *database.Database,
	domain, val string) (err error) {

	zoneId, err := g.DnsZoneFind(db, domain)
	if err != nil {
		return
	}

	ctx := context.Background()

	domain = cleanDomain(domain)

	// Get existing TXT record
	listCall := g.service.ResourceRecordSets.List(g.project, zoneId)
	listCall.Name(domain + ".")
	listCall.Type("TXT")

	existingRecords, err := listCall.Context(ctx).Do()
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: GCP TXT record list error"),
		}
		return
	}

	change := &dns.Change{}

	// If existing record exists, delete it first
	var existingRrSet *dns.ResourceRecordSet
	for _, rrSet := range existingRecords.Rrsets {
		if matchDomains(rrSet.Name, domain+".") && rrSet.Type == "TXT" {
			existingRrSet = rrSet
			break
		}
	}

	if existingRrSet != nil {
		change.Deletions = append(change.Deletions, existingRrSet)
	}

	// Add new record
	ttl := int64(settings.Acme.DnsGcpTtl)
	if existingRrSet != nil && existingRrSet.Ttl > 0 {
		ttl = existingRrSet.Ttl
	}

	newRrSet := &dns.ResourceRecordSet{
		Name:    domain + ".",
		Type:    "TXT",
		Ttl:     ttl,
		Rrdatas: []string{val},
	}
	change.Additions = append(change.Additions, newRrSet)

	_, err = g.service.Changes.Create(g.project, zoneId, change).Context(ctx).Do()
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: GCP TXT record upsert error"),
		}
		return
	}

	return
}

func (g *Gcp) DnsTxtDelete(db *database.Database,
	domain, val string) (err error) {

	zoneId, err := g.DnsZoneFind(db, domain)
	if err != nil {
		return
	}

	ctx := context.Background()

	domain = cleanDomain(domain)

	// Get existing TXT record
	listCall := g.service.ResourceRecordSets.List(g.project, zoneId)
	listCall.Name(domain + ".")
	listCall.Type("TXT")

	existingRecords, err := listCall.Context(ctx).Do()
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: GCP TXT record list error"),
		}
		return
	}

	// Find and delete matching record
	for _, rrSet := range existingRecords.Rrsets {
		if matchDomains(rrSet.Name, domain+".") && rrSet.Type == "TXT" {
			// Check if value matches
			for _, rdata := range rrSet.Rrdatas {
				if matchTxt(rdata, val) {
					change := &dns.Change{
						Deletions: []*dns.ResourceRecordSet{rrSet},
					}

					_, err = g.service.Changes.Create(g.project, zoneId, change).Context(ctx).Do()
					if err != nil {
						err = &errortypes.ApiError{
							errors.Wrap(err, "acme: GCP TXT record delete error"),
						}
						return
					}

					return
				}
			}
		}
	}

	return
}
