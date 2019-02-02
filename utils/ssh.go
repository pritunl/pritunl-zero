package utils

import (
	"encoding/json"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"golang.org/x/crypto/ssh"
)

type sshCertPermissions struct {
	CriticalOptions map[string]string `json:"critical_options"`
	Extensions      map[string]string `json:"extensions"`
}

type sshCertSignature struct {
	Format string `json:"format"`
	Blob   []byte `json:"blob"`
}

type sshCert struct {
	Nonce           []byte             `json:"nonce"`
	Key             []byte             `json:"key"`
	Serial          uint64             `json:"serial"`
	CertType        uint32             `json:"cert_type"`
	KeyId           string             `json:"key_id"`
	ValidPrincipals []string           `json:"valid_principals"`
	ValidAfter      uint64             `json:"valid_after"`
	ValidBefore     uint64             `json:"valid_before"`
	Permissions     sshCertPermissions `json:"permissions"`
	Reserved        []byte             `json:"reserved"`
	SignatureKey    []byte             `json:"signature_key"`
	Signature       *sshCertSignature  `json:"signature"`
}

func MarshalSshCertificate(inCert *ssh.Certificate) (data []byte, err error) {
	outCert := &sshCert{
		Nonce:           inCert.Nonce,
		Serial:          inCert.Serial,
		CertType:        inCert.CertType,
		KeyId:           inCert.KeyId,
		ValidPrincipals: inCert.ValidPrincipals,
		ValidAfter:      inCert.ValidAfter,
		ValidBefore:     inCert.ValidBefore,
		Permissions: sshCertPermissions{
			CriticalOptions: inCert.Permissions.CriticalOptions,
			Extensions:      inCert.Permissions.Extensions,
		},
		Reserved: inCert.Reserved,
	}

	if inCert.Key != nil {
		outCert.Key = inCert.Key.Marshal()
	}

	if inCert.SignatureKey != nil {
		outCert.SignatureKey = inCert.SignatureKey.Marshal()
	}

	if inCert.Signature != nil {
		outCert.Signature = &sshCertSignature{
			Format: inCert.Signature.Format,
			Blob:   inCert.Signature.Blob,
		}
	}

	data, err = json.Marshal(outCert)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "utils: Failed to marshal certificate"),
		}
		return
	}

	return
}

func UnmarshalSshCertificate(data []byte) (
	outCert *ssh.Certificate, err error) {

	inCert := &sshCert{}
	err = json.Unmarshal(data, inCert)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "utils: Failed to unmarshal certificate"),
		}
		return
	}

	outCert = &ssh.Certificate{
		Nonce:           inCert.Nonce,
		Serial:          inCert.Serial,
		CertType:        inCert.CertType,
		KeyId:           inCert.KeyId,
		ValidPrincipals: inCert.ValidPrincipals,
		ValidAfter:      inCert.ValidAfter,
		ValidBefore:     inCert.ValidBefore,
		Permissions: ssh.Permissions{
			CriticalOptions: inCert.Permissions.CriticalOptions,
			Extensions:      inCert.Permissions.Extensions,
		},
		Reserved: inCert.Reserved,
	}

	if inCert.Key != nil {
		out, e := ssh.ParsePublicKey(inCert.Key)
		if e != nil {
			err = &errortypes.ParseError{
				errors.Wrap(e, "utils: Failed to unmarshal key"),
			}
			return
		}

		outCert.Key = out
	}

	if inCert.SignatureKey != nil {
		out, e := ssh.ParsePublicKey(inCert.SignatureKey)
		if e != nil {
			err = &errortypes.ParseError{
				errors.Wrap(e, "utils: Failed to unmarshal key"),
			}
			return
		}

		outCert.SignatureKey = out
	}

	if inCert.Signature != nil {
		outCert.Signature = &ssh.Signature{
			Format: inCert.Signature.Format,
			Blob:   inCert.Signature.Blob,
		}
	}

	return
}
