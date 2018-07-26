package letsencrypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/pritunl/pritunl-zero/letsencrypt/internal/base64"
	"github.com/square/go-jose"
)

const (
	ChallengeDNS    = "dns-01"
	ChallengeHTTP   = "http-01"
	ChallengeTLSSNI = "tls-sni-01"
)

// HTTP returns a URL path and HTTP response body that the ACME server will
// check when verifying the challenge.
func (chal Challenge) HTTP(accountKey interface{}) (urlPath, resource string, err error) {
	if chal.Type != ChallengeHTTP {
		return "", "", fmt.Errorf("challenge type is %s not %s", chal.Type, ChallengeHTTP)
	}

	urlPath = path.Join("/.well-known/acme-challenge", chal.Token)
	resource, err = keyAuth(accountKey, chal.Token)
	return
}

// TLSSNI returns TLS certificates for a set of server names.
// The ACME server will make a TLS Server Name Indication handshake with the
// given domain. The domain must present the returned certifiate for each name.
func (chal Challenge) TLSSNI(accountKey interface{}) (map[string]*tls.Certificate, error) {
	if chal.Type != ChallengeTLSSNI {
		return nil, fmt.Errorf("challenge type is %s not %s", chal.Type, ChallengeTLSSNI)
	}

	auth, err := keyAuth(accountKey, chal.Token)
	if err != nil {
		return nil, err
	}

	// private key for creating self-signed certificates
	// TODO: Make configurable?
	tlsKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	compute := func(content string) string {
		hash := sha256.Sum256([]byte(content))
		return hex.EncodeToString(hash[:])
	}

	// compute z0 ... zN+1 for the challenge
	z := make([]string, chal.N+1)
	z[0] = compute(auth)
	for i := 0; i < chal.N; i++ {
		z[i+1] = compute(z[i])
	}

	// crypto/tls library takes a PEM and DER encoded private key
	certKeyPEM := pemEncode(x509.MarshalPKCS1PrivateKey(tlsKey), "RSA PRIVATE KEY")

	certs := make(map[string]*tls.Certificate)
	for _, zi := range z {
		name := zi[:32] + "." + zi[32:64] + ".acme.invalid"

		// initialize server certificate template with sane defaults
		tmpl, err := certTmpl()
		if err != nil {
			return nil, err
		}

		tmpl.SignatureAlgorithm = x509.SHA256WithRSA
		tmpl.DNSNames = []string{name}

		derCert, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &tlsKey.PublicKey, tlsKey)
		if err != nil {
			return nil, fmt.Errorf("create self-signed certificate: %v", err)
		}
		certPEM := pemEncode(derCert, "CERTIFICATE")
		cert, err := tls.X509KeyPair(certPEM, certKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("loading x509 key pair: %v", err)
		}

		certs[name] = &cert
	}
	return certs, nil
}

// DNS returns the subdomain name and the TXT value you need to set for that
// subdomain. The ACME server will make DNS TXT lookup on that subdomain and
// verify that the value matches. Keep in mind that DNS TTL's might prevent
// the lookup from working correctly the first few times and ChallengeReady
// will continue to loop if the record is missing/invalid. It is recommended
// that you set the record to the lowest TTL allowed by your provider.
func (chal Challenge) DNS(accountKey interface{}) (subdomain, txt string, err error) {
	if chal.Type != ChallengeDNS {
		return "", "", fmt.Errorf("challenge type is %s not %s", chal.Type, ChallengeDNS)
	}
	auth, err := keyAuth(accountKey, chal.Token)
	if err != nil {
		return "", "", err
	}
	hash := sha256.Sum256([]byte(auth))
	txt = base64.RawURLEncoding.EncodeToString(hash[:])
	subdomain = "_acme-challenge"
	return
}

// Not yet implemented
func (chal Challenge) ProofOfPossession(accountKey, certKey interface{}) (Challenge, error) {
	return Challenge{}, errors.New("proof of possession not implemented")
}

// ChallengeReady informs the server that the provided challenge is ready
// for verification.
//
// The client then begins polling the server for confirmation on the
// result of the status.
func (c *Client) ChallengeReady(accountKey interface{}, chal Challenge) error {
	switch chal.Type {
	case ChallengeHTTP, ChallengeTLSSNI, ChallengeDNS:
	default:
		return fmt.Errorf("unsupported challenge type '%s'", chal.Type)
	}
	auth, err := keyAuth(accountKey, chal.Token)
	if err != nil {
		return err
	}
	data := struct {
		Resource string `json:"resource"`
		KeyAuth  string `json:"keyAuthorization"`
		Type     string `json:"type"`
		Token    string `json:"token"`
	}{resourceChallenge, auth, chal.Type, chal.Token}
	sig, err := c.signObject(accountKey, &data)
	if err != nil {
		return err
	}
	resp, err := c.client.Post(chal.URI, jwsContentType, strings.NewReader(sig))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := checkHTTPError(resp, http.StatusAccepted); err != nil {
		return err
	}

	// Begin polling the server and check if the challenge is completed.
	pollInterval := c.PollInterval
	if pollInterval == 0 {
		pollInterval = 500 * time.Millisecond
	}
	pollTimeout := c.PollTimeout
	if pollTimeout == 0 {
		pollTimeout = 30 * time.Second
	}
	start := time.Now()
	for {
		if time.Now().Sub(start) > pollTimeout {
			if chal.Error != nil {
				return chal.Error
			}
			return errors.New("polling pending challenge timed out")
		}
		chal, err = c.Challenge(chal.URI)
		if err != nil {
			return err
		}
		switch chal.Status {
		case StatusPending, "":
			time.Sleep(pollInterval)
		case StatusInvalid:
			if chal.Error == nil {
				return errors.New("challenge returned status 'invalid' without explicit error")
			}
			// if the challenge was DNS, keep trying since the value is cached
			if chal.Type == ChallengeDNS {
				// unauthorized is the TXT value is wrong or not found
				// connection if NXDOMAIN
				if chal.Error.Typ == "urn:acme:error:unauthorized" || chal.Error.Typ == "urn:acme:error:connection" {
					time.Sleep(pollInterval)
					continue
				}
			}
			return chal.Error
		case StatusValid:
			return nil
		default:
			return fmt.Errorf("unexpected challenge status: %s", chal.Status)
		}
	}
}

func keyAuth(key interface{}, token string) (string, error) {
	thumbprint, err := (&jose.JsonWebKey{Key: key}).Thumbprint(crypto.SHA256)
	if err != nil {
		return "", fmt.Errorf("compute key thumbprint: %v", err)
	}
	return token + "." + base64.RawURLEncoding.EncodeToString(thumbprint), nil
}

func certTmpl() (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.New("failed to generate serial number: " + err.Error())
	}

	// TODO: make certificate expire when the authorization challenge expires
	tmpl := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24),
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	return &tmpl, nil
}

func pemEncode(block []byte, typ string) []byte {
	b := pem.Block{Type: typ, Bytes: block}
	return pem.EncodeToMemory(&b)
}
