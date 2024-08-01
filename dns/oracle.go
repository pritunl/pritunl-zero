package dns

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/oracle/oci-go-sdk/v65/dns"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/secret"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/sirupsen/logrus"
)

type Oracle struct {
	token       string
	cacheZoneId map[string]string
	provider    *secret.OracleProvider
}

func (o *Oracle) OracleUser() string {
	return ""
}

func (o *Oracle) OraclePrivateKey() string {
	return ""
}

func (o *Oracle) Connect(db *database.Database,
	secr *secret.Secret) (err error) {

	if secr.Type != secret.OracleCloud {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Secret type not Oracle Cloud"),
		}
		return
	}

	o.cacheZoneId = map[string]string{}

	o.provider, err = secr.GetOracleProvider()
	if err != nil {
		return
	}

	return
}

func (o *Oracle) DnsZoneFind(db *database.Database, domain string) (
	zoneId string, err error) {

	domain = extractDomain(domain)

	zoneId = o.cacheZoneId[domain]
	if zoneId != "" {
		return
	}

	compartmentId, err := o.provider.CompartmentOCID()
	if err != nil {
		return
	}

	req := dns.ListZonesRequest{
		CompartmentId: &compartmentId,
	}

	client, err := o.provider.GetDnsClient()
	if err != nil {
		return
	}

	zones, err := client.ListZones(db, req)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Oracle zone list error"),
		}
		return
	}

	for _, zone := range zones.Items {
		if matchDomains(*zone.Name, domain) {
			zoneId = *zone.Id
			break
		}
	}

	if zoneId == "" {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Oracle zone not found"),
		}
		return
	}

	o.cacheZoneId[domain] = zoneId

	return
}

func (o *Oracle) DnsCommit(db *database.Database,
	domain, recordType string, ops []*Operation) (err error) {

	zoneName := extractDomain(domain)
	domain = cleanDomain(domain)

	items := []dns.RecordOperation{}

	values := set.NewSet()
	oracleOps := []string{}

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

		values.Add(op.Value)

		switch op.Operation {
		case RETAIN:
			break
		case UPSERT:
			oracleOps = append(oracleOps, "add:"+op.Value)
			items = append(items, dns.RecordOperation{
				Domain: &domain,
				Rtype:  utils.PointerString(recordType),
				Ttl: utils.PointerInt(
					settings.Acme.DnsOracleCloudTtl),
				Rdata:     utils.PointerString(op.Value),
				Operation: dns.RecordOperationOperationAdd,
			})
			break
		case DELETE:
			oracleOps = append(oracleOps, "remove:"+op.Value)
			items = append(items, dns.RecordOperation{
				Domain: &domain,
				Rtype:  utils.PointerString(recordType),
				Ttl: utils.PointerInt(
					settings.Acme.DnsOracleCloudTtl),
				Rdata:     utils.PointerString(op.Value),
				Operation: dns.RecordOperationOperationRemove,
			})
			break
		}
	}

	client, err := o.provider.GetDnsClient()
	if err != nil {
		return
	}

	getReq := dns.GetZoneRecordsRequest{
		ZoneNameOrId: utils.PointerString(zoneName),
		Domain:       utils.PointerString(domain),
	}

	resp, err := client.GetZoneRecords(db, getReq)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Oracle zone record get error"),
		}
		return
	}

	for _, record := range resp.Items {
		if record.Rtype != nil && *record.Rtype == recordType &&
			record.Rdata != nil {

			val := *record.Rdata
			if recordType == "AAAA" {
				val = normalizeIp(val)
			}

			if val == "" {
				continue
			}

			if values.Contains(val) {
				continue
			}

			oracleOps = append(oracleOps, "remove_unknown:"+*record.Rdata)
			items = append(items, dns.RecordOperation{
				Domain: &domain,
				Rtype:  utils.PointerString(recordType),
				Ttl: utils.PointerInt(
					settings.Acme.DnsOracleCloudTtl),
				Rdata:     utils.PointerString(*record.Rdata),
				Operation: dns.RecordOperationOperationRemove,
			})
		}
	}

	logrus.WithFields(logrus.Fields{
		"domain":     domain,
		"operations": oracleOps,
	}).Info("domain: Oracle dns batch operation")

	req := dns.PatchZoneRecordsRequest{
		ZoneNameOrId: utils.PointerString(zoneName),
		PatchZoneRecordsDetails: dns.PatchZoneRecordsDetails{
			Items: items,
		},
	}

	_, err = client.PatchZoneRecords(db, req)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Oracle zone patch error"),
		}
		return
	}

	return
}

func (o *Oracle) DnsFind(db *database.Database,
	domain, recordType string) (vals []string, err error) {

	zoneName := extractDomain(domain)
	domain = cleanDomain(domain)

	req := dns.GetZoneRecordsRequest{
		ZoneNameOrId: utils.PointerString(zoneName),
		Domain:       utils.PointerString(domain),
	}

	client, err := o.provider.GetDnsClient()
	if err != nil {
		return
	}

	resp, err := client.GetZoneRecords(db, req)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Oracle zone record get error"),
		}
		return
	}

	for _, record := range resp.Items {
		if record.Rtype != nil && *record.Rtype == recordType &&
			record.Rdata != nil {

			val := *record.Rdata
			if recordType == "AAAA" {
				val = normalizeIp(val)
			}

			if val == "" {
				continue
			}

			vals = append(vals, val)
		}
	}

	return
}

func (o *Oracle) DnsTxtGet(db *database.Database,
	domain string) (vals []string, err error) {

	zoneName := extractDomain(domain)
	domain = cleanDomain(domain)

	req := dns.GetZoneRecordsRequest{
		ZoneNameOrId: utils.PointerString(zoneName),
		Domain:       utils.PointerString(domain),
	}

	client, err := o.provider.GetDnsClient()
	if err != nil {
		return
	}

	resp, err := client.GetZoneRecords(db, req)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Oracle zone record get error"),
		}
		return
	}

	for _, record := range resp.Items {
		if record.Rtype != nil && *record.Rtype == "TXT" &&
			record.Rdata != nil {

			vals = append(vals, *record.Rdata)
		}
	}

	return
}

func (o *Oracle) DnsTxtUpsert(db *database.Database,
	domain, val string) (err error) {

	zoneName := extractDomain(domain)
	domain = cleanDomain(domain)

	req := dns.PatchZoneRecordsRequest{
		ZoneNameOrId: utils.PointerString(zoneName),
		PatchZoneRecordsDetails: dns.PatchZoneRecordsDetails{
			Items: []dns.RecordOperation{
				{
					Domain: &domain,
					Rtype:  utils.PointerString("TXT"),
					Ttl: utils.PointerInt(
						settings.Acme.DnsOracleCloudTtl),
					Rdata:     utils.PointerString(val),
					Operation: dns.RecordOperationOperationAdd,
				},
			},
		},
	}

	client, err := o.provider.GetDnsClient()
	if err != nil {
		return
	}

	_, err = client.PatchZoneRecords(db, req)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Oracle zone patch error"),
		}
		return
	}

	return
}

func (o *Oracle) DnsTxtDelete(db *database.Database,
	domain, val string) (err error) {

	zoneName := extractDomain(domain)
	domain = cleanDomain(domain)

	req := dns.PatchZoneRecordsRequest{
		ZoneNameOrId: utils.PointerString(zoneName),
		PatchZoneRecordsDetails: dns.PatchZoneRecordsDetails{
			Items: []dns.RecordOperation{
				{
					Domain:    &domain,
					Rtype:     utils.PointerString("TXT"),
					Rdata:     utils.PointerString(val),
					Operation: dns.RecordOperationOperationRemove,
				},
			},
		},
	}

	client, err := o.provider.GetDnsClient()
	if err != nil {
		return
	}

	_, err = client.PatchZoneRecords(db, req)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Oracle zone patch error"),
		}
		return
	}

	return
}
