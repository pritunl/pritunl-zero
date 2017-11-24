package authority

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
)

func MarshalCertificate(cert *ssh.Certificate, comment string) []byte {
	b := &bytes.Buffer{}
	b.WriteString(cert.Type())
	b.WriteByte(' ')
	e := base64.NewEncoder(base64.StdEncoding, b)
	e.Write(cert.Marshal())
	e.Close()
	b.WriteByte(' ')
	b.Write([]byte(comment))
	return b.Bytes()
}

func MarshalAuthorizedKey(key ssh.PublicKey) []byte {
	b := &bytes.Buffer{}
	b.WriteString(key.Type())
	b.WriteByte(' ')
	e := base64.NewEncoder(base64.StdEncoding, b)
	e.Write(key.Marshal())
	e.Close()
	return b.Bytes()
}

func GenerateRsaKey() (encodedPriv, encodedPub []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "authority: Failed to generate rsa key"),
		}
		return
	}

	pubKey, err := ssh.NewPublicKey(privateKey.Public())
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to parse rsa key"),
		}
		return
	}

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	encodedPriv = pem.EncodeToMemory(block)
	encodedPub = MarshalAuthorizedKey(pubKey)

	return
}

func GenerateEcKey() (encodedPriv, encodedPub []byte, err error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "authority: Failed to generate ec key"),
		}
		return
	}

	pubKey, err := ssh.NewPublicKey(privateKey.Public())
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to parse ec key"),
		}
		return
	}

	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to marshal ec key"),
		}
		return
	}

	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	}

	encodedPriv = pem.EncodeToMemory(block)
	encodedPub = MarshalAuthorizedKey(pubKey)

	return
}

func ParsePemKey(data string) (key crypto.PrivateKey, err error) {
	block, _ := pem.Decode([]byte(data))

	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "authority: Failed to parse rsa key"),
			}
			return
		}
		break
	case "EC PRIVATE KEY":
		key, err = x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "authority: Failed to parse ec key"),
			}
			return
		}
		break
	default:
		err = &errortypes.ParseError{
			errors.Newf("authority: Unknown key type '%s'", block.Type),
		}
		return
	}

	return
}

func Get(db *database.Database, authrId bson.ObjectId) (
	authr *Authority, err error) {

	coll := db.Authorities()
	authr = &Authority{}

	err = coll.FindOneId(authrId, authr)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database) (authrs []*Authority, err error) {
	coll := db.Authorities()
	authrs = []*Authority{}

	cursor := coll.Find(bson.M{}).Iter()

	authr := &Authority{}
	for cursor.Next(authr) {
		authrs = append(authrs, authr)
		authr = &Authority{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, authrId bson.ObjectId) (err error) {
	coll := db.Authorities()

	_, err = coll.RemoveAll(&bson.M{
		"_id": authrId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
