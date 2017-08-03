package router

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

func generateCert(certPath, keyPath string) (err error) {
	logrus.Info("router: Generating self signed certificate")

	certKey, err := ecdsa.GenerateKey(
		elliptic.P384(),
		rand.Reader,
	)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "router: Failed to generate private key"),
		}
		return
	}

	certKeyByte, err := x509.MarshalECPrivateKey(certKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "router: Failed to parse private key"),
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

	serialLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serial, err := rand.Int(rand.Reader, serialLimit)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "router: Failed to generate certificate serial"),
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

	cert, err := x509.CreateCertificate(rand.Reader, certTempl, certTempl,
		certKey.Public(), certKey)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "router: Failed to create certificate"),
		}
		return
	}

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}

	certPem := string(pem.EncodeToMemory(certBlock))

	err = utils.CreateWrite(certPath, certPem, 0644)
	if err != nil {
		return
	}

	return
}
