package letsencrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Specified in boulder's configuration
// See $GOATH/src/github.com/letsencrypt/boulder/test/boulder-config.json
var (
	httpPort      int    = 5002
	httpsPort     int    = 5001
	boulderDNSSrv string = "http://localhost:8055/set-txt"
)

func TestHTTPChallenge(t *testing.T) {
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
}

func TestTLSSNIChallenge(t *testing.T) {
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

	chals := auth.Combinations(ChallengeTLSSNI)
	if len(chals) == 0 || len(chals[0]) != 1 {
		t.Fatal("no supported challenges")
	}
	chal := chals[0][0]
	certs, err := chal.TLSSNI(priv)
	if err != nil {
		t.Fatal(err)
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{},
		GetCertificate: func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert, ok := certs[clientHello.ServerName]
			if ok {
				return cert, nil
			}
			t.Errorf("got unknown SNI server name: %v", clientHello.ServerName)
			return nil, nil
		},
	}
	for _, cert := range certs {
		tlsConf.Certificates = append(tlsConf.Certificates, *cert)
	}

	list, err := tls.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", httpsPort), tlsConf)
	if err != nil {
		t.Errorf("listening on port %d: %v", httpsPort, err)
		return
	}
	defer list.Close()
	go func() {
		for {
			conn, err := list.Accept()
			if err != nil {
				return
			}
			if conn, ok := conn.(*tls.Conn); ok {
				// must get past the handshake
				if err := conn.Handshake(); err != nil {
					t.Errorf("handshake error: %v", err)
				}
			} else {
				t.Errorf("connection is not a tls connection")
			}
			conn.Close()
		}
	}()

	if err := cli.ChallengeReady(priv, chal); err != nil {
		t.Fatal(err)
	}
}

func TestDNSChallenge(t *testing.T) {
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

	chals := auth.Combinations(ChallengeDNS)
	if len(chals) == 0 || len(chals[0]) != 1 {
		t.Fatal("no supported challenges")
	}
	chal := chals[0][0]
	subdomain, txt, err := chal.DNS(priv)
	if err != nil {
		t.Fatal(err)
	}

	body := struct {
		Host  string `json:"host"`
		Value string `json:"value"`
	}{
		// end host in a period so its fqdn for dns question
		Host:  strings.Join([]string{subdomain, testDomain, ""}, "."),
		Value: txt,
	}
	bodyb, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", boulderDNSSrv, bytes.NewReader(bodyb))
	if err != nil {
		t.Fatal(err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if err := cli.ChallengeReady(priv, chal); err != nil {
		t.Fatal(err)
	}
}
