package authority

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"hash/fnv"
	"math/big"
	"net"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/nonce"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var (
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
)

type validateData struct {
	PublicKey string `bson:"public_key" json:"public_key"`
}

type Info struct {
	KeyAlg string `bson:"key_alg" json:"key_alg"`
}

type Authority struct {
	Id                 bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name               string        `bson:"name" json:"name"`
	Type               string        `bson:"type" json:"type"`
	Info               *Info         `bson:"info" json:"info"`
	MatchRoles         bool          `bson:"match_roles" json:"match_roles"`
	Roles              []string      `bson:"roles" json:"roles"`
	Expire             int           `bson:"expire" json:"expire"`
	HostExpire         int           `bson:"host_expire" json:"host_expire"`
	Algorithm          string        `bson:"algorithm" json:"algorithm"`
	KeyIdFormat        string        `bson:"key_id_format" json:"key_id_format"`
	PrivateKey         string        `bson:"private_key" json:"-"`
	PublicKey          string        `bson:"public_key" json:"public_key"`
	PublicKeyPem       string        `bson:"public_key_pem" json:"public_key_pem"`
	RootCertificate    string        `bson:"root_certificate" json:"root_certificate"`
	ProxyJump          string        `bson:"-" json:"proxy_jump"`
	ProxyPrivateKey    string        `bson:"proxy_private_key" json:"-"`
	ProxyPublicKey     string        `bson:"proxy_public_key" json:"proxy_public_key"`
	ProxyHosting       bool          `bson:"proxy_hosting" json:"proxy_hosting"`
	ProxyHostname      string        `bson:"proxy_hostname" json:"proxy_hostname"`
	ProxyPort          int           `bson:"proxy_port" json:"proxy_port"`
	HostDomain         string        `bson:"host_domain" json:"host_domain"`
	HostSubnets        []string      `bson:"host_subnets" json:"host_subnets"`
	HostMatches        []string      `bson:"host_matches" json:"host_matches"`
	HostProxy          string        `bson:"host_proxy" json:"host_proxy"`
	HostCertificates   bool          `bson:"host_certificates" json:"host_certificates"`
	StrictHostChecking bool          `bson:"strict_host_checking" json:"strict_host_checking"`
	HostTokens         []string      `bson:"host_tokens" json:"host_tokens"`
	HsmToken           string        `bson:"hsm_token" json:"hsm_token"`
	HsmSecret          string        `bson:"hsm_secret" json:"hsm_secret"`
	HsmSerial          string        `bson:"hsm_serial" json:"hsm_serial"`
	HsmStatus          string        `bson:"hsm_status" json:"hsm_status"`
	HsmTimestamp       time.Time     `bson:"hsm_timestamp" json:"hsm_timestamp"`
}

func (a *Authority) GetDomain(hostname string) string {
	return hostname + "." + a.HostDomain
}

func (a *Authority) GenerateProxyKey() (err error) {
	logrus.WithFields(logrus.Fields{
		"authority_id": a.Id.Hex(),
	}).Info("authority: Generating bastion host key")

	if a.Algorithm == ECP384 {
		err = a.generateEcProxyPrivateKey()
		if err != nil {
			return
		}
	} else {
		err = a.generateEdProxyPrivateKey()
		if err != nil {
			return
		}
	}

	return
}

func (a *Authority) generateEdProxyPrivateKey() (err error) {
	privKeyBytes, pubKeyBytes, err := GenerateEdKey()
	if err != nil {
		return
	}

	a.Info = &Info{
		KeyAlg: "Ed25519",
	}
	a.ProxyPrivateKey = strings.TrimSpace(string(privKeyBytes))
	a.ProxyPublicKey = strings.TrimSpace(string(pubKeyBytes))

	return
}

func (a *Authority) generateEcProxyPrivateKey() (err error) {
	privKeyBytes, pubKeyBytes, err := GenerateEcKey()
	if err != nil {
		return
	}

	a.Info = &Info{
		KeyAlg: "EC P384",
	}
	a.ProxyPrivateKey = strings.TrimSpace(string(privKeyBytes))
	a.ProxyPublicKey = strings.TrimSpace(string(pubKeyBytes))

	return
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

	err = a.SetPublicKeyPem()
	if err != nil {
		return
	}

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

	err = a.SetPublicKeyPem()
	if err != nil {
		return
	}

	return
}

func (a *Authority) GeneratePrivateKey() (err error) {
	if a.Algorithm == ECP384 {
		err = a.GenerateEcPrivateKey()
	} else {
		err = a.GenerateRsaPrivateKey()
	}

	return
}

func (a *Authority) GenerateHsmToken() (err error) {
	a.PublicKey = ""

	a.HsmToken, err = utils.RandStr(32)
	if err != nil {
		return
	}

	a.HsmSecret, err = utils.RandStr(64)
	if err != nil {
		return
	}

	return
}

func (a *Authority) GetHostDomain() string {
	if a.HostDomain == "" {
		return ""
	}

	domain := "*." + a.HostDomain

	bastionDomain := a.GetBastionDomain()
	if bastionDomain != "" {
		domain += " !" + bastionDomain
	}

	return domain
}

func (a *Authority) GetBastionDomain() string {
	jumpProxy := a.JumpProxy()
	if jumpProxy == "" {
		return ""
	}

	hostProxy := strings.SplitN(jumpProxy, "@", 2)
	hostProxy = strings.SplitN(hostProxy[len(hostProxy)-1], ":", 2)
	return hostProxy[0]
}

func (a *Authority) GetCertAuthority() string {
	if a.HostDomain == "" {
		return ""
	}
	return fmt.Sprintf("@cert-authority *.%s %s", a.HostDomain, a.PublicKey)
}

func (a *Authority) GetBastionCertAuthority() string {
	bastionDomain := a.GetBastionDomain()
	if bastionDomain == "" {
		return ""
	}

	return fmt.Sprintf("@cert-authority %s %s", bastionDomain, a.PublicKey)
}

func (a *Authority) UserHasAccess(usr *user.User) bool {
	if !a.MatchRoles {
		return true
	}
	return usr.RolesMatch(a.Roles)
}

func (a *Authority) HostnameValidate(hostname string, port int,
	pubKey string) bool {

	domain := a.GetDomain(hostname)

	ipsNet, err := net.LookupIP(domain)
	if err != nil {
		err = &errortypes.RequestError{
			errors.Wrap(err, "authority: Failed to lookup host"),
		}

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("authority: Failed to lookup host")

		return false
	}

	ips := []net.IP{}
	for _, ip := range ipsNet {
		if ip.To4() != nil {
			ips = append(ips, ip)
		}
	}

	if len(ips) == 0 {
		logrus.WithFields(logrus.Fields{
			"host": domain,
		}).Error("authority: No IPv4 addresses found for host")
		return false
	}

	valid := false
	url := ""
	if port == 0 {
		port = 9748
	}

	for _, ip := range ips {
		url = fmt.Sprintf("http://%s:%d/challenge", ip, port)
		req, e := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if e != nil {
			err = &errortypes.RequestError{
				errors.Wrap(e, "authority: Validation request failed"),
			}
			continue
		}

		resp, e := client.Do(req)
		if e != nil {
			err = &errortypes.RequestError{
				errors.Wrap(e, "authority: Validation request failed"),
			}
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			err = &errortypes.RequestError{
				errors.Newf("authority: Validation request bad status %d",
					resp.StatusCode),
			}
			continue
		}

		data := &validateData{}
		e = json.NewDecoder(resp.Body).Decode(data)
		if e != nil {
			err = &errortypes.ParseError{
				errors.Wrap(e, "authority: Failed to parse response"),
			}
			break
		}

		hostPubKey := strings.TrimSpace(data.PublicKey)
		if len(hostPubKey) > settings.System.SshPubKeyLen {
			err = errortypes.ParseError{
				errors.New("authority: Public key too long"),
			}
			break
		}

		if subtle.ConstantTimeCompare([]byte(pubKey),
			[]byte(hostPubKey)) != 1 {

			err = errortypes.AuthenticationError{
				errors.New("authority: Public key does not match"),
			}
			break
		}

		valid = true
		err = nil
		break
	}

	if err != nil || !valid {
		logrus.WithFields(logrus.Fields{
			"host":  domain,
			"url":   url,
			"error": err,
		}).Error("authority: Host validation failed")
		return false
	}

	return true
}

func (a *Authority) createCertificateLocal(
	usr *user.User, keyId, sshPubKey string) (
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
	_, _ = serialHash.Write([]byte(bson.NewObjectID().Hex()))
	serial := serialHash.Sum64()

	expire := a.Expire
	if expire == 0 {
		expire = 600
	}
	validAfter := time.Now().Add(-3 * time.Minute).Unix()
	validBefore := time.Now().Add(
		time.Duration(expire) * time.Minute).Unix()

	if len(usr.Roles) == 0 {
		err = &errortypes.AuthenticationError{
			errors.Wrap(err, "authority: User has no roles"),
		}
		return
	}

	roles := usr.Roles
	if a.JumpProxy() != "" {
		hasBastion := slices.Contains(roles, "bastion")

		if !hasBastion {
			roles = append(usr.Roles, "bastion")
		}
	}

	cert = &ssh.Certificate{
		Key:             pubKey,
		Serial:          serial,
		CertType:        ssh.UserCert,
		KeyId:           keyId,
		ValidPrincipals: roles,
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

	certMarshaledByt, err := MarshalCertificate(cert, comment)
	if err != nil {
		return
	}

	certMarshaled = string(certMarshaledByt)

	return
}

func (a *Authority) createCertificateHsm(db *database.Database,
	usr *user.User, keyId, sshPubKey string) (cert *ssh.Certificate,
	certMarshaled string, err error) {

	pubKey, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(sshPubKey))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to parse ssh public key"),
		}
		return
	}

	expire := a.Expire
	if expire == 0 {
		expire = 600
	}
	validAfter := time.Now().Add(-3 * time.Minute).Unix()
	validBefore := time.Now().Add(
		time.Duration(expire) * time.Minute).Unix()

	if len(usr.Roles) == 0 {
		err = &errortypes.AuthenticationError{
			errors.Wrap(err, "authority: User has no roles"),
		}
		return
	}

	roles := usr.Roles
	if a.JumpProxy() != "" {
		hasBastion := slices.Contains(roles, "bastion")

		if !hasBastion {
			roles = append(usr.Roles, "bastion")
		}
	}

	cert = &ssh.Certificate{
		Key:             pubKey,
		CertType:        ssh.UserCert,
		KeyId:           keyId,
		ValidPrincipals: roles,
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

	certData, err := utils.MarshalSshCertificate(cert)
	if err != nil {
		return
	}

	data := SshRequest{
		Serial:      a.HsmSerial,
		Certificate: certData,
	}

	cipData, err := json.Marshal(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to marshal certificate"),
		}
		return
	}

	pad := 16 - len(cipData)%16
	for i := 0; i < pad; i++ {
		cipData = append(cipData, 0)
	}

	encKeyHash := sha256.New()
	encKeyHash.Write([]byte(a.HsmSecret))
	cipKey := encKeyHash.Sum(nil)

	cipIv, err := utils.RandBytes(aes.BlockSize)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(cipKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to load cipher"),
		}
		return
	}

	mode := cipher.NewCBCEncrypter(block, cipIv)
	mode.CryptBlocks(cipData, cipData)

	if a.HsmSecret == "" {
		err = &errortypes.ReadError{
			errors.Wrap(err, "session: Empty secret"),
		}
		return
	}

	hashFunc := hmac.New(sha512.New, []byte(a.HsmSecret))
	hashFunc.Write(cipData)
	rawSignature := hashFunc.Sum(nil)
	sig := base64.StdEncoding.EncodeToString(rawSignature)

	payloadId := bson.NewObjectID().Hex()
	payload := &HsmPayload{
		Id:        payloadId,
		Token:     a.HsmToken,
		Iv:        cipIv,
		Signature: sig,
		Type:      "ssh_certificate",
		Data:      cipData,
	}

	waiter := sync.WaitGroup{}
	waiter.Add(1)

	timeout := time.Duration(settings.System.HsmResponseTimeout) * time.Second
	start := time.Now()
	var eventErr error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithFields(logrus.Fields{
					"error": errors.New(fmt.Sprintf("%s", r)),
				}).Error("authority: Parse hsm panic")

				eventErr = &errortypes.UnknownError{
					errors.New("authority: Parse hsm panic"),
				}
			}
			waiter.Done()
		}()

		event.SubscribeType([]string{"pritunl_hsm_recv"}, 5*time.Second,
			func() event.CustomEvent {
				return &HsmEvent{}
			},
			func(msgInf event.CustomEvent, e error) bool {
				if e != nil {
					eventErr = e
					return false
				}

				if msgInf == nil || msgInf.GetData() == nil {
					if time.Since(start) < timeout {
						return true
					}

					eventErr = &errortypes.UnknownError{
						errors.New("authority: Timeout waiting for hsm"),
					}
					return false
				}

				msg := msgInf.(*HsmEvent)

				if msg.Data.Id != payloadId ||
					msg.Data.Type != "ssh_certificate" {

					if time.Since(start) < timeout {
						return true
					}
					eventErr = &errortypes.UnknownError{
						errors.New("authority: Timeout waiting for hsm"),
					}
					return false
				}

				payloadData, e := UnmarshalPayload(
					a.HsmToken, a.HsmSecret, msg.Data)
				if e != nil {
					eventErr = e
					return false
				}

				respData := &SshResponse{}
				e = json.Unmarshal(payloadData, respData)
				if e != nil {
					eventErr = &errortypes.ParseError{
						errors.Wrap(e,
							"authority: Failed to unmarshal payload data"),
					}
					return false
				}

				cert, e = utils.UnmarshalSshCertificate(
					respData.Certificate)
				if e != nil {
					eventErr = &errortypes.ParseError{
						errors.Wrap(e,
							"authority: Failed to unmarshal payload data"),
					}
					return false
				}

				certMarshaledByt, err := MarshalCertificate(cert, comment)
				if err != nil {
					eventErr = err
					return false
				}

				certMarshaled = string(certMarshaledByt)

				return false
			})
	}()

	err = event.Publish(db, "pritunl_hsm_send", payload)
	if err != nil {
		return
	}

	waiter.Wait()
	if eventErr != nil {
		cert = nil
		certMarshaled = ""
		logrus.WithFields(logrus.Fields{
			"error": eventErr,
		}).Error("authority: Error getting hsm certificate")
		return
	}

	return
}

func (a *Authority) CreateCertificate(db *database.Database, usr *user.User,
	sshPubKey string) (cert *ssh.Certificate, certMarshaled string,
	err error) {

	keyId := ""
	switch a.KeyIdFormat {
	case Username:
		keyId = usr.Username
		break
	case UsernameId:
		keyId = fmt.Sprintf("%s-%s", usr.Id.Hex(), usr.Username)
		break
	case UsernameStripDomain:
		keyId = strings.Split(usr.Username, "@")[0]
		break
	default:
		keyId = usr.Id.Hex()
	}

	keyId = utils.FilterUnixStr(keyId, 32)

	if a.Type == PritunlHsm {
		cert, certMarshaled, err = a.createCertificateHsm(
			db, usr, keyId, sshPubKey)
	} else {
		cert, certMarshaled, err = a.createCertificateLocal(
			usr, keyId, sshPubKey)
	}

	return
}

func (a *Authority) createHostCertificate(
	hostname, domain string, sshPubKey string) (
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
	_, _ = serialHash.Write([]byte(bson.NewObjectID().Hex()))
	serial := serialHash.Sum64()

	expire := a.HostExpire
	if expire == 0 {
		expire = 600
	}
	validAfter := time.Now().Add(-3 * time.Minute).Unix()
	validBefore := time.Now().Add(
		time.Duration(expire) * time.Minute).Unix()

	cert = &ssh.Certificate{
		Key:             pubKey,
		Serial:          serial,
		CertType:        ssh.HostCert,
		KeyId:           hostname,
		ValidPrincipals: []string{domain},
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

	certMarshaledByt, err := MarshalCertificate(cert, comment)
	if err != nil {
		return
	}

	certMarshaled = string(certMarshaledByt)

	return
}

func (a *Authority) createHostCertificateHsm(db *database.Database,
	hostname, domain string, sshPubKey string) (cert *ssh.Certificate,
	certMarshaled string, err error) {

	pubKey, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(sshPubKey))
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to parse ssh public key"),
		}
		return
	}

	expire := a.HostExpire
	if expire == 0 {
		expire = 600
	}
	validAfter := time.Now().Add(-3 * time.Minute).Unix()
	validBefore := time.Now().Add(
		time.Duration(expire) * time.Minute).Unix()

	cert = &ssh.Certificate{
		Key:             pubKey,
		CertType:        ssh.HostCert,
		KeyId:           hostname,
		ValidPrincipals: []string{domain},
		ValidAfter:      uint64(validAfter),
		ValidBefore:     uint64(validBefore),
	}

	certData, err := utils.MarshalSshCertificate(cert)
	if err != nil {
		return
	}

	data := SshRequest{
		Serial:      a.HsmSerial,
		Certificate: certData,
	}

	cipData, err := json.Marshal(data)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to marshal certificate"),
		}
		return
	}

	pad := 16 - len(cipData)%16
	for i := 0; i < pad; i++ {
		cipData = append(cipData, 0)
	}

	encKeyHash := sha256.New()
	encKeyHash.Write([]byte(a.HsmSecret))
	cipKey := encKeyHash.Sum(nil)

	cipIv, err := utils.RandBytes(aes.BlockSize)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(cipKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to load cipher"),
		}
		return
	}

	mode := cipher.NewCBCEncrypter(block, cipIv)
	mode.CryptBlocks(cipData, cipData)

	if a.HsmSecret == "" {
		err = &errortypes.ReadError{
			errors.Wrap(err, "session: Empty secret"),
		}
		return
	}

	hashFunc := hmac.New(sha512.New, []byte(a.HsmSecret))
	hashFunc.Write(cipData)
	rawSignature := hashFunc.Sum(nil)
	sig := base64.StdEncoding.EncodeToString(rawSignature)

	payloadId := bson.NewObjectID().Hex()
	payload := &HsmPayload{
		Id:        payloadId,
		Token:     a.HsmToken,
		Iv:        cipIv,
		Signature: sig,
		Type:      "ssh_certificate",
		Data:      cipData,
	}

	waiter := sync.WaitGroup{}
	waiter.Add(1)

	timeout := time.Duration(settings.System.HsmResponseTimeout) * time.Second
	start := time.Now()
	var eventErr error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithFields(logrus.Fields{
					"error": errors.New(fmt.Sprintf("%s", r)),
				}).Error("authority: Parse hsm panic")

				eventErr = &errortypes.UnknownError{
					errors.New("authority: Parse hsm panic"),
				}
			}
			waiter.Done()
		}()

		event.SubscribeType([]string{"pritunl_hsm_recv"}, 5*time.Second,
			func() event.CustomEvent {
				return &HsmEvent{}
			},
			func(msgInf event.CustomEvent, e error) bool {
				if e != nil {
					eventErr = e
					return false
				}

				if msgInf == nil || msgInf.GetData() == nil {
					if time.Since(start) < timeout {
						return true
					}

					eventErr = &errortypes.UnknownError{
						errors.New("authority: Timeout waiting for hsm"),
					}
					return false
				}

				msg := msgInf.(*HsmEvent)

				if msg.Data.Id != payloadId ||
					msg.Data.Type != "ssh_certificate" {

					if time.Since(start) < timeout {
						return true
					}
					eventErr = &errortypes.UnknownError{
						errors.New("authority: Timeout waiting for hsm"),
					}
					return false
				}

				payloadData, e := UnmarshalPayload(
					a.HsmToken, a.HsmSecret, msg.Data)
				if e != nil {
					eventErr = e
					return false
				}

				respData := &SshResponse{}
				e = json.Unmarshal(payloadData, respData)
				if e != nil {
					eventErr = &errortypes.ParseError{
						errors.Wrap(e,
							"authority: Failed to unmarshal payload data"),
					}
					return false
				}

				cert, e = utils.UnmarshalSshCertificate(
					respData.Certificate)
				if e != nil {
					eventErr = &errortypes.ParseError{
						errors.Wrap(e,
							"authority: Failed to unmarshal payload data"),
					}
					return false
				}

				certMarshaledByt, err := MarshalCertificate(cert, comment)
				if err != nil {
					eventErr = err
					return false
				}

				certMarshaled = string(certMarshaledByt)

				return false
			})
	}()

	err = event.Publish(db, "pritunl_hsm_send", payload)
	if err != nil {
		return
	}

	waiter.Wait()
	if eventErr != nil {
		cert = nil
		certMarshaled = ""
		logrus.WithFields(logrus.Fields{
			"error": eventErr,
		}).Error("authority: Error getting hsm certificate")
		return
	}

	return
}

func (a *Authority) CreateHostCertificate(db *database.Database,
	hostname string, sshPubKey string) (
	cert *ssh.Certificate, certMarshaled string, err error) {

	if a.Type == PritunlHsm {
		cert, certMarshaled, err = a.createHostCertificateHsm(
			db, hostname, a.GetDomain(hostname), sshPubKey)
	} else {
		cert, certMarshaled, err = a.createHostCertificate(
			hostname, a.GetDomain(hostname), sshPubKey)
	}

	return
}

func (a *Authority) CreateBastionHostCertificate(db *database.Database,
	hostname string, sshPubKey string) (
	cert *ssh.Certificate, certMarshaled string, err error) {

	if a.Type == PritunlHsm {
		cert, certMarshaled, err = a.createHostCertificateHsm(
			db, hostname, hostname, sshPubKey)
	} else {
		cert, certMarshaled, err = a.createHostCertificate(
			hostname, hostname, sshPubKey)
	}

	return
}

func (a *Authority) createRootCertificateLocal() (err error) {
	privateKey, err := ParsePemKey(a.PrivateKey)
	if err != nil {
		return
	}

	pubKey, err := ParseSshPubKey(a.PublicKey)
	if err != nil {
		return
	}

	serialHash := fnv.New64a()
	_, _ = serialHash.Write([]byte(bson.NewObjectID().Hex()))
	serial := &big.Int{}
	serial.SetUint64(serialHash.Sum64())

	notBefore := time.Now().Add(-90 * time.Second)
	notAfter := time.Now().Add(87600 * time.Hour)

	template := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: a.Id.Hex(),
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage: x509.KeyUsageCertSign |
			x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
	}

	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		template,
		template,
		pubKey,
		privateKey,
	)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to create certificate"),
		}
		return
	}

	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}

	a.RootCertificate = strings.TrimSpace(string(pem.EncodeToMemory(block)))

	return
}

func (a *Authority) CreateRootCertificate(db *database.Database) (
	err error) {

	if a.Type == PritunlHsm {
		err = &errortypes.UnknownError{
			errors.Wrap(err,
				"authority: Root certificate not available on HSM"),
		}
	} else {
		err = a.createRootCertificateLocal()
	}

	return
}

func (a *Authority) createClientCertificateLocal() (
	clientCert *tls.Certificate, err error) {

	privateKey, err := ParsePemKey(a.PrivateKey)
	if err != nil {
		return
	}

	block, _ := pem.Decode([]byte(a.RootCertificate))
	if block == nil {
		err = &errortypes.ParseError{
			errors.New("authority: Failed to decode root certificate"),
		}
		return
	}

	rootCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to parse root certificate"),
		}
		return
	}

	clientPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "authority: Failed to generate rsa key"),
		}
		return
	}

	serialHash := fnv.New64a()
	_, _ = serialHash.Write([]byte(bson.NewObjectID().Hex()))
	serial := &big.Int{}
	serial.SetUint64(serialHash.Sum64())

	notBefore := time.Now().Add(-120 * time.Second)
	notAfter := time.Now().Add(120 * time.Second)

	template := &x509.Certificate{
		SerialNumber: serial,
		Issuer: pkix.Name{
			CommonName: a.Id.Hex(),
		},
		Subject: pkix.Name{
			CommonName: node.Self.Id.Hex(),
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
		},
	}

	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		template,
		rootCert,
		clientPrivateKey.Public(),
		privateKey,
	)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err,
				"authority: Failed to create client certificate"),
		}
		return
	}

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}

	certPem := pem.EncodeToMemory(certBlock)

	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(clientPrivateKey),
	}

	clientKeyPem := pem.EncodeToMemory(privateBlock)

	cert, err := tls.X509KeyPair(certPem, clientKeyPem)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err,
				"authority: Failed to create x509 key pair"),
		}
		return
	}

	clientCert = &cert

	return
}

func (a *Authority) CreateClientCertificate(db *database.Database) (
	cert *tls.Certificate, err error) {

	cert = clientCertCacheGet(a)
	if cert != nil {
		return
	}

	if a.Type == PritunlHsm {
		err = &errortypes.UnknownError{
			errors.Wrap(err,
				"authority: Client certificate not available on HSM"),
		}
		return
	} else {
		cert, err = a.createClientCertificateLocal()
		if err != nil {
			return
		}
	}

	clientCertCacheSet(a, cert)

	return
}

func (a *Authority) TokenNew() (err error) {
	if a.HostTokens == nil {
		a.HostTokens = []string{}
	}

	token, err := utils.RandStr(48)
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

func (a *Authority) HandleHsmStatus(db *database.Database,
	payload *HsmPayload) (err error) {

	payloadData, err := UnmarshalPayload(
		a.HsmToken, a.HsmSecret, payload)
	if err != nil {
		return
	}

	respData := &HsmStatus{}
	err = json.Unmarshal(payloadData, respData)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to unmarshal payload data"),
		}
		return
	}

	sendEvent := false
	fields := set.NewSet("hsm_timestamp", "hsm_status")
	a.HsmTimestamp = time.Now()

	status := Disconnected
	if respData.Status == "online" {
		status = Connected
	}

	if a.HsmStatus != status {
		sendEvent = true
	}
	a.HsmStatus = status

	if a.PublicKey != respData.SshPublicKey {
		sendEvent = true
		fields.Add("public_key")
		a.PublicKey = respData.SshPublicKey
	}

	err = a.CommitFields(db, fields)
	if err != nil {
		return
	}

	if sendEvent {
		_ = event.PublishDispatch(db, "authority.change")
	}

	return
}

func (a *Authority) ValidateHsmSignature(
	db *database.Database, token, sig, timeStr, nonc, method,
	path string) (err error) {

	timestampInt, _ := strconv.ParseInt(timeStr, 10, 64)
	if timestampInt == 0 {
		err = &errortypes.AuthenticationError{
			errors.New("authority: Invalid authentication timestamp"),
		}
		return
	}

	timestamp := time.Unix(timestampInt, 0)

	if token == "" || token != a.HsmToken {
		err = &errortypes.AuthenticationError{
			errors.New("authority: Invalid authentication token"),
		}
		return
	}

	if len(nonc) < 16 || len(nonc) > 128 {
		err = &errortypes.AuthenticationError{
			errors.New("authority: Invalid authentication nonce"),
		}
		return
	}

	if utils.SinceAbs(timestamp) > time.Duration(
		settings.Auth.Window)*time.Second {

		err = &errortypes.AuthenticationError{
			errors.New("authority: Authentication timestamp outside window"),
		}
		return
	}

	authString := strings.Join([]string{
		a.HsmToken,
		strconv.FormatInt(timestamp.Unix(), 10),
		nonc,
		method,
		path,
	}, "&")

	err = nonce.Validate(db, nonc)
	if err != nil {
		return
	}

	if a.HsmSecret == "" {
		err = &errortypes.ReadError{
			errors.Wrap(err, "session: Empty secret"),
		}
		return
	}

	hashFunc := hmac.New(sha512.New, []byte(a.HsmSecret))
	hashFunc.Write([]byte(authString))
	rawSignature := hashFunc.Sum(nil)
	testSig := base64.StdEncoding.EncodeToString(rawSignature)

	if subtle.ConstantTimeCompare([]byte(sig), []byte(testSig)) != 1 {
		err = &errortypes.AuthenticationError{
			errors.New("signature: Invalid signature"),
		}
		return
	}

	return
}

func (a *Authority) Export(passphrase string) (encKey string, err error) {
	block, _ := pem.Decode([]byte(a.PrivateKey))
	if block == nil {
		err = &errortypes.ParseError{
			errors.New("authority: Failed to decode private key"),
		}
		return
	}

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

func (a *Authority) JumpProxy() string {
	if a.ProxyHosting {
		return fmt.Sprintf(
			"bastion@%s:%d",
			a.ProxyHostname,
			a.ProxyPort,
		)
	} else {
		return a.HostProxy
	}
}

func (a *Authority) GetMatches() (matches []string, err error) {
	matches = []string{}

	hostDomain := a.GetHostDomain()
	if hostDomain != "" {
		matches = append(matches, hostDomain)
	}

	if a.HostMatches != nil {
		for _, match := range a.HostMatches {
			matches = append(matches, match)
		}
	}

	if a.HostSubnets != nil {
		for _, hostSubnet := range a.HostSubnets {
			match, e := parseSubnetMatch(hostSubnet)
			if e != nil {
				logrus.WithFields(logrus.Fields{
					"authority_id": a.Id.Hex(),
					"subnet":       hostSubnet,
					"error":        e,
				}).Error("authority: Failed to parse subnet match")
				continue
			}

			matches = append(matches, match)
		}

	}

	return
}

func (a *Authority) SetPublicKeyPem() (err error) {
	pubKey, err := ParseSshPubKey(a.PublicKey)
	if err != nil {
		return
	}

	switch pubKey := pubKey.(type) {
	case *rsa.PublicKey:
		keyBytes := x509.MarshalPKCS1PublicKey(pubKey)

		block := &pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: keyBytes,
		}

		a.PublicKeyPem = strings.TrimSpace(string(pem.EncodeToMemory(block)))

		break
	case *ecdsa.PublicKey:
		keyBytes, e := x509.MarshalPKIXPublicKey(pubKey)
		if e != nil {
			err = &errortypes.ParseError{
				errors.Wrap(e, "authority: Failed to parse public key"),
			}
			return
		}

		block := &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: keyBytes,
		}

		a.PublicKeyPem = strings.TrimSpace(string(pem.EncodeToMemory(block)))

		break
	}

	return
}

func (a *Authority) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	a.Name = utils.FilterName(a.Name)

	if a.Type == "" {
		a.Type = Local
	}

	if !a.MatchRoles {
		a.Roles = []string{}
	}

	if a.HostMatches == nil {
		a.HostMatches = []string{}
	}

	if a.HostSubnets == nil {
		a.HostSubnets = []string{}
	}

	switch a.Algorithm {
	case RSA4096:
		break
	case ECP384:
		break
	case "":
		a.Algorithm = RSA4096
		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "invalid_algorithm",
			Message: "Invalid algorithm",
		}
		return
	}

	switch a.Type {
	case Local:
		a.HsmToken = ""
		a.HsmSecret = ""
		a.HsmSerial = ""

		if a.PrivateKey == "" {
			err = a.GeneratePrivateKey()
			if err != nil {
				return
			}
		}

		break
	case PritunlHsm:
		a.PrivateKey = ""

		if a.HsmSerial == "" {
			errData = &errortypes.ErrorData{
				Error:   "missing_hsm_serial",
				Message: "Missing authority HSM serial",
			}
			return
		}

		if a.HsmToken == "" {
			err = a.GenerateHsmToken()
			if err != nil {
				return
			}
		}

		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "invalid_type",
			Message: "Authority type is invalid",
		}
		return
	}

	switch a.KeyIdFormat {
	case UserId:
		break
	case Username:
		break
	case UsernameId:
		break
	case UsernameStripDomain:
		break
	case "":
		a.KeyIdFormat = UserId
	default:
		errData = &errortypes.ErrorData{
			Error:   "invalid_key_id_format",
			Message: "Key ID format is invalid",
		}
		return
	}

	if a.PublicKeyPem == "" {
		err = a.SetPublicKeyPem()
		if err != nil {
			return
		}
	}

	if a.RootCertificate == "" {
		err = a.CreateRootCertificate(db)
		if err != nil {
			return
		}
	}

	if a.Expire < 1 {
		a.Expire = 600
	} else if a.Expire > 1440 {
		a.Expire = 1440
	}

	if a.HostExpire < 1 {
		a.HostExpire = 600
	} else if a.HostExpire > 1440 {
		a.HostExpire = 1440
	} else if a.HostExpire < 15 {
		a.HostExpire = 15
	}

	if a.HostDomain == "" && len(a.HostSubnets) == 0 &&
		len(a.HostMatches) == 0 {

		a.HostCertificates = false
		a.HostProxy = ""
	}

	if a.HostDomain == "" {
		a.StrictHostChecking = false
	}

	if !a.ProxyHosting {
		a.ProxyPort = 0
		a.ProxyHostname = ""

		err = RemoveNode(db, a.Id)
		if err != nil {
			return
		}
	} else {
		a.HostProxy = ""
		if a.ProxyHostname == "" {
			errData = &errortypes.ErrorData{
				Error:   "proxy_hostname_missing",
				Message: "Bastion hosting hostname required",
			}
			return
		}

		if a.ProxyPort < 1 || a.ProxyPort > 65535 {
			errData = &errortypes.ErrorData{
				Error:   "proxy_port_invalid",
				Message: "Bastion hosting port is invalid",
			}
			return
		}
	}

	if a.HostCertificates && a.HostDomain == "" {
		errData = &errortypes.ErrorData{
			Error:   "host_domain_required",
			Message: "Host domain must be set for host certificates",
		}
		return
	}

	if !a.HostCertificates {
		a.StrictHostChecking = false
		a.HostTokens = []string{}
	}

	if a.HostTokens == nil || !a.HostCertificates {
		a.HostTokens = []string{}
	}

	for _, hostSubnet := range a.HostSubnets {
		_, e := parseSubnetMatch(hostSubnet)
		if e != nil {
			errData = &errortypes.ErrorData{
				Error:   "subnet_invalid",
				Message: "Match subnet is invalid",
			}
			return
		}
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

func (a *Authority) Json() {
	if a.Type == PritunlHsm {
		if time.Since(a.HsmTimestamp) > 45*time.Second {
			a.HsmStatus = Disconnected
		}
	}

	a.ProxyJump = a.JumpProxy()
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

	if !a.Id.IsZero() {
		err = &errortypes.DatabaseError{
			errors.New("authority: Authority already exists"),
		}
		return
	}

	_, err = coll.InsertOne(db, a)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func init() {
	module := requires.New("authority")
	module.After("settings")

	module.Handler = func() (err error) {
		db := database.GetDatabase()
		defer db.Close()

		clientCertCacheWatch()

		authrs, err := GetAll(db)
		if err != nil {
			return
		}

		for _, authr := range authrs {
			if !authr.HostCertificates && authr.HostDomain != "" &&
				authr.HostTokens != nil && len(authr.HostTokens) > 0 {

				authr.HostCertificates = true
				err = authr.CommitFields(db, set.NewSet("host_certificates"))
				if err != nil {
					return
				}
			}

			if authr.PublicKeyPem == "" {
				err = authr.SetPublicKeyPem()
				if err != nil {
					return
				}

				err = authr.CommitFields(db, set.NewSet("public_key_pem"))
				if err != nil {
					return
				}
			}

			if authr.RootCertificate == "" {
				err = authr.CreateRootCertificate(db)
				if err != nil {
					return
				}

				err = authr.CommitFields(db, set.NewSet("root_certificate"))
				if err != nil {
					return
				}
			}
		}

		return
	}
}
