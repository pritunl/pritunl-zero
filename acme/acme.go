package acme

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/ericchiang/letsencrypt"
	"github.com/pritunl/pritunl-zero/certificate"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"golang.org/x/crypto/acme"
	"strings"
	"time"
)

func Get(db *database.Database, cert *certificate.Certificate) (err error) {
	if cert.AcmeDomains == nil || len(cert.AcmeDomains) == 0 {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "acme: No acme domains"),
		}
		return
	}

	cli, err := letsencrypt.NewClient(settings.Acme.Url + "/directory")
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "acme: Failed to create acme client"),
		}
		return
	}

	var acctKey *rsa.PrivateKey

	if cert.AcmeAccount != "" {
		acctBlock, _ := pem.Decode([]byte(cert.AcmeAccount))
		if acctBlock == nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "acme: Failed to read account key"),
			}
			return
		}

		acctKey, err = x509.ParsePKCS1PrivateKey(acctBlock.Bytes)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "acme: Failed to parse account key"),
			}
			return
		}
	} else {
		acctKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			err = &errortypes.ReadError{
				errors.Wrap(err, "acme: Failed to generate account key"),
			}
			return
		}

		acctBlock := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(acctKey),
		}

		cert.AcmeAccount = string(pem.EncodeToMemory(acctBlock))
		err = cert.CommitFields(db, set.NewSet("acme_account"))
		if err != nil {
			return
		}
	}

	_, err = cli.NewRegistration(acctKey)
	if err != nil {
		switch err.(type) {
		case *acme.Error:
			acmeErr := err.(*acme.Error)
			if acmeErr.StatusCode == 409 {
				err = nil
			}
			break
		}

		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrap(err, "acme: Failed to create registration"),
			}
			return
		}
	}

	for _, domain := range cert.AcmeDomains {
		auth, _, e := cli.NewAuthorization(acctKey, "dns", domain)
		if e != nil {
			err = &errortypes.RequestError{
				errors.Wrapf(e, "acme: Failed to authorize %s", domain),
			}
			return
		}

		challenges := auth.Combinations(letsencrypt.ChallengeHTTP)
		if len(challenges) == 0 || len(challenges[0]) == 0 {
			err = &errortypes.ParseError{
				errors.Wrap(err, "acme: No supported challenges"),
			}
			return
		}

		challenge := challenges[0][0]

		path, resource, e := challenge.HTTP(acctKey)
		if e != nil {
			err = &errortypes.ParseError{
				errors.Wrap(e, "acme: Failed to generate challenge path"),
			}
			return
		}

		token := ParsePath(path)
		if token == "" {
			err = &errortypes.ParseError{
				errors.Wrap(err, "acme: Failed to parse challenge path"),
			}
			return
		}

		chal := &Challenge{
			Id:        token,
			Resource:  resource,
			Timestamp: time.Now(),
		}

		err = chal.Insert(db)
		if err != nil {
			return
		}

		err = cli.ChallengeReady(acctKey, challenge)
		if err != nil {
			err = &errortypes.RequestError{
				errors.Wrapf(e, "acme: Failed to challenge %s", domain),
			}
			return
		}

		err = chal.Remove(db)
		if err != nil {
			return
		}
	}

	csr, certKey, err := newCsr(cert.AcmeDomains)
	if err != nil {
		return
	}

	certResp, err := cli.NewCertificate(acctKey, csr)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "acme: Failed to get certificate"),
		}
		return
	}

	certKeyByte, err := x509.MarshalECPrivateKey(certKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "acme: Failed to parse private key"),
		}
		return
	}

	certKeyBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: certKeyByte,
	}

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certResp.Certificate.Raw,
	}

	certPem := string(pem.EncodeToMemory(certBlock))
	certPem = strings.Trim(certPem, "\n")
	certPem += AcmeChain

	cert.Key = string(pem.EncodeToMemory(certKeyBlock))
	cert.Certificate = certPem
	err = cert.CommitFields(db, set.NewSet("key", "certificate"))
	if err != nil {
		return
	}

	err = cert.Write()
	if err != nil {
		return
	}

	return
}
