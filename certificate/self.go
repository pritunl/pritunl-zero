package certificate

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/utils"
	"math/big"
	"time"
)

var selfCert = false

func selfGenerateCert(parent *x509.Certificate, parentKey *ecdsa.PrivateKey) (
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
			errors.Wrap(err, "certificate: Failed to generate certificate serial"),
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

func SelfGenerateCert(certPath, keyPath string) (err error) {
	if selfCert {
		return
	}

	logrus.Info("certificate: Generating self signed certificate")

	caCert, _, caKey, err := selfGenerateCert(nil, nil)
	if err != nil {
		return
	}

	_, certByt, certKey, err := selfGenerateCert(caCert, caKey)
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

	certKeyPem := string(pem.EncodeToMemory(certKeyBlock))

	err = utils.CreateWrite(keyPath, certKeyPem, 0600)
	if err != nil {
		return
	}

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certByt,
	}

	certPem := string(pem.EncodeToMemory(certBlock))

	err = utils.CreateWrite(certPath, certPem, 0644)
	if err != nil {
		return
	}

	selfCert = true

	return
}
