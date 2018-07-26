package letsencrypt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pritunl/pritunl-zero/letsencrypt/internal/base64"
	"gopkg.in/square/go-jose.v1"
)

const jwsContentType = "application/jose+jws"

const (
	resourceNewRegistration      = "new-reg"
	resourceRecoverRegistation   = "recover-reg"
	resourceNewAuthorization     = "new-authz"
	resourceNewCertificate       = "new-cert"
	resourceNewRevokeCertificate = "revoke-cert"
	resourceRegistration         = "reg"
	resourceAuthorization        = "authz"
	resourceChallenge            = "challenge"
	resourceCertificate          = "cert"
)

type directory struct {
	NewRegistration      string `json:"new-reg"`
	RecoverRegistation   string `json:"recover-reg"`
	NewAuthorization     string `json:"new-authz"`
	NewCertificate       string `json:"new-cert"`
	NewRevokeCertificate string `json:"revoke-cert"`
	Registration         string `json:"reg"`
	Authorization        string `json:"authz"`
	Challenge            string `json:"challenge"`
	Certificate          string `json:"cert"`
	Terms                string `json:"terms"`
}

// Paths taken directly from boulder's source code.
// There are quite a few paths missing in the /directory object
// for boulder's current implementation.
// When those are missing default to these.
// See: https://github.com/letsencrypt/boulder/issues/754
const (
	boulderDirectoryPath  = "/directory"
	boulderNewRegPath     = "/acme/new-reg"
	boulderRegPath        = "/acme/reg/"
	boulderNewAuthzPath   = "/acme/new-authz"
	boulderAuthzPath      = "/acme/authz/"
	boulderNewCertPath    = "/acme/new-cert"
	boulderCertPath       = "/acme/cert/"
	boulderRevokeCertPath = "/acme/revoke-cert"
	boulderTermsPath      = "/terms"
	boulderIssuerPath     = "/acme/issuer-cert"
	boulderBuildIDPath    = "/build"
)

var (
	// errUnsupportedRSABitLen reports whether an unrecognized RSA key size was used
	errUnsupportedRSABitLen = errors.New("unsupported RSA bit length")
	// errUnsupportedECDSACurve reports whether an unrecognized ECDSA curve was used
	errUnsupportedECDSACurve = errors.New("unsupported ECDSA curve")
)

func newDefaultDirectory(baseURL *url.URL) directory {
	pathToURL := func(path string) string {
		var u url.URL
		u = *baseURL
		u.Path = path
		return u.String()
	}

	return directory{
		NewRegistration:      pathToURL(boulderNewRegPath),
		NewAuthorization:     pathToURL(boulderNewAuthzPath),
		NewCertificate:       pathToURL(boulderNewCertPath),
		NewRevokeCertificate: pathToURL(boulderRevokeCertPath),
		Registration:         pathToURL(boulderRegPath),
		Authorization:        pathToURL(boulderAuthzPath),
		Certificate:          pathToURL(boulderCertPath),
		Terms:                pathToURL(boulderTermsPath),
	}
}

// Client is a client for a single ACME server.
type Client struct {
	// PollInterval determines how quickly the client will
	// request updates on a challenge from the ACME server.
	// If unspecified, it defaults to 500 milliseconds.
	PollInterval time.Duration
	// Amount of time after the client notifies the server a challenge is
	// ready, and when it will stop checking for updates.
	// If unspecified, it defaults to 30 seconds.
	PollTimeout time.Duration

	resources directory

	client      *http.Client
	nonceSource jose.NonceSource

	terms string
}

// Terms returns the URL of the server's terms of service.
// All accounts registered using this client automatically
// accept these terms.
func (c *Client) Terms() string {
	return c.terms
}

// NewClient creates a client of a ACME server by querying the server's
// resource directory and attempting to resolve the URL of the terms of service.
func NewClient(directoryURL string) (*Client, error) {
	return NewClientWithTransport(directoryURL, nil)
}

// NewClientWithTransport creates a client of a ACME server by querying the server's
// resource directory and attempting to resolve the URL of the terms of service.
func NewClientWithTransport(directoryURL string, t http.RoundTripper) (*Client, error) {
	u, err := url.Parse(directoryURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL %s: %v", directoryURL, err)
	}
	if u.Path == "" {
		u.Path = boulderDirectoryPath
	}
	nrt := newNonceRoundTripper(t)

	c := &Client{
		client:      &http.Client{Transport: nrt},
		resources:   newDefaultDirectory(u),
		nonceSource: nrt,
	}

	resp, err := c.client.Get(directoryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkHTTPError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&c.resources); err != nil {
		return nil, fmt.Errorf("could not decode body: %v %s", err, resp.Body)
	}

	termsResp, err := c.client.Get(c.resources.Terms)
	if err != nil {
		return nil, fmt.Errorf("GET failed: %v", err)
	}
	defer termsResp.Body.Close()
	if err := checkHTTPError(termsResp, http.StatusOK); err != nil {
		return nil, fmt.Errorf("failed to get terms of service: %v", err)
	}
	c.terms = termsResp.Request.URL.String()

	return c, nil
}

// UpdateRegistration sends the updated registration object to the server.
func (c *Client) UpdateRegistration(accountKey interface{}, reg Registration) (Registration, error) {
	url := c.resources.Registration + strconv.Itoa(reg.Id)
	return c.registration(accountKey, reg, resourceRegistration, url)
}

// NewRegistration registers a key pair with the ACME server.
// If the key pair is already registered, the registration object is recovered.
func (c *Client) NewRegistration(accountKey interface{}) (reg Registration, err error) {
	reg, err = c.registration(accountKey, Registration{}, resourceNewRegistration, c.resources.NewRegistration)
	if err != nil || reg.Agreement == c.Terms() {
		return
	}
	reg.Agreement = c.Terms()
	reg, err = c.UpdateRegistration(accountKey, reg)
	return reg, err
}

func (c *Client) registration(accountKey interface{}, reg Registration, resource, url string) (Registration, error) {
	reg.Resource = resource
	sig, err := c.signObject(accountKey, &reg)
	if err != nil {
		return Registration{}, err
	}
	resp, err := c.client.Post(url, jwsContentType, strings.NewReader(sig))
	if err != nil {
		return Registration{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict && resource == resourceNewRegistration {
		// We must send our agreement in order to get a non-empty registration back
		return c.registration(accountKey, Registration{
			Agreement: c.Terms(),
		}, resourceRegistration, resp.Header.Get("Location"))
	}

	statusExp := http.StatusCreated
	if resource == resourceRegistration {
		statusExp = http.StatusAccepted
	}

	if err := checkHTTPError(resp, statusExp); err != nil {
		return Registration{}, err
	}

	var updatedReg Registration
	if err := json.NewDecoder(resp.Body).Decode(&updatedReg); err != nil {
		return Registration{}, fmt.Errorf("unmarshalling response body: %v", err)
	}

	return updatedReg, nil
}

// NewAuthorization requests a set of challenges from the server to prove
// ownership of a given resource.
// Only known type is 'dns'.
//
// NOTE: Currently the only way to recover an authorization object is with
// the returned authorization URL.
func (c *Client) NewAuthorization(accountKey interface{}, typ, val string) (auth Authorization, authURL string, err error) {
	type Identifier struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	data := struct {
		Resource   string     `json:"resource"`
		Identifier Identifier `json:"identifier"`
	}{
		resourceNewAuthorization,
		Identifier{typ, val},
	}
	payload, err := c.signObject(accountKey, &data)
	if err != nil {
		return auth, "", err
	}
	resp, err := c.client.Post(c.resources.NewAuthorization, jwsContentType, strings.NewReader(payload))
	if err != nil {
		return auth, "", err
	}
	defer resp.Body.Close()
	if err = checkHTTPError(resp, http.StatusCreated); err != nil {
		return auth, "", err
	}

	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return auth, "", fmt.Errorf("decoding response body: %v", err)
	}
	return auth, resp.Header.Get("Location"), nil
}

// Authorization returns the authorization object associated with
// the given authorization URI.
func (c *Client) Authorization(authURI string) (Authorization, error) {
	var auth Authorization
	resp, err := c.client.Get(authURI)
	if err != nil {
		return auth, err
	}
	defer resp.Body.Close()
	if err = checkHTTPError(resp, http.StatusOK); err != nil {
		return auth, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return auth, fmt.Errorf("decoding response body: %v", err)
	}
	return auth, nil
}

// Challenge returns the challenge object associated with the
// given challenge URI.
func (c *Client) Challenge(chalURI string) (Challenge, error) {
	var chal Challenge
	resp, err := c.client.Get(chalURI)
	if err != nil {
		return chal, err
	}
	defer resp.Body.Close()
	if err = checkHTTPError(resp, http.StatusAccepted); err != nil {
		return chal, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&chal); err != nil {
		return chal, fmt.Errorf("decoding response body: %v", err)
	}
	return chal, nil
}

// NewCertificate requests a certificate from the ACME server.
//
// csr must have already been signed by a private key.
func (c *Client) NewCertificate(accountKey interface{}, csr *x509.CertificateRequest) (*CertificateResponse, error) {
	if csr == nil || csr.Raw == nil {
		return nil, errors.New("invalid certificate request object")
	}
	payload := struct {
		Resource string `json:"resource"`
		CSR      string `json:"csr"`
	}{
		resourceNewCertificate,
		base64.RawURLEncoding.EncodeToString(csr.Raw),
	}
	data, err := c.signObject(accountKey, &payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(c.resources.NewCertificate, jwsContentType, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := checkHTTPError(resp, http.StatusCreated); err != nil {
		return nil, err
	}

	return handleCertificateResponse(resp)
}

// RenewCertificate attempts to renew an existing certificate.
// Let's Encrypt may return the same certificate. You should load your
// current x509.Certificate and use the Equal method to compare to the "new"
// certificate. If it's identical, you'll need to run NewCertificate and/or
// start a new certificate flow.
func (c *Client) RenewCertificate(certURI string) (*CertificateResponse, error) {
	resp, err := c.client.Get(certURI)
	if err != nil {
		return nil, fmt.Errorf("renew certificate http error: %s", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, errors.New("certificate not available. Start a new certificate flow")
	}

	certResp, err := handleCertificateResponse(resp)
	if err == nil {
		if certResp.URI == "" {
			certResp.URI = certURI
		}
	}

	return certResp, err
}

// RevokeCertificate takes a PEM encoded certificate or bundle and
// attempts to revoke it.
func (c *Client) RevokeCertificate(accountKey interface{}, pemBytes []byte) error {
	certificates, err := parsePEMBundle(pemBytes)
	if err != nil {
		return err
	}

	cert := certificates[0]
	if cert.IsCA {
		return errors.New("Certificate bundle starts with a CA certificate")
	}

	// cert.Raw holds DERbytes, which need to be encoded to base64 per acme spec
	encoded := base64.URLEncoding.EncodeToString(cert.Raw)
	payload := struct {
		Resource    string `json:"resource"`
		Certificate string `json:"certificate"`
	}{
		resourceNewRevokeCertificate,
		encoded,
	}
	data, err := c.signObject(accountKey, &payload)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(c.resources.NewRevokeCertificate, jwsContentType, strings.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := checkHTTPError(resp, http.StatusOK); err != nil {
		return err
	}

	return nil
}

func handleCertificateResponse(resp *http.Response) (*CertificateResponse, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %v", err)
	}
	defer resp.Body.Close()

	// Certificate is not yet available. Gather data and retry later
	if len(body) == 0 {
		retryAfter, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			return nil, fmt.Errorf("Error parsing retry-after header: %s", err)
		}

		return &CertificateResponse{
			RetryAfter: retryAfter,
			URI:        resp.Header.Get("Location"),
		}, nil
	}

	// Certificate was available in response body
	x509Cert, err := x509.ParseCertificate(body)
	if err != nil {
		return nil, fmt.Errorf("Error parsing x509 certificate: %s", err)
	}

	links := parseLinks(resp.Header["Link"])
	return &CertificateResponse{
		Certificate: x509Cert,
		URI:         resp.Header.Get("Location"),
		StableURI:   resp.Header.Get("Content-Location"),
		Issuer:      links["up"],
	}, nil
}

// TODO: doesn't need to be a function on the client struct
func (c *Client) signObject(accountKey interface{}, v interface{}) (string, error) {
	var (
		signer jose.Signer
		alg    jose.SignatureAlgorithm
		err    error
	)

	switch accountKey := accountKey.(type) {
	case *rsa.PrivateKey:
		modulus := accountKey.N
		bitLen := modulus.BitLen()

		switch bitLen {
		case 2048:
			alg = jose.RS256
		// Not yet supported by LetsEncrypt's Boulder service: https://github.com/letsencrypt/boulder/issues/1592
		// case 3072:
		// 	alg = jose.RS384
		// case 4096:
		// 	alg = jose.RS512
		default:
			return "", errUnsupportedRSABitLen
		}
	case *ecdsa.PrivateKey:
		switch accountKey.Params() {
		case elliptic.P256().Params():
			alg = jose.ES256
		case elliptic.P384().Params():
			alg = jose.ES384
		// Not yet supported by LetsEncrypt's Boulder service: https://github.com/letsencrypt/boulder/issues/1592
		// case elliptic.P521().Params():
		// 	alg = jose.ES512
		default:
			return "", errUnsupportedECDSACurve
		}
	default:
		err = errors.New("acme: unsupported private key type")
	}

	signer, err = jose.NewSigner(alg, accountKey)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	signer.SetNonceSource(c.nonceSource)
	sig, err := signer.Sign(data)
	if err != nil {
		return "", err
	}
	return sig.FullSerialize(), nil
}

var aBrkt = regexp.MustCompile("[<>]")
var slver = regexp.MustCompile("(.+) *= *\"(.+)\"")

func parseLinks(links []string) map[string]string {
	linkMap := make(map[string]string)

	for _, link := range links {
		link = aBrkt.ReplaceAllString(link, "")
		parts := strings.Split(link, ";")

		matches := slver.FindStringSubmatch(parts[1])
		if len(matches) > 0 {
			linkMap[matches[2]] = parts[0]
		}
	}

	return linkMap
}

// parsePEMBundle parses a certificate bundle from top to bottom and returns
// a slice of x509 certificates. This function will error if no certificates are found.
// Credit: github.com/xenolf/lego
func parsePEMBundle(bundle []byte) ([]*x509.Certificate, error) {
	var certificates []*x509.Certificate

	remaining := bundle
	for len(remaining) != 0 {
		certBlock, rem := pem.Decode(remaining)
		// Thanks golang for having me do this :[
		remaining = rem
		if certBlock == nil {
			return nil, errors.New("Could not decode certificate.")
		}

		cert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return nil, err
		}

		certificates = append(certificates, cert)
	}

	if len(certificates) == 0 {
		return nil, errors.New("No certificates were found while parsing the bundle.")
	}

	return certificates, nil
}
