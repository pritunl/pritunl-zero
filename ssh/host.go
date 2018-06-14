package ssh

import (
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/authority"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"time"
)

func NewHostCertificate(db *database.Database, hostname string, port int,
	tokens []string, r *http.Request, pubKey string) (
	cert *Certificate, errData *errortypes.ErrorData, err error) {

	pubKey = strings.TrimSpace(pubKey)

	if len(tokens) > settings.System.SshHostTokenLen {
		err = errortypes.ParseError{
			errors.New("ssh: Too many tokens"),
		}
		return
	}

	if len(pubKey) > settings.System.SshPubKeyLen {
		err = errortypes.ParseError{
			errors.New("ssh: Public key too long"),
		}
		return
	}

	agnt, err := agent.Parse(db, r)
	if err != nil {
		return
	}

	cert = &Certificate{
		Id:               bson.NewObjectId(),
		AuthorityIds:     []bson.ObjectId{},
		Timestamp:        time.Now(),
		PubKey:           pubKey,
		Certificates:     []string{},
		CertificatesInfo: []*Info{},
		Agent:            agnt,
	}

	authrs, err := authority.GetTokens(db, tokens)
	if err != nil {
		return
	}

	if len(authrs) == 0 {
		errData = &errortypes.ErrorData{
			Error:   "invalid_tokens",
			Message: "All tokens are invalid",
		}
		return
	}

	for _, authr := range authrs {
		if !authr.HostnameValidate(hostname, port, pubKey) {
			continue
		}

		crt, certStr, e := authr.CreateHostCertificate(hostname, pubKey)
		if e != nil {
			err = e
			return
		}

		info := &Info{
			Expires:    time.Unix(int64(crt.ValidBefore), 0),
			Serial:     fmt.Sprintf("%d", crt.Serial),
			Principals: crt.ValidPrincipals,
			Extensions: []string{},
		}

		for permission := range crt.Permissions.Extensions {
			info.Extensions = append(info.Extensions, permission)
		}

		cert.AuthorityIds = append(cert.AuthorityIds, authr.Id)
		cert.Certificates = append(cert.Certificates, certStr)
		cert.CertificatesInfo = append(cert.CertificatesInfo, info)
	}

	if len(cert.Certificates) == 0 {
		errData = &errortypes.ErrorData{
			Error:   "certificate_unavailable",
			Message: "No certificates are available",
		}
		return
	}

	err = cert.Insert(db)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func NewBastionHostCertificate(db *database.Database, hostname,
	pubKey string, authr *authority.Authority) (
	cert *Certificate, err error) {

	pubKey = strings.TrimSpace(pubKey)

	if len(pubKey) > settings.System.SshPubKeyLen {
		err = errortypes.ParseError{
			errors.New("ssh: Public key too long"),
		}
		return
	}

	cert = &Certificate{
		Id:               bson.NewObjectId(),
		AuthorityIds:     []bson.ObjectId{},
		Timestamp:        time.Now(),
		PubKey:           pubKey,
		Certificates:     []string{},
		CertificatesInfo: []*Info{},
	}

	crt, certStr, err := authr.CreateHostCertificate(hostname, pubKey)
	if err != nil {
		return
	}

	info := &Info{
		Expires:    time.Unix(int64(crt.ValidBefore), 0),
		Serial:     fmt.Sprintf("%d", crt.Serial),
		Principals: crt.ValidPrincipals,
		Extensions: []string{},
	}

	for permission := range crt.Permissions.Extensions {
		info.Extensions = append(info.Extensions, permission)
	}

	cert.AuthorityIds = append(cert.AuthorityIds, authr.Id)
	cert.Certificates = append(cert.Certificates, certStr)
	cert.CertificatesInfo = append(cert.CertificatesInfo, info)

	err = cert.Insert(db)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
