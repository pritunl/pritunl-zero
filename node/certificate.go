package node

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/sirupsen/logrus"
)

func selfCert(parent *x509.Certificate, parentKey *ecdsa.PrivateKey) (
	cert *x509.Certificate, certByt []byte, certKey *ecdsa.PrivateKey,
	err error) {

	certKey, err = ecdsa.GenerateKey(
		elliptic.P384(),
		rand.Reader,
	)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "certificate: Failed to generate private key"),
		}
		return
	}

	serialLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serial, err := rand.Int(rand.Reader, serialLimit)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(
				err,
				"certificate: Failed to generate certificate serial",
			),
		}
		return
	}

	certTempl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"Pritunl Zero"},
		},
		NotBefore: time.Now().Add(-24 * time.Hour),
		NotAfter:  time.Now().Add(26280 * time.Hour),
		KeyUsage: x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
	}

	if parent == nil {
		parent = certTempl
		parentKey = certKey
	}

	certByt, err = x509.CreateCertificate(rand.Reader, certTempl, parent,
		certKey.Public(), parentKey)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "certificate: Failed to create certificate"),
		}
		return
	}

	cert, err = x509.ParseCertificate(certByt)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "certificate: Failed to parse certificate"),
		}
		return
	}

	return
}

func SelfCert() (certPem, keyPem []byte, err error) {
	if Self.SelfCertificate != "" && Self.SelfCertificateKey != "" {
		certPem = []byte(Self.SelfCertificate)
		keyPem = []byte(Self.SelfCertificateKey)
		return
	}

	logrus.Info("certificate: Generating self signed certificate")

	caCert, _, caKey, err := selfCert(nil, nil)
	if err != nil {
		return
	}

	_, certByt, certKey, err := selfCert(caCert, caKey)
	if err != nil {
		return
	}

	certKeyByte, err := x509.MarshalECPrivateKey(certKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "certificate: Failed to parse private key"),
		}
		return
	}

	certKeyBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: certKeyByte,
	}
	keyPem = pem.EncodeToMemory(certKeyBlock)

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certByt,
	}
	certPem = pem.EncodeToMemory(certBlock)

	db := database.GetDatabase()
	defer db.Close()

	Self.SelfCertificate = string(certPem)
	Self.SelfCertificateKey = string(keyPem)
	err = Self.CommitFields(db, set.NewSet(
		"self_certificate", "self_certificate_key"))
	if err != nil {
		return
	}

	return
}
