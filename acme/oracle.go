package acme

import (
	"strings"

	"github.com/dropbox/godropbox/errors"
	"github.com/oracle/oci-go-sdk/v55/dns"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/secret"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
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
		if strings.Trim(*zone.Name, ".") == strings.Trim(domain, ".") {
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
					Rdata:     utils.PointerString("\"" + val + "\""),
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
					Rdata:     utils.PointerString("\"" + val + "\""),
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
