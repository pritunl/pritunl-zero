package authority

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
	"hash/fnv"
	"net"
	"sort"
	"strings"
	"time"
)

type Info struct {
	KeyAlg string `bson:"key_alg" json:"key_alg"`
}

type Authority struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Type       string        `bson:"type" json:"type"`
	Info       *Info         `bson:"info" json:"info"`
	MatchRoles bool          `bson:"match_roles" json:"match_roles"`
	Roles      []string      `bson:"roles" json:"roles"`
	Expire     int           `bson:"expire" json:"expire"`
	PrivateKey string        `bson:"private_key" json:"-"`
	PublicKey  string        `bson:"public_key" json:"public_key"`
	HostDomain string        `bson:"host_domain" json:"host_domain"`
	HostTokens []string      `bson:"host_tokens" json:"host_tokens"`
}

func (a *Authority) GetDomain(hostname string) string {
	return hostname + "." + a.HostDomain
}

func (a *Authority) GenerateRsaPrivateKey() (err error) {
	privKeyBytes, pubKeyBytes, err := GenerateRsaKey()
	if err != nil {
		return
	}

	a.Info = &Info{
		KeyAlg: "RSA 4096",
	}
	a.PrivateKey = strings.TrimSpace(string(privKeyBytes))
	a.PublicKey = strings.TrimSpace(string(pubKeyBytes))

	return
}

func (a *Authority) GenerateEcPrivateKey() (err error) {
	privKeyBytes, pubKeyBytes, err := GenerateEcKey()
	if err != nil {
		return
	}

	a.Info = &Info{
		KeyAlg: "EC P384",
	}
	a.PrivateKey = strings.TrimSpace(string(privKeyBytes))
	a.PublicKey = strings.TrimSpace(string(pubKeyBytes))

	return
}

func (a *Authority) UserHasAccess(usr *user.User) bool {
	if !a.MatchRoles {
		return true
	}
	return usr.RolesMatch(a.Roles)
}

func (a *Authority) HostnameValidate(hostname string, port int,
	pubKey string) bool {

	ips, err := net.LookupIP(a.GetDomain(hostname))
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "authority: Failed to lookup host"),
		}

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("authority: Failed to lookup host")

		return false
	}

	for _, ip := range ips {
		// TODO
		println(ip.String())
	}

	return true
}

func (a *Authority) CreateCertificate(usr *user.User, sshPubKey string) (
	cert *ssh.Certificate, certMarshaled string, err error) {

	privateKey, err := ParsePemKey(a.PrivateKey)
	if err != nil {
		return
	}

	pubKey, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(sshPubKey))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to parse ssh public key"),
		}
		return
	}

	serialHash := fnv.New64a()
	serialHash.Write([]byte(bson.NewObjectId().Hex()))
	serial := serialHash.Sum64()

	validAfter := time.Now().Add(-5 * time.Minute).Unix()
	validBefore := time.Now().Add(
		time.Duration(a.Expire) * time.Minute).Unix()

	cert = &ssh.Certificate{
		Key:             pubKey,
		Serial:          serial,
		CertType:        ssh.UserCert,
		KeyId:           usr.Id.Hex(),
		ValidPrincipals: usr.Roles,
		ValidAfter:      uint64(validAfter),
		ValidBefore:     uint64(validBefore),
		Permissions: ssh.Permissions{
			Extensions: map[string]string{
				"permit-X11-forwarding":   "",
				"permit-agent-forwarding": "",
				"permit-port-forwarding":  "",
				"permit-pty":              "",
				"permit-user-rc":          "",
			},
		},
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		return
	}

	err = cert.SignCert(rand.Reader, signer)
	if err != nil {
		return
	}

	certMarshaled = string(MarshalCertificate(cert, comment))

	return
}

func (a *Authority) CreateHostCertificate(hostname string, sshPubKey string) (
	cert *ssh.Certificate, certMarshaled string, err error) {

	privateKey, err := ParsePemKey(a.PrivateKey)
	if err != nil {
		return
	}

	pubKey, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(sshPubKey))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to parse ssh public key"),
		}
		return
	}

	serialHash := fnv.New64a()
	serialHash.Write([]byte(bson.NewObjectId().Hex()))
	serial := serialHash.Sum64()

	validAfter := time.Now().Add(-5 * time.Minute).Unix()
	validBefore := time.Now().Add(
		time.Duration(a.Expire) * time.Minute).Unix()

	cert = &ssh.Certificate{
		Key:             pubKey,
		Serial:          serial,
		CertType:        ssh.HostCert,
		KeyId:           hostname,
		ValidPrincipals: []string{a.GetDomain(hostname)},
		ValidAfter:      uint64(validAfter),
		ValidBefore:     uint64(validBefore),
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		return
	}

	err = cert.SignCert(rand.Reader, signer)
	if err != nil {
		return
	}

	certMarshaled = string(MarshalCertificate(cert, comment))

	return
}

func (a *Authority) TokenNew() (err error) {
	if a.HostTokens == nil {
		a.HostTokens = []string{}
	}

	token, err := utils.RandStr(32)
	if err != nil {
		return
	}

	a.HostTokens = append(a.HostTokens, token)

	return
}

func (a *Authority) TokenDelete(token string) (err error) {
	if a.HostTokens == nil {
		a.HostTokens = []string{}
	}

	for i, tokn := range a.HostTokens {
		if tokn == token {
			a.HostTokens = append(
				a.HostTokens[:i], a.HostTokens[i+1:]...)
			break
		}
	}

	return
}

func (a *Authority) Export(passphrase string) (encKey string, err error) {
	block, _ := pem.Decode([]byte(a.PrivateKey))

	encBlock, err := x509.EncryptPEMBlock(
		rand.Reader,
		block.Type,
		block.Bytes,
		[]byte(passphrase),
		x509.PEMCipherAES256,
	)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to encrypt private key"),
		}
		return
	}

	encodedBlock := pem.EncodeToMemory(encBlock)

	encKey = string(encodedBlock)

	return
}

func (a *Authority) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if a.Type == "" {
		a.Type = Local
	}

	if !a.MatchRoles {
		a.Roles = []string{}
	}

	if a.PrivateKey == "" {
		err = a.GenerateRsaPrivateKey()
		if err != nil {
			return
		}
	}

	if a.Expire < 1 {
		a.Expire = 600
	} else if a.Expire > 1440 {
		a.Expire = 1440
	}

	if a.HostTokens == nil || a.HostDomain == "" {
		a.HostTokens = []string{}
	}

	a.Format()

	return
}

func (a *Authority) Format() {
	roles := []string{}
	rolesSet := set.NewSet()

	for _, role := range a.Roles {
		rolesSet.Add(role)
	}

	for role := range rolesSet.Iter() {
		roles = append(roles, role.(string))
	}

	sort.Strings(roles)

	a.Roles = roles

	sort.Strings(a.HostTokens)
}

func (a *Authority) Commit(db *database.Database) (err error) {
	coll := db.Authorities()

	err = coll.Commit(a.Id, a)
	if err != nil {
		return
	}

	return
}

func (a *Authority) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Authorities()

	err = coll.CommitFields(a.Id, a, fields)
	if err != nil {
		return
	}

	return
}

func (a *Authority) Insert(db *database.Database) (err error) {
	coll := db.Authorities()

	if a.Id != "" {
		err = &errortypes.DatabaseError{
			errors.New("authority: Authority already exists"),
		}
		return
	}

	err = coll.Insert(a)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
