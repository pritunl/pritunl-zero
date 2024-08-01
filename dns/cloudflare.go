package dns

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/secret"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/sirupsen/logrus"
)

type Cloudflare struct {
	sess        *cloudflare.API
	token       string
	cacheZoneId map[string]string
}

func (c *Cloudflare) Connect(db *database.Database,
	secr *secret.Secret) (err error) {

	if secr.Type != secret.Cloudflare {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Secret type not cloudflare"),
		}
		return
	}

	c.sess, err = cloudflare.NewWithAPIToken(utils.FilterStr(secr.Key, 256))
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "dns: Failed to connect to cloudflare api"),
		}
		return
	}

	c.cacheZoneId = map[string]string{}

	return
}

func (c *Cloudflare) DnsZoneFind(db *database.Database, domain string) (
	zoneId string, err error) {

	domain = extractDomain(domain)

	zoneId = c.cacheZoneId[domain]
	if zoneId != "" {
		return
	}

	zones, err := c.sess.ListZones(db)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Failed to list cloudflare zones"),
		}
		return
	}

	for _, zone := range zones {
		if matchDomains(zone.Name, domain) {
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

func (c *Cloudflare) DnsCommit(db *database.Database,
	domain, recordType string, ops []*Operation) (err error) {

	domain = cleanDomain(domain)

	zoneId, err := c.DnsZoneFind(db, domain)
	if err != nil {
		return
	}

	listParams := cloudflare.ListDNSRecordsParams{
		Type: recordType,
		Name: domain,
	}

	records, _, err := c.sess.ListDNSRecords(
		db,
		cloudflare.ZoneIdentifier(zoneId),
		listParams,
	)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "acme: Failed to get DNS records"),
		}
		return
	}

	recordIds := map[string]string{}
	for _, record := range records {
		if record.Type == recordType && matchDomains(record.Name, domain) {
			val := record.Content
			if recordType == "AAAA" {
				val = normalizeIp(val)
			}

			if val == "" {
				continue
			}

			recordIds[val] = record.ID
		}
	}

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
	}

	for _, op := range ops {
		if op.Operation != DELETE {
			continue
		}

		recordId := recordIds[op.Value]
		if recordId == "" {
			continue
		}
		delete(recordIds, op.Value)

		logrus.WithFields(logrus.Fields{
			"operation": "delete",
			"record_id": recordId,
			"domain":    domain,
			"value":     op.Value,
		}).Info("domain: Cloudflare dns operation")

		err = c.sess.DeleteDNSRecord(
			db,
			cloudflare.ZoneIdentifier(zoneId),
			recordId,
		)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "dns: Failed to delete record"),
			}
			return
		}
	}

	for _, op := range ops {
		if op.Operation != RETAIN && op.Operation != UPSERT {
			continue
		}

		recordId := recordIds[op.Value]
		if recordId == "" {
			continue
		}
		delete(recordIds, op.Value)

		op.Operation = ""
	}

	for _, op := range ops {
		if op.Operation != RETAIN && op.Operation != UPSERT {
			continue
		}

		updateVal := ""
		recordId := ""
		for val, recId := range recordIds {
			updateVal = val
			recordId = recId
			break
		}
		if recordId != "" {
			delete(recordIds, updateVal)
		}

		if recordId == "" {
			logrus.WithFields(logrus.Fields{
				"operation": "create",
				"domain":    domain,
				"value":     op.Value,
			}).Info("domain: Cloudflare dns operation")

			createParams := cloudflare.CreateDNSRecordParams{
				Type:    recordType,
				Name:    domain,
				Content: op.Value,
				TTL:     settings.Acme.DnsCloudflareTtl,
			}

			_, err = c.sess.CreateDNSRecord(
				db,
				cloudflare.ZoneIdentifier(zoneId),
				createParams,
			)
			if err != nil {
				err = &errortypes.ParseError{
					errors.Wrap(err, "acme: Failed to create record"),
				}
				return
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"operation": "update",
				"record_id": recordId,
				"domain":    domain,
				"value":     op.Value,
			}).Info("domain: Cloudflare dns operation")

			updateParams := cloudflare.UpdateDNSRecordParams{
				ID:      recordId,
				Type:    recordType,
				Name:    domain,
				Content: op.Value,
				TTL:     settings.Acme.DnsCloudflareTtl,
			}

			_, err = c.sess.UpdateDNSRecord(
				db,
				cloudflare.ZoneIdentifier(zoneId),
				updateParams,
			)
			if err != nil {
				err = &errortypes.ParseError{
					errors.Wrap(err, "acme: Failed to update record"),
				}
				return
			}
		}
	}

	for val, recordId := range recordIds {
		logrus.WithFields(logrus.Fields{
			"operation": "delete_unknown",
			"record_id": recordId,
			"domain":    domain,
			"value":     val,
		}).Info("domain: Cloudflare dns operation")

		err = c.sess.DeleteDNSRecord(
			db,
			cloudflare.ZoneIdentifier(zoneId),
			recordId,
		)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "dns: Failed to delete record"),
			}
			return
		}
	}

	return
}

func (c *Cloudflare) DnsFind(db *database.Database,
	domain, recordType string) (vals []string, err error) {

	vals = []string{}
	domain = cleanDomain(domain)

	zoneId, err := c.DnsZoneFind(db, domain)
	if err != nil {
		return
	}

	listParams := cloudflare.ListDNSRecordsParams{
		Type: recordType,
		Name: domain,
	}

	records, _, err := c.sess.ListDNSRecords(
		db,
		cloudflare.ZoneIdentifier(zoneId),
		listParams,
	)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "acme: Failed to get DNS records"),
		}
		return
	}

	for _, record := range records {
		if record.Type == recordType && matchDomains(record.Name, domain) {
			val := record.Content
			if recordType == "AAAA" {
				val = normalizeIp(val)
			}

			if val == "" {
				continue
			}

			vals = append(vals, val)
			break
		}
	}

	return
}

func (c *Cloudflare) DnsTxtGet(db *database.Database,
	domain string) (vals []string, err error) {

	vals = []string{}

	domain = cleanDomain(domain)

	zoneId, err := c.DnsZoneFind(db, domain)
	if err != nil {
		return
	}

	listParams := cloudflare.ListDNSRecordsParams{
		Type: "TXT",
		Name: domain,
	}

	records, _, err := c.sess.ListDNSRecords(
		db,
		cloudflare.ZoneIdentifier(zoneId),
		listParams,
	)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "acme: Failed to get DNS records"),
		}
		return
	}

	for _, record := range records {
		if record.Type == "TXT" && matchDomains(record.Name, domain) {
			vals = append(vals, record.Content)
			break
		}
	}

	return
}

func (c *Cloudflare) DnsTxtUpsert(db *database.Database,
	domain, val string) (err error) {

	domain = cleanDomain(domain)

	zoneId, err := c.DnsZoneFind(db, domain)
	if err != nil {
		return
	}

	listParams := cloudflare.ListDNSRecordsParams{
		Type: "TXT",
		Name: domain,
	}

	records, _, err := c.sess.ListDNSRecords(
		db,
		cloudflare.ZoneIdentifier(zoneId),
		listParams,
	)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "acme: Failed to get DNS records"),
		}
		return
	}

	recordId := ""
	for _, record := range records {
		if record.Type == "TXT" && matchDomains(record.Name, domain) {
			recordId = record.ID
			break
		}
	}

	if recordId == "" {
		createParams := cloudflare.CreateDNSRecordParams{
			Type:    "TXT",
			Name:    domain,
			Content: val,
			TTL:     settings.Acme.DnsCloudflareTtl,
		}

		_, err = c.sess.CreateDNSRecord(
			db,
			cloudflare.ZoneIdentifier(zoneId),
			createParams,
		)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "acme: Failed to create DNS record"),
			}
			return
		}
	} else {
		updateParams := cloudflare.UpdateDNSRecordParams{
			Type:    "TXT",
			Name:    domain,
			Content: val,
			TTL:     settings.Acme.DnsCloudflareTtl,
		}

		_, err = c.sess.UpdateDNSRecord(
			db,
			cloudflare.ResourceIdentifier(recordId),
			updateParams,
		)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "acme: Failed to update DNS record"),
			}
			return
		}
	}

	return
}

func (c *Cloudflare) DnsTxtDelete(db *database.Database,
	domain, val string) (err error) {

	domain = cleanDomain(domain)

	zoneId, err := c.DnsZoneFind(db, domain)
	if err != nil {
		return
	}

	listParams := cloudflare.ListDNSRecordsParams{
		Type: "TXT",
		Name: domain,
	}

	records, _, err := c.sess.ListDNSRecords(
		db,
		cloudflare.ZoneIdentifier(zoneId),
		listParams,
	)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "acme: Failed to get DNS records"),
		}
		return
	}

	recordId := ""
	for _, record := range records {
		if record.Type == "TXT" &&
			matchDomains(record.Name, domain) &&
			matchTxt(record.Content, val) {

			recordId = record.ID
			break
		}
	}

	if recordId != "" {
		err = c.sess.DeleteDNSRecord(
			db,
			cloudflare.ZoneIdentifier(zoneId),
			recordId,
		)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "acme: Failed to delete DNS record"),
			}
			return
		}
	}

	return
}
