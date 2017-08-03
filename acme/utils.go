package acme

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"strings"
)

func ParsePath(path string) string {
	split := strings.SplitN(path, AcmePath, 2)
	if len(split) == 2 {
		return split[1]
	}
	return ""
}

func GetChallenge(token string) (challenge *Challenge, err error) {
	db := database.GetDatabase()
	defer db.Close()

	coll := db.AcmeChallenges()
	challenge = &Challenge{}

	err = coll.FindOneId(token, challenge)
	if err != nil {
		return
	}

	return
}

func newCsr(domains []string) (csr *x509.CertificateRequest,
	key *ecdsa.PrivateKey, err error) {

	key, err = ecdsa.GenerateKey(
		elliptic.P384(),
		rand.Reader,
	)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "acme: Failed to generate private key"),
		}
		return
	}

	csrReq := &x509.CertificateRequest{
		SignatureAlgorithm: x509.ECDSAWithSHA256,
		PublicKeyAlgorithm: x509.ECDSA,
		PublicKey:          key.Public(),
		Subject: pkix.Name{
			CommonName: domains[0],
		},
		DNSNames: domains,
	}

	csrData, err := x509.CreateCertificateRequest(rand.Reader, csrReq, key)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "acme: Failed to create certificate request"),
		}
		return
	}

	csr, err = x509.ParseCertificateRequest(csrData)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "acme: Failed to parse certificate request"),
		}
		return
	}

	return
}
