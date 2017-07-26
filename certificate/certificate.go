package certificate

import (
	"crypto/md5"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"io"
	"time"
)

type Info struct {
	Hash         string    `bson:"hash" json:"hash"`
	SignatureAlg string    `bson:"signature_alg" json:"signature_alg"`
	PublicKeyAlg string    `bson:"public_key_alg" json:"public_key_alg"`
	IssuedOn     time.Time `bson:"issued_on" json:"issued_on"`
	ExpiresOn    time.Time `bson:"expires_on" json:"expires_on"`
	DnsNames     []string  `bson:"dns_names" json:"dns_names"`
}

type Certificate struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Type        string        `bson:"type" json:"type"`
	Key         string        `bson:"key" json:"key"`
	Certificate string        `bson:"certificate" json:"certificate"`
	Info        *Info         `bson:"info" json:"info"`
	AcmeHash    string        `bson:"acme_hash" json:"acme_hash"`
	AcmeAccount string        `bson:"acme_account" json:"acme_account"`
	AcmeDomains []string      `bson:"acme_domains" json:"acme_domains"`
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

	return
}

func (c *Certificate) UpdateInfo() (err error) {
	hash := c.Hash()

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

	info := &Info{
		Hash:         hash,
		SignatureAlg: cert.SignatureAlgorithm.String(),
		PublicKeyAlg: publicKeyAlg,
		IssuedOn:     cert.NotBefore,
		ExpiresOn:    cert.NotAfter,
		DnsNames:     dnsNames,
	}
	c.Info = info

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

	if c.Id != "" {
		err = &errortypes.DatabaseError{
			errors.New("certificate: Certificate already exists"),
		}
		return
	}

	err = coll.Insert(c)
	if err != nil {
		return
	}

	return
}

func (c *Certificate) Hash() string {
	hash := md5.New()
	io.WriteString(hash, c.Type)
	io.WriteString(hash, c.Key)
	io.WriteString(hash, c.Certificate)
	io.WriteString(hash, c.AcmeAccount)
	if c.AcmeDomains != nil {
		for _, domain := range c.AcmeDomains {
			io.WriteString(hash, domain)
		}
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (c *Certificate) Write() (err error) {
	err = utils.CreateWrite(constants.KeyPath, c.Key, 0600)
	if err != nil {
		return
	}

	err = utils.CreateWrite(constants.CertPath, c.Certificate, 0666)
	if err != nil {
		return
	}

	return
}
