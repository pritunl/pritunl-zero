package authority

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"strings"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/mongo-go-driver/mongo/options"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/utils"
	"golang.org/x/crypto/ssh"
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

func MarshalCertificate(cert *ssh.Certificate, comment string) (
	data []byte, err error) {

	buffer := &bytes.Buffer{}
	buffer.WriteString(cert.Type())
	buffer.WriteByte(' ')

	encoder := base64.NewEncoder(base64.StdEncoding, buffer)

	_, err = encoder.Write(cert.Marshal())
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to write certificate"),
		}
		return
	}

	err = encoder.Close()
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to close certificate"),
		}
		return
	}

	if comment != "" {
		buffer.WriteByte(' ')
		buffer.Write([]byte(comment))
	}

	data = buffer.Bytes()
	return
}

func MarshalPublicKey(key ssh.PublicKey) (data []byte, err error) {
	buffer := &bytes.Buffer{}
	buffer.WriteString(key.Type())
	buffer.WriteByte(' ')

	encoder := base64.NewEncoder(base64.StdEncoding, buffer)

	_, err = encoder.Write(key.Marshal())
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to write public key"),
		}
		return
	}

	err = encoder.Close()
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to close public key"),
		}
		return
	}

	data = buffer.Bytes()
	return
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

	encodedPub, err = MarshalPublicKey(pubKey)
	if err != nil {
		return
	}

	return
}

func GenerateEdKey() (encodedPriv, encodedPub []byte, err error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "authority: Failed to generate ed key"),
		}
		return
	}

	pubKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to parse ed key"),
		}
		return
	}

	block := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: utils.MarshalED25519PrivateKey(privateKey),
	}

	encodedPriv = pem.EncodeToMemory(block)

	encodedPub, err = MarshalPublicKey(pubKey)
	if err != nil {
		return
	}

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

	encodedPub, err = MarshalPublicKey(pubKey)
	if err != nil {
		return
	}

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

func ParseSshPubKey(data string) (pubKey crypto.PublicKey, err error) {
	sshPubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(data))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to parse ssh public key"),
		}
		return
	}

	cryptoPubKey, ok := sshPubKey.(ssh.CryptoPublicKey)
	if !ok {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to parse ssh public key type"),
		}
		return
	}

	pubKey = cryptoPubKey.CryptoPublicKey()

	return
}

func Get(db *database.Database, authrId primitive.ObjectID) (
	authr *Authority, err error) {

	coll := db.Authorities()
	authr = &Authority{}

	err = coll.FindOneId(authrId, authr)
	if err != nil {
		return
	}

	return
}

func GetOne(db *database.Database, query *bson.M) (
	authr *Authority, err error) {

	coll := db.Authorities()
	authr = &Authority{}

	err = coll.FindOne(db, query).Decode(authr)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetHsmToken(db *database.Database, token string) (
	authr *Authority, err error) {

	coll := db.Authorities()
	authr = &Authority{}

	err = coll.FindOne(db, &bson.M{
		"hsm_token": token,
	}).Decode(authr)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetMulti(db *database.Database, authrIds []primitive.ObjectID) (
	authrs []*Authority, err error) {

	coll := db.Authorities()
	authrs = []*Authority{}

	cursor, err := coll.Find(db, &bson.M{
		"_id": &bson.M{"$in": authrIds},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		authr := &Authority{}
		err = cursor.Decode(authr)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		authrs = append(authrs, authr)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAll(db *database.Database) (authrs []*Authority, err error) {
	coll := db.Authorities()
	authrs = []*Authority{}

	cursor, err := coll.Find(db, &bson.M{})
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		authr := &Authority{}
		err = cursor.Decode(authr)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		authrs = append(authrs, authr)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAllPaged(db *database.Database, query *bson.M,
	page, pageCount int64) (authorities []*Authority, count int64, err error) {

	coll := db.Authorities()
	authorities = []*Authority{}

	if len(*query) == 0 {
		count, err = coll.EstimatedDocumentCount(db)
		if err != nil {
			err = database.ParseError(err)
			return
		}
	} else {
		count, err = coll.CountDocuments(db, query)
		if err != nil {
			err = database.ParseError(err)
			return
		}
	}

	maxPage := count / pageCount
	if count == pageCount {
		maxPage = 0
	}
	page = utils.Min64(page, maxPage)
	skip := utils.Min64(page*pageCount, count)

	cursor, err := coll.Find(
		db,
		query,
		&options.FindOptions{
			Sort: &bson.D{
				{"name", 1},
			},
			Skip:  &skip,
			Limit: &pageCount,
		},
	)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		authority := &Authority{}
		err = cursor.Decode(authority)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		authorities = append(authorities, authority)
	}

	err = cursor.Err()
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

	cursor, err := coll.Find(db, &bson.M{
		"host_tokens": &bson.M{
			"$in": tokens,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		authr := &Authority{}
		err = cursor.Decode(authr)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		authrs = append(authrs, authr)
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, authrId primitive.ObjectID) (
	errData *errortypes.ErrorData, err error) {

	err = RemoveNode(db, authrId)
	if err != nil {
		return
	}

	coll := db.Services()

	count, err := coll.CountDocuments(db, &bson.M{
		"client_authority": authrId,
	})
	if err != nil {
		err = database.ParseError(err)
		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		default:
			return
		}
	}

	if count > 0 {
		errData = &errortypes.ErrorData{
			Error:   "authority_in_use_service",
			Message: "Cannot delete authority in use by service",
		}
		return
	}

	coll = db.Authorities()

	_, err = coll.DeleteOne(db, &bson.M{
		"_id": authrId,
	})
	if err != nil {
		err = database.ParseError(err)
		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		default:
			return
		}
	}

	return
}

func RemoveMulti(db *database.Database, authorityIds []primitive.ObjectID) (
	err error) {
	coll := db.Authorities()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": &bson.M{
			"$in": authorityIds,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func RemoveNode(db *database.Database,
	authrId primitive.ObjectID) (err error) {

	coll := db.Nodes()

	_, err = coll.UpdateMany(db, &bson.M{
		"authorities": authrId,
	}, &bson.M{
		"$pull": &bson.M{
			"authorities": authrId,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		default:
			return
		}
	}

	return
}
