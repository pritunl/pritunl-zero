package acme

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/certificate"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme"
)

func Generate(db *database.Database, cert *certificate.Certificate) (
	err error) {

	logrus.WithFields(logrus.Fields{
		"certificate": cert.Name,
		"domains":     cert.AcmeDomains,
	}).Info("acme: Generating acme certificate")

	if cert.AcmeDomains == nil || len(cert.AcmeDomains) == 0 {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "acme: No acme domains"),
		}
		return
	}

	var acctKey *rsa.PrivateKey

	if cert.AcmeAccount != "" {
		acctBlock, _ := pem.Decode([]byte(cert.AcmeAccount))
		if acctBlock == nil {
			err = &errortypes.ParseError{
				errors.New("acme: Failed to decode account key"),
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

	acct := &acme.Account{}

	client := &acme.Client{
		DirectoryURL: AcmeDirectory,
		Key:          acctKey,
	}

	_, err = client.Register(context.Background(), acct, prompt)
	if err != nil {
		if err == acme.ErrAccountAlreadyExists {
			err = nil
		} else {
			err = &errortypes.RequestError{
				errors.Wrap(err, "acme: Failed to register account"),
			}
			return
		}
	}

	order, err := client.AuthorizeOrder(
		context.Background(), acme.DomainIDs(cert.AcmeDomains...))
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "acme: Failed to authorize order"),
		}
		return
	}

	if order.Status == acme.StatusReady {
		err = create(db, cert, client, order)
		if err != nil {
			return
		}

		return
	} else if order.Status != acme.StatusPending {
		err = &errortypes.RequestError{
			errors.Newf(
				"acme: Authorize order status '%s' not pending",
				order.Status,
			),
		}
		return
	}

	authzUrls := order.AuthzURLs

	for _, authzUrl := range authzUrls {
		authz, e := client.GetAuthorization(
			context.Background(), authzUrl)
		if e != nil {
			err = &errortypes.RequestError{
				errors.Wrap(e, "acme: Failed to get authorization"),
			}
			return
		}

		if authz.Status != acme.StatusPending {
			continue
		}

		var authzChal *acme.Challenge
		for _, c := range authz.Challenges {
			if c.Type == "http-01" {
				authzChal = c
				break
			}
		}

		if authzChal == nil {
			revoke(client, authzUrls)

			err = &errortypes.RequestError{
				errors.New(
					"acme: Authorization HTTP challenge not available"),
			}
			return
		}

		resp, e := client.HTTP01ChallengeResponse(authzChal.Token)
		if e != nil {
			revoke(client, authzUrls)

			err = &errortypes.RequestError{
				errors.Wrap(e, "acme: Challenge response failed"),
			}
			return
		}

		chal := &Challenge{
			Id:        authzChal.Token,
			Resource:  resp,
			Timestamp: time.Now(),
		}

		err = chal.Insert(db)
		if err != nil {
			return
		}

		_, err = client.Accept(context.Background(), authzChal)
		if err != nil {
			revoke(client, authzUrls)

			err = &errortypes.RequestError{
				errors.Wrap(err, "acme: Authorization accept failed"),
			}
			return
		}

		_, err = client.WaitAuthorization(
			context.Background(), authzChal.URI)
		if err != nil {
			revoke(client, authzUrls)

			err = &errortypes.RequestError{
				errors.Wrap(err, "acme: Authorization wait failed"),
			}
			return
		}

		err = chal.Remove(db)
		if err != nil {
			revoke(client, authzUrls)

			return
		}
	}

	order, err = client.WaitOrder(context.Background(), order.URI)
	if err != nil {
		revoke(client, authzUrls)

		err = &errortypes.RequestError{
			errors.Wrap(err, "acme: Order wait failed"),
		}
		return
	}

	if order.Status != acme.StatusReady {
		err = &errortypes.RequestError{
			errors.Newf(
				"acme: Authorize order status '%s' not ready",
				order.Status,
			),
		}
		return
	}

	err = create(db, cert, client, order)
	if err != nil {
		return
	}

	return
}

func create(db *database.Database, cert *certificate.Certificate,
	client *acme.Client, order *acme.Order) (err error) {

	var csr []byte
	var keyPem []byte

	if settings.System.AcmeKeyAlgorithm == "ec" {
		csr, keyPem, err = newEcCsr(cert.AcmeDomains)
		if err != nil {
			return
		}
	} else {
		csr, keyPem, err = newRsaCsr(cert.AcmeDomains)
		if err != nil {
			return
		}
	}

	derChain, _, err := client.CreateOrderCert(
		context.Background(),
		order.FinalizeURL,
		csr,
		true,
	)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "acme: Create order cert failed"),
		}
		return
	}

	certPem := ""

	for _, der := range derChain {
		certBlock := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: der,
		}

		if certPem != "" {
			certPem += "\n"
		}
		certPem += strings.TrimSpace(string(pem.EncodeToMemory(certBlock)))
	}

	cert.Key = strings.TrimSpace(string(keyPem))
	cert.Certificate = certPem
	cert.AcmeHash = cert.Hash()

	_, err = cert.Validate(db)
	if err != nil {
		return
	}

	err = cert.CommitFields(db, set.NewSet(
		"key", "certificate", "acme_hash", "info"))
	if err != nil {
		return
	}

	return
}

func Update(db *database.Database, cert *certificate.Certificate) (err error) {
	if cert.Type != certificate.LetsEncrypt {
		return
	}

	if cert.AcmeHash != cert.Hash() {
		err = Generate(db, cert)
		if err != nil {
			return
		}
	}

	return
}

func Renew(db *database.Database, cert *certificate.Certificate) (
	err error) {

	if cert.Type != certificate.LetsEncrypt {
		return
	}

	if cert.Info != nil && !cert.Info.ExpiresOn.IsZero() &&
		time.Until(cert.Info.ExpiresOn) < 168*time.Hour {

		err = Generate(db, cert)
		if err != nil {
			return
		}
	}

	return
}
