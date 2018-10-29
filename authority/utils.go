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
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
	"net"
	"strings"
)

func parseSubnetMatch(subnetMatch string) (
	match string, err error) {

	_, subnet, err := net.ParseCIDR(subnetMatch)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to parse cidr"),
		}
		return
	}

	cidr, _ := subnet.Mask.Size()

	subnetNet := strings.SplitN(subnet.String(), "/", 2)[0]
	parts := strings.Split(subnetNet, ".")

	if strings.Contains(subnetMatch, ":") {
		if !strings.HasSuffix(subnetNet, "::") {
			err = &errortypes.ParseError{
				errors.New("authority: IPv6 subnet suffix invalid"),
			}
			return
		}

		if len(subnetNet) < 6 {
			err = &errortypes.ParseError{
				errors.New("authority: IPv6 subnet length invalid"),
			}
			return
		}

		switch cidr {
		case 56:
			match = fmt.Sprintf(
				"%s*",
				subnetNet[:len(subnetNet)-4],
			)
			break
		case 64:
			match = fmt.Sprintf(
				"%s*",
				subnetNet[:len(subnetNet)-2],
			)
			break
		default:
			err = &errortypes.ParseError{
				errors.New("authority: Unsupported subnet size"),
			}
			return
		}
	} else {
		if len(parts) != 4 {
			err = &errortypes.ParseError{
				errors.New("authority: Failed to split subnet parts"),
			}
			return
		}

		switch cidr {
		case 8:
			match = fmt.Sprintf(
				"%s.*.*.*",
				parts[0],
			)
			break
		case 16:
			match = fmt.Sprintf(
				"%s.%s.*.*",
				parts[0],
				parts[1],
			)
			break
		case 24:
			match = fmt.Sprintf(
				"%s.%s.%s.*",
				parts[0],
				parts[1],
				parts[2],
			)
			break
		case 32:
			match = fmt.Sprintf(
				"%s.%s.%s.%s",
				parts[0],
				parts[1],
				parts[2],
				parts[3],
			)
			break
		default:
			err = &errortypes.ParseError{
				errors.New("authority: Unsupported subnet size"),
			}
			return
		}
	}

	return
}

func MarshalCertificate(cert *ssh.Certificate, comment string) []byte {
	b := &bytes.Buffer{}
	b.WriteString(cert.Type())
	b.WriteByte(' ')
	e := base64.NewEncoder(base64.StdEncoding, b)
	e.Write(cert.Marshal())
	e.Close()
	if comment != "" {
		b.WriteByte(' ')
		b.Write([]byte(comment))
	}
	return b.Bytes()
}

func MarshalPublicKey(key ssh.PublicKey) []byte {
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
	encodedPub = MarshalPublicKey(pubKey)

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
	encodedPub = MarshalPublicKey(pubKey)

	return
}

func ParsePemKey(data string) (key crypto.PrivateKey, err error) {
	block, _ := pem.Decode([]byte(data))
	if block == nil {
		err = &errortypes.ParseError{
			errors.New("authority: Failed to decode private key"),
		}
		return
	}

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

func GetHsmToken(db *database.Database, token string) (
	authr *Authority, err error) {

	coll := db.Authorities()
	authr = &Authority{}

	err = coll.FindOne(&bson.M{
		"hsm_token": token,
	}, authr)
	if err != nil {
		return
	}

	return
}

func GetMulti(db *database.Database, authrIds []bson.ObjectId) (
	authrs []*Authority, err error) {

	coll := db.Authorities()
	authrs = []*Authority{}

	cursor := coll.Find(&bson.M{
		"_id": &bson.M{"$in": authrIds},
	}).Iter()

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

func GetTokens(db *database.Database, tokens []string) (
	authrs []*Authority, err error) {

	coll := db.Authorities()
	authrs = []*Authority{}

	cursor := coll.Find(&bson.M{
		"host_tokens": &bson.M{
			"$in": tokens,
		},
	}).Iter()

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
