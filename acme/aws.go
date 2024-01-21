package acme

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/secret"
	"github.com/pritunl/pritunl-zero/settings"
	"strings"
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
		if strings.Trim(*zone.Name, ".") == strings.Trim(domain, ".") {
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
								Value: aws.String("\"" + val + "\""),
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
								Value: aws.String("\"" + val + "\""),
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
		err = &errortypes.ApiError{
			errors.Wrap(err, "acme: AWS route53 record set error"),
		}
		return
	}

	return
}
