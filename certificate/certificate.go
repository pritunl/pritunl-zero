package certificate

import (
	"crypto/md5"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
)

type Info struct {
	Hash         string    `bson:"hash" json:"hash"`
	SignatureAlg string    `bson:"signature_alg" json:"signature_alg"`
	PublicKeyAlg string    `bson:"public_key_alg" json:"public_key_alg"`
	Issuer       string    `bson:"issuer" json:"issuer"`
	IssuedOn     time.Time `bson:"issued_on" json:"issued_on"`
	ExpiresOn    time.Time `bson:"expires_on" json:"expires_on"`
	DnsNames     []string  `bson:"dns_names" json:"dns_names"`
}

type Certificate struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Type        string             `bson:"type" json:"type"`
	Key         string             `bson:"key" json:"key"`
	Certificate string             `bson:"certificate" json:"certificate"`
	Info        *Info              `bson:"info" json:"info"`
	AcmeHash    string             `bson:"acme_hash" json:"-"`
	AcmeAccount string             `bson:"acme_account" json:"-"`
	AcmeDomains []string           `bson:"acme_domains" json:"acme_domains"`
}

func (c *Certificate) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if c.Type == "" {
		c.Type = Text
	}

	if c.Type != LetsEncrypt {
		c.AcmeAccount = ""
		c.AcmeDomains = []string{}
	}

	if c.AcmeDomains == nil {
		c.AcmeDomains = []string{}
	}

	for i, domain := range c.AcmeDomains {
		if strings.HasSuffix(domain, ".") {
			c.AcmeDomains[i] = domain[:len(domain)-1]
		}
	}

	if c.Type == LetsEncrypt && len(c.AcmeDomains) == 0 {
		errData = &errortypes.ErrorData{
			Error:   "missing_acme_domains",
			Message: "Lets Encrypt domains required",
		}
		return
	}

	err = c.UpdateInfo()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("certificate: Failed to update certificate info")
		err = nil
	}

	return
}

func (c *Certificate) UpdateInfo() (err error) {
	hash := c.Hash()

	if c.Certificate == "" {
		c.Info = &Info{
			DnsNames: []string{},
		}
		return
	}

	if c.Info != nil && hash == c.Info.Hash {
		return
	}

	certBlock, _ := pem.Decode([]byte(c.Certificate))
	if certBlock == nil {
		c.Info = nil
		err = &errortypes.ParseError{
			errors.New("certificate: Failed to decode certificate"),
		}
		return
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		c.Info = nil
		err = &errortypes.ParseError{
			errors.Wrap(err, "certificate: Failed to parse certificate"),
		}
		return
	}

	publicKeyAlg := ""
	switch cert.PublicKeyAlgorithm {
	case x509.RSA:
		publicKeyAlg = "RSA"
		break
	case x509.DSA:
		publicKeyAlg = "DSA"
		break
	case x509.ECDSA:
		publicKeyAlg = "ECDSA"
		break
	default:
		publicKeyAlg = "Unknown"
	}

	dnsNames := cert.DNSNames
	if len(dnsNames) == 0 && cert.Subject.CommonName != "" {
		dnsNames = append(dnsNames, cert.Subject.CommonName)
	}

	c.Info = &Info{
		Hash:         hash,
		SignatureAlg: cert.SignatureAlgorithm.String(),
		PublicKeyAlg: publicKeyAlg,
		Issuer:       cert.Issuer.CommonName,
		IssuedOn:     cert.NotBefore,
		ExpiresOn:    cert.NotAfter,
		DnsNames:     dnsNames,
	}

	return
}

func (c *Certificate) Commit(db *database.Database) (err error) {
	coll := db.Certificates()

	err = coll.Commit(c.Id, c)
	if err != nil {
		return
	}

	return
}

func (c *Certificate) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Certificates()

	err = coll.CommitFields(c.Id, c, fields)
	if err != nil {
		return
	}

	return
}

func (c *Certificate) Insert(db *database.Database) (err error) {
	coll := db.Certificates()

	if !c.Id.IsZero() {
		err = &errortypes.DatabaseError{
			errors.New("certificate: Certificate already exists"),
		}
		return
	}

	_, err = coll.InsertOne(db, c)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (c *Certificate) Hash() string {
	hash := md5.New()
	hash.Write([]byte(c.Type))
	hash.Write([]byte(c.Key))
	hash.Write([]byte(c.Certificate))
	hash.Write([]byte(c.AcmeAccount))
	if c.AcmeDomains != nil {
		for _, domain := range c.AcmeDomains {
			io.WriteString(hash, domain)
		}
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}
