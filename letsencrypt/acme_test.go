package letsencrypt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	testDomain = "example.org"
	testURL    = "http://localhost:4000/directory"
)

func TestNewClient(t *testing.T) {
	if _, err := NewClient(testURL); err != nil {
		t.Fatal(err)
	}
}

func TestNewClientWithTransport(t *testing.T) {
	rt := http.DefaultTransport
	if _, err := NewClientWithTransport(testURL, rt); err != nil {
		t.Fatal(err)
	}
}

func TestRegistration(t *testing.T) {
	tests := []struct {
		keyType string
		bitSize int
		err     error
	}{
		{"ecdsa", 256, nil},
		{"ecdsa", 384, nil},
		//{"ecdsa", 512, nil},
		{"rsa", 2048, nil},
		//{"rsa", 3072, nil},
		//{"rsa", 4096, nil},
	}

	cli, err := NewClient(testURL)
	//cli, err := NewClient("https://acme-staging.api.letsencrypt.org/directory")
	//cli, err := NewClient("https://acme-v01.api.letsencrypt.org/directory")
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		var priv interface{}
		var err error

		switch tt.keyType {
		case "rsa":
			priv, err = rsa.GenerateKey(rand.Reader, tt.bitSize)
		case "ecdsa":
			var curve elliptic.Curve
			switch tt.bitSize {
			case 256:
				curve = elliptic.P256()
			case 384:
				curve = elliptic.P384()
			case 512:
				curve = elliptic.P521()
			}
			priv, err = ecdsa.GenerateKey(curve, rand.Reader)
		}

		if err != nil {
			t.Fatal(err)
		}

		_, err = cli.NewRegistration(priv)
		if err != tt.err {
			log.Printf("key type: %s\n", tt.keyType)
			log.Printf("bit size: %d\n", tt.bitSize)
			t.Fatal(err)
		}
	}
}

func TestUpdateRegistration(t *testing.T) {
	cli, err := NewClient(testURL)
	if err != nil {
		t.Fatal(err)
	}
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	reg, err := cli.NewRegistration(priv)
	if err != nil {
		t.Fatal(err)
	}
	reg.Contact = []string{"mailto:cert-admin@example.com"}
	updatedReg, err := cli.UpdateRegistration(priv, reg)
	if err != nil {
		t.Fatal(err)
	}
	if len(updatedReg.Contact) != 1 {
		t.Errorf("expected update to add one contact, got %s", updatedReg.Contact)
	}
	recoveredReg, err := cli.NewRegistration(priv)
	if err != nil {
		t.Errorf("expected recovered reg to have two contacts, got %s", recoveredReg.Contact)
	}
}

func TestNewAuthorization(t *testing.T) {
	cli, err := NewClient(testURL)
	if err != nil {
		t.Fatal(err)
	}
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cli.NewRegistration(priv); err != nil {
		t.Fatal(err)
	}

	auth, authURL, err := cli.NewAuthorization(priv, "dns", testDomain)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = cli.NewAuthorization(priv, "dns", testDomain); err != nil {
		t.Fatal(err)
	}
	if _, err := cli.NewRegistration(priv); err != nil {
		t.Fatal(err)
	}
	recoveredAuth, err := cli.Authorization(authURL)
	if err != nil {
		t.Fatal(err)
	}
	if auth.Identifier.Value != recoveredAuth.Identifier.Value {
		t.Error("recovered auth did not match original auth")
	}
}

func TestAuthorizationChallenges(t *testing.T) {
	cli, err := NewClient(testURL)
	if err != nil {
		t.Fatal(err)
	}
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cli.NewRegistration(priv); err != nil {
		t.Fatal(err)
	}

	auth, _, err := cli.NewAuthorization(priv, "dns", testDomain)
	if err != nil {
		t.Fatal(err)
	}
	for _, chal := range auth.Challenges {
		if _, err := cli.Challenge(chal.URI); err != nil {
			t.Errorf("failed to get challenge from URI %s: %v", chal.URI, err)
		}
	}
}

func TestNewCertificate(t *testing.T) {
	requiresEtcHostsEdits(t)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	cli, err := NewClient(testURL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cli.NewRegistration(priv); err != nil {
		t.Fatal(err)
	}
	auth, _, err := cli.NewAuthorization(priv, "dns", testDomain)
	if err != nil {
		t.Fatal(err)
	}

	chals := auth.Combinations(ChallengeHTTP)
	if len(chals) == 0 || len(chals[0]) != 1 {
		t.Fatal("no supported challenges")
	}
	chal := chals[0][0]
	urlPath, resource, err := chal.HTTP(priv)
	if err != nil {
		t.Fatal(err)
	}

	hf := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != urlPath {
			t.Error("server request did not match path. expecting", urlPath, "got", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		io.WriteString(w, resource)
	}

	list, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", httpPort))
	if err != nil {
		t.Fatal("listening on port 5002", err)
	}

	s := &httptest.Server{
		Listener: list,
		Config:   &http.Server{Handler: http.HandlerFunc(hf)},
	}
	s.Start()
	defer s.Close()

	if err := cli.ChallengeReady(priv, chal); err != nil {
		t.Fatal(err)
	}

	certKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	template := &x509.CertificateRequest{
		SignatureAlgorithm: x509.SHA256WithRSA,
		PublicKeyAlgorithm: x509.RSA,
		PublicKey:          &certKey.PublicKey,
		Subject:            pkix.Name{CommonName: testDomain},
		DNSNames:           []string{testDomain},
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, template, certKey)
	if err != nil {
		t.Fatal(err)
	}
	csr, err := x509.ParseCertificateRequest(csrDER)
	if err != nil {
		t.Fatal(err)
	}

	certResp, err := cli.NewCertificate(priv, csr)
	if err != nil {
		t.Fatal(err)
	}
	contains := func(sli []string, ele string) bool {
		for _, e := range sli {
			if ele == e {
				return true
			}
		}
		return false
	}
	if !contains(certResp.Certificate.DNSNames, testDomain) {
		t.Errorf("returned cert was not for test domain")
	}

	certPEM := pemEncode(certResp.Certificate.Raw, "CERTIFICATE")
	certKeyPEM := pemEncode(x509.MarshalPKCS1PrivateKey(certKey), "RSA PRIVATE KEY")
	if _, err := tls.X509KeyPair(certPEM, certKeyPEM); err != nil {
		t.Errorf("private key did not match returned cert")
	}

	if !certResp.IsAvailable() {
		t.Error("Expected certificate to be available in CertificateResponse")
	}

	if certResp.Issuer == "" {
		t.Error("Expected issuer to be non empty.")
	}

	pemBundle, err := cli.Bundle(certResp)
	if err != nil {
		t.Errorf("Expected bundling certificate to return no errors. Error: %s", err)
	}

	parsed, err := parsePEMBundle(pemBundle)
	if err != nil {
		t.Errorf("Expected parsePEMBundle to return no errors. Error: %s", err)
	}

	if len(parsed) != 2 {
		t.Errorf("Expected bundled response to have two certificates. Given: %d", len(parsed))
	}

	if !reflect.DeepEqual(parsed[0].Raw, certResp.Certificate.Raw) {
		t.Error("Expected first certificate in bundle to match original certificate")
	}

	if parsed[1].IsCA != true {
		t.Error("Expected second certificate in bundle to be CA")
	}

	// Revoke
	if err := cli.RevokeCertificate(priv, certPEM); err != nil {
		t.Errorf("Revoke certificate failed. Error: %s", err)
	}
}

func TestRenewCertificate(t *testing.T) {
	cli, err := NewClient(testURL)
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	block, _ := pem.Decode(pemCert)
	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("Parsing pem certificate for test failed. Error: %s", err)
	}

	endpoint := server.URL + "/acme/cert/asdfouadfs"
	mux.HandleFunc("/acme/cert/asdfouadfs", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("fail") != "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/pkix-cert")
		w.Header().Set("Link", fmt.Sprintf(`<%s>;rel="up"`, endpoint))
		w.WriteHeader(http.StatusOK)
		w.Write(x509Cert.Raw)
	})

	certResp, err := cli.RenewCertificate(endpoint)
	if !certResp.IsAvailable() {
		t.Error("Expected RenewCertificate would return a certificate")
	}

	if !reflect.DeepEqual(certResp.Certificate.Raw, x509Cert.Raw) {
		t.Error("Expected RenewCertificate to return the same certificate")
	}

	if certResp.RetryAfter != 0 {
		t.Errorf("Expected RetryAfter to equal %d, given %d", 0, certResp.RetryAfter)
	}

	if certResp.Issuer != endpoint {
		t.Errorf("Expected issuer to equal %q, given %q", endpoint, certResp.Issuer)
	}

	// Test wrong status code returns an error
	if _, err := cli.RenewCertificate(endpoint + "?fail=1"); err == nil {
		t.Error("Expected invalid status code to return error")
	}
}

func TestParseLinks(t *testing.T) {
	tests := []struct {
		header http.Header
		want   map[string]string
	}{
		{
			header: map[string][]string{
				"Link": {
					`<https://example.com/acme/new-authz>;rel="next"`,
					`<https://example.com/acme/recover-reg>;rel="recover"`,
					`<https://example.com/acme/terms>;rel="terms-of-service"`,
				},
			},
			want: map[string]string{
				"next":             "https://example.com/acme/new-authz",
				"recover":          "https://example.com/acme/recover-reg",
				"terms-of-service": "https://example.com/acme/terms",
			},
		},
		{
			header: map[string][]string{
				"Link": []string{`<https://example.com/acme/new-cert>;rel="next"`},
			},
			want: map[string]string{
				"next": "https://example.com/acme/new-cert",
			},
		},
		{
			header: map[string][]string{
				"Link": {
					`<https://example.com/acme/ca-cert>;rel="up";title="issuer"`,
					`<https://example.com/acme/revoke-cert>;rel="revoke"`,
					`<https://example.com/acme/reg/asdf>;rel="author"`,
				},
				"Location":         {"https://example.com/acme/cert/asdf"},
				"Content-Location": {"https://example.com/acme/cert-seq/12345"},
			},
			want: map[string]string{
				"up":     "https://example.com/acme/ca-cert",
				"revoke": "https://example.com/acme/revoke-cert",
				"author": "https://example.com/acme/reg/asdf",
			},
		},
	}

	for i, test := range tests {
		links := parseLinks(test.header["Link"])

		for key, want := range test.want {
			given, ok := links[key]
			if !ok {
				t.Errorf("TestParseLinks (%d): want rel of %q to be present", i, key)
			}

			if given != want {
				t.Errorf("TestParseLinks (%d): want rel of %q to equal %s, given %s", i, key, want, given)
			}
		}
	}
}

func requiresEtcHostsEdits(t *testing.T) {
	addrs, err := net.LookupHost(testDomain)
	if err != nil || len(addrs) == 0 || (addrs[0] != "127.0.0.1" && addrs[0] != "::1") {
		addr := "NXDOMAIN"
		if len(addrs) > 0 {
			addr = addrs[0]
		}
		t.Skipf("/etc/hosts file not properly configured, skipping test. see README for required edits. %s resolved to %s", testDomain, addr)
	}
	return
}
