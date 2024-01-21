package secret

import (
	"crypto/rsa"
	"fmt"

	"github.com/dropbox/godropbox/errors"
	"github.com/oracle/oci-go-sdk/v55/common"
	"github.com/oracle/oci-go-sdk/v55/core"
	"github.com/oracle/oci-go-sdk/v55/dns"
	"github.com/pritunl/pritunl-zero/errortypes"
)

type OracleProvider struct {
	privateKey    *rsa.PrivateKey
	tenancy       string
	user          string
	fingerprint   string
	region        string
	compartment   string
	dnsClient     *dns.DnsClient
	computeClient *core.ComputeClient
}

func (p *OracleProvider) AuthType() (common.AuthConfig, error) {
	return common.AuthConfig{
		AuthType:         common.UserPrincipal,
		IsFromConfigFile: false,
		OboToken:         nil,
	}, nil
}

func (p *OracleProvider) PrivateRSAKey() (*rsa.PrivateKey, error) {
	return p.privateKey, nil
}

func (p *OracleProvider) KeyID() (string, error) {
	return fmt.Sprintf("%s/%s/%s", p.tenancy, p.user, p.fingerprint), nil
}

func (p *OracleProvider) TenancyOCID() (string, error) {
	return p.tenancy, nil
}

func (p *OracleProvider) UserOCID() (string, error) {
	return p.user, nil
}

func (p *OracleProvider) KeyFingerprint() (string, error) {
	return p.fingerprint, nil
}

func (p *OracleProvider) Region() (string, error) {
	return p.region, nil
}

func (p *OracleProvider) CompartmentOCID() (string, error) {
	return p.compartment, nil
}

func (p *OracleProvider) GetDnsClient() (
	dnsClient *dns.DnsClient, err error) {

	if p.dnsClient != nil {
		dnsClient = p.dnsClient
		return
	}

	client, err := dns.NewDnsClientWithConfigurationProvider(p)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "secret: Failed to create oracle client"),
		}
		return
	}

	p.dnsClient = &client
	dnsClient = p.dnsClient

	return
}

func NewOracleProvider(secr *Secret) (prov *OracleProvider, err error) {
	privateKey, fingerprint, err := loadPrivateKey(secr)
	if err != nil {
		return
	}

	prov = &OracleProvider{
		privateKey:  privateKey,
		tenancy:     secr.Key,
		user:        secr.Value,
		fingerprint: fingerprint,
		region:      secr.Region,
		compartment: secr.Key,
	}

	return
}
