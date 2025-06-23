package dns

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/secret"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/sirupsen/logrus"
)

type Aws struct {
	sess        *session.Session
	sessRoute53 *route53.Route53
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

	a.sess, err = session.NewSession(&aws.Config{
		Region: aws.String(secr.Region),
		Credentials: credentials.NewStaticCredentials(
			secr.Key, secr.Value, ""),
	})
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: AWS session error"),
		}
		return
	}

	a.sessRoute53 = route53.New(a.sess)

	return
}

func (a *Aws) DnsZoneFind(domain string) (zoneId string, err error) {
	domain = extractDomain(domain)

	zoneId = a.cacheZoneId[domain]
	if zoneId != "" {
		return
	}

	input := &route53.ListHostedZonesInput{}

	result, err := a.sessRoute53.ListHostedZones(input)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: AWS route53 zone lookup error"),
		}
		return
	}

	for _, zone := range result.HostedZones {
		if matchDomains(*zone.Name, domain) {
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

	action := aws.String("DELETE")
	resourceRecs := []*route53.ResourceRecord{}
	values := []string{}
	for _, op := range ops {
		if op.Operation == UPSERT || op.Operation == RETAIN {
			action = aws.String("UPSERT")
		}

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

		resourceRecs = append(resourceRecs, &route53.ResourceRecord{
			Value: aws.String(op.Value),
		})
		values = append(values, op.Value)
	}

	logrus.WithFields(logrus.Fields{
		"operation": *action,
		"domain":    domain,
		"values":    values,
	}).Info("domain: AWS dns batch operation")

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: action,
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(domain),
						Type: aws.String(recordType),
						TTL: aws.Int64(int64(
							settings.Acme.DnsAwsTtl)),
						ResourceRecords: resourceRecs,
					},
				},
			},
			Comment: aws.String("Pritunl update record"),
		},
		HostedZoneId: aws.String(zoneId),
	}

	_, err = a.sessRoute53.ChangeResourceRecordSets(input)
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
		StartRecordType: aws.String(recordType),
	}

	result, err := a.sessRoute53.ListResourceRecordSets(input)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: AWS record list error"),
		}
		return
	}

	for _, recordSet := range result.ResourceRecordSets {
		if recordSet.Type != nil && *recordSet.Type == recordType &&
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
		StartRecordType: aws.String("TXT"),
	}

	result, err := a.sessRoute53.ListResourceRecordSets(input)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: AWS route53 record set error"),
		}
		return
	}

	for _, recordSet := range result.ResourceRecordSets {
		if recordSet.Type != nil && *recordSet.Type == "TXT" &&
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
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(domain),
						Type: aws.String("TXT"),
						TTL:  aws.Int64(int64(settings.Acme.DnsAwsTtl)),
						ResourceRecords: []*route53.ResourceRecord{
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

	_, err = a.sessRoute53.ChangeResourceRecordSets(input)
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
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("DELETE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(domain),
						Type: aws.String("TXT"),
						TTL:  aws.Int64(int64(settings.Acme.DnsAwsTtl)),
						ResourceRecords: []*route53.ResourceRecord{
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

	_, err = a.sessRoute53.ChangeResourceRecordSets(input)
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
