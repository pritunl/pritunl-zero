package letsencrypt

import (
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var pemCert = []byte(`-----BEGIN CERTIFICATE-----
MIIEVzCCAz+gAwIBAgITAP+xdyV4gP42OOW1GTK5/Dqh+zANBgkqhkiG9w0BAQsF
ADAfMR0wGwYDVQQDExRoYXBweSBoYWNrZXIgZmFrZSBDQTAeFw0xNTEyMTQyMjM4
MDBaFw0xNjAzMTMyMjM4MDBaMBYxFDASBgNVBAMTC2V4YW1wbGUub3JnMIIBIjAN
BgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyESSzaSOZeLoeVBABygcKMUuzKNR
SgdmNWP7GsjgNFmj+eiz/nia+BNmtSuqHK47YBLyCzqiAmt3bJt55vjnYx8Hra2z
W5TVwhpyzlOM0sbC0NzWOogDKC5woIEZGqdUTkaVEkIELQvNtOCjjjiIilk2A0K6
WLdFFNQWqLSnTYGCo53Gh5RIZrIqQEUbgJ7MJemj6SuPuCyussrF+WtvaBb7xQjN
LuVvprVv7NBuh8uz2cRuBLAsR03Fng1ItfuEtuLoKJC7Vfh2qBqwuEFZ8WoPP0H5
HMXK8pLJ5I4jmpdLeOik75yfhmxcJjdT57SgApSlU5XidOXmmuUbuk/8vQIDAQAB
o4IBkzCCAY8wDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggr
BgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBRMn3zVIo3wd1JgV+leFojm
6JCgWzAfBgNVHSMEGDAWgBT7eE8S+WAVgyyfF380GbMuNupBiTBqBggrBgEFBQcB
AQReMFwwJgYIKwYBBQUHMAGGGmh0dHA6Ly8xMjcuMC4wLjE6NDAwMi9vY3NwMDIG
CCsGAQUFBzAChiZodHRwOi8vMTI3LjAuMC4xOjQwMDAvYWNtZS9pc3N1ZXItY2Vy
dDAWBgNVHREEDzANggtleGFtcGxlLm9yZzAnBgNVHR8EIDAeMBygGqAYhhZodHRw
Oi8vZXhhbXBsZS5jb20vY3JsMGMGA1UdIARcMFowCgYGZ4EMAQIBMAAwTAYDKgME
MEUwIgYIKwYBBQUHAgEWFmh0dHA6Ly9leGFtcGxlLmNvbS9jcHMwHwYIKwYBBQUH
AgIwEwwRRG8gV2hhdCBUaG91IFdpbHQwDQYJKoZIhvcNAQELBQADggEBAKOOjKpt
xYUh+Ttun2OLR0RU7vk9wYvQPc8LpqpjqUQROMNzVQ9fO6im/Em0oTVNtaX+h5QU
493MCCyhkspajU8sTULA9f2l6Et7e03JO1K4lIc7hDDYOvlqwLZJ0+71OBjRnYPT
9KRmq6fizfrvvpmyBIQKZkTeCWO/9IQahMgnvpSDWXvVtjQPqgJg8vsIWHG+yC6W
RdXRo8GegNPrmc3wWV8mFDQ09j0ordRODtmnbl4ltiR7GKqPOkVck/hVTGZAm4KF
uD31pj4Nn3CPhxPbjMw9LsAHLqC+N1g6B2mog9uHLZxB1A1i3h8mKqz6YuzIrRj1
QLn0CFtToNBm2vY=
-----END CERTIFICATE-----`)

// TestRetry sets up a test server. The first request should return a
// Retry-After of 120 seconds with a 202 status code. The second
// request should return a certificate.
func TestRetry(t *testing.T) {
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

	delay := true
	mux.HandleFunc("/acme/cert/asdfouadfs", func(w http.ResponseWriter, r *http.Request) {
		if delay {
			w.Header().Set("Retry-After", "120")
			w.WriteHeader(http.StatusAccepted)

			delay = false
			return
		}

		w.Header().Set("Content-Type", "application/pkix-cert")
		w.WriteHeader(http.StatusOK)
		w.Write(x509Cert.Raw)
	})

	certResp := &CertificateResponse{
		URI:        server.URL + "/acme/cert/asdfouadfs",
		RetryAfter: 10,
	}

	cli.Retry(certResp)
	if certResp.IsAvailable() {
		t.Error("Expected first call to Retry would not return a certificate")
	}

	if certResp.RetryAfter != 120 {
		t.Errorf("Expected RetryAfter of %d, given  %d", 120, certResp.RetryAfter)
	}

	cli.Retry(certResp)
	if !certResp.IsAvailable() {
		t.Error("Expected second call to retry would return a certificate")
	}

	if !reflect.DeepEqual(certResp.Certificate.Raw, x509Cert.Raw) {
		t.Error("Expected certificate returned from Retry to be the same")
	}

	if certResp.RetryAfter != 0 {
		t.Errorf("Expected RetryAfter to equal %d, given %d", 0, certResp.RetryAfter)
	}
}
