package dns

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/secret"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/sirupsen/logrus"
)

type Aws struct {
	client      *route53.Client
	cacheZoneId map[string]string
}

func (a *Aws) Connect(db *database.Database,
	secr *secret.Secret) (err error) {

	if secr.Type != secret.AWS {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: Secret type not AWS"),
		}
		return
	}

	a.cacheZoneId = map[string]string{}

	a.client = route53.NewFromConfig(aws.Config{
		Region: secr.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			secr.Key, secr.Value, ""),
	})

	return
}

func (a *Aws) DnsZoneFind(domain string) (zoneId string, err error) {
	domain = extractDomain(domain)

	zoneId = a.cacheZoneId[domain]
	if zoneId != "" {
		return
	}

	input := &route53.ListHostedZonesInput{}

	result, err := a.client.ListHostedZones(context.Background(), input)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: AWS route53 zone lookup error"),
		}
		return
	}

	for _, zone := range result.HostedZones {
		if zone.Name != nil && matchDomains(*zone.Name, domain) {
			zoneId = *zone.Id
			break
		}
	}

	if zoneId == "" {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: AWS route53 zone not found"),
		}
		return
	}

	a.cacheZoneId[domain] = zoneId

	return
}

func (a *Aws) DnsCommit(db *database.Database,
	domain, recordType string, ops []*Operation) (err error) {

	domain = cleanDomain(domain)

	zoneId, err := a.DnsZoneFind(domain)
	if err != nil {
		return
	}

	deleteResourceRecs := []types.ResourceRecord{}
	updateResourceRecs := []types.ResourceRecord{}
	operations := []string{}
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

		resourceRec := types.ResourceRecord{
			Value: aws.String(op.Value),
		}

		switch op.Operation {
		case UPSERT, RETAIN:
			operations = append(operations, "add:"+op.Value)
			updateResourceRecs = append(updateResourceRecs, resourceRec)
		case DELETE:
			curVals, e := a.DnsFind(db, domain, recordType)
			if e != nil {
				err = e
				return
			}

			exists := false
			for _, val := range curVals {
				if val == op.Value {
					exists = true
					break
				}
			}

			if !exists {
				logrus.WithFields(logrus.Fields{
					"domain":    domain,
					"operation": "remove:" + op.Value,
				}).Info("domain: Skipping delete on changed record")
				continue
			}

			operations = append(operations, "remove:"+op.Value)
			deleteResourceRecs = append(deleteResourceRecs, resourceRec)
		}
	}

	logrus.WithFields(logrus.Fields{
		"domain":     domain,
		"operations": operations,
	}).Info("domain: AWS dns batch operation")

	if len(updateResourceRecs) == 0 && len(deleteResourceRecs) > 0 {
		input := &route53.ChangeResourceRecordSetsInput{
			ChangeBatch: &types.ChangeBatch{
				Changes: []types.Change{
					{
						Action: types.ChangeActionDelete,
						ResourceRecordSet: &types.ResourceRecordSet{
							Name: aws.String(domain),
							Type: types.RRType(recordType),
							TTL: aws.Int64(int64(
								settings.Acme.DnsAwsTtl)),
							ResourceRecords: deleteResourceRecs,
						},
					},
				},
				Comment: aws.String("Pritunl delete record"),
			},
			HostedZoneId: aws.String(zoneId),
		}

		_, err = a.client.ChangeResourceRecordSets(
			context.Background(), input)
		if err != nil {
			if strings.Contains(err.Error(), "delete") &&
				strings.Contains(err.Error(), "not found") {

				err = nil
			} else {
				err = &errortypes.ApiError{
					errors.Wrap(err, "acme: AWS record delete error"),
				}
				return
			}
		}
	}

	if len(updateResourceRecs) > 0 {
		input := &route53.ChangeResourceRecordSetsInput{
			ChangeBatch: &types.ChangeBatch{
				Changes: []types.Change{
					{
						Action: types.ChangeActionUpsert,
						ResourceRecordSet: &types.ResourceRecordSet{
							Name: aws.String(domain),
							Type: types.RRType(recordType),
							TTL: aws.Int64(int64(
								settings.Acme.DnsAwsTtl)),
							ResourceRecords: updateResourceRecs,
						},
					},
				},
				Comment: aws.String("Pritunl update record"),
			},
			HostedZoneId: aws.String(zoneId),
		}

		_, err = a.client.ChangeResourceRecordSets(
			context.Background(), input)
		if err != nil {
			err = &errortypes.ApiError{
				errors.Wrap(err, "acme: AWS record update error"),
			}
			return
		}
	}

	return
}

func (a *Aws) DnsFind(db *database.Database, domain, recordType string) (
	vals []string, err error) {

	vals = []string{}
	domain = cleanDomain(domain)

	zoneId, err := a.DnsZoneFind(domain)
	if err != nil {
		return
	}

	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneId),
		StartRecordName: aws.String(domain),
		StartRecordType: types.RRType(recordType),
	}

	result, err := a.client.ListResourceRecordSets(
		context.Background(), input)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: AWS record list error"),
		}
		return
	}

	for _, recordSet := range result.ResourceRecordSets {
		if string(recordSet.Type) == recordType &&
			recordSet.Name != nil && matchDomains(*recordSet.Name, domain) {

			for _, record := range recordSet.ResourceRecords {
				if record.Value != nil {
					val := *record.Value

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
	}

	return
}

func (a *Aws) DnsTxtGet(db *database.Database, domain string) (
	vals []string, err error) {

	vals = []string{}

	zoneId, err := a.DnsZoneFind(domain)
	if err != nil {
		return
	}

	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneId),
		StartRecordName: aws.String(domain),
		StartRecordType: types.RRTypeTxt,
	}

	result, err := a.client.ListResourceRecordSets(
		context.Background(), input)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: AWS route53 record set error"),
		}
		return
	}

	for _, recordSet := range result.ResourceRecordSets {
		if recordSet.Type == types.RRTypeTxt &&
			recordSet.Name != nil && matchDomains(*recordSet.Name, domain) {

			for _, record := range recordSet.ResourceRecords {
				if record.Value != nil {
					vals = append(vals, *record.Value)
				}
			}
		}
	}

	return
}

func (a *Aws) DnsTxtUpsert(db *database.Database,
	domain, val string) (err error) {

	zoneId, err := a.DnsZoneFind(domain)
	if err != nil {
		return
	}

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(domain),
						Type: types.RRTypeTxt,
						TTL:  aws.Int64(int64(settings.Acme.DnsAwsTtl)),
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(val),
							},
						},
					},
				},
			},
			Comment: aws.String("Pritunl update TXT record"),
		},
		HostedZoneId: aws.String(zoneId),
	}

	_, err = a.client.ChangeResourceRecordSets(context.Background(), input)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: AWS route53 record set error"),
		}
		return
	}

	return
}

func (a *Aws) DnsTxtDelete(db *database.Database,
	domain, val string) (err error) {

	zoneId, err := a.DnsZoneFind(domain)
	if err != nil {
		return
	}

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionDelete,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(domain),
						Type: types.RRTypeTxt,
						TTL:  aws.Int64(int64(settings.Acme.DnsAwsTtl)),
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(val),
							},
						},
					},
				},
			},
			Comment: aws.String("Pritunl delete TXT record"),
		},
		HostedZoneId: aws.String(zoneId),
	}

	_, err = a.client.ChangeResourceRecordSets(context.Background(), input)
	if err != nil {
		if strings.Contains(err.Error(), "delete") &&
			strings.Contains(err.Error(), "not found") {

			err = nil
		} else {
			err = &errortypes.ApiError{
				errors.Wrap(err, "acme: AWS record change error"),
			}
			return
		}
	}

	return
}
