package letsencrypt

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// CertificateResponse holds response items after requesting a Certificate.
type CertificateResponse struct {
	Certificate *x509.Certificate
	RetryAfter  int
	URI         string
	StableURI   string
	Issuer      string
}

// Bundle bundles the certificate with the issuer certificate.
func (c *Client) Bundle(certResp *CertificateResponse) (bundledPEM []byte, err error) {
	if !certResp.IsAvailable() {
		return nil, errors.New("Cannot bundle without certificate")
	}

	if certResp.Issuer == "" {
		return nil, errors.New("Could not bundle certificates. Issuer not found")
	}

	resp, err := c.client.Get(certResp.Issuer)
	if err != nil {
		return nil, fmt.Errorf("Error requesting issuer certificate: %s", err)
	}
	defer resp.Body.Close()

	issuerDER, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading issuer certificate: %s", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certResp.Certificate.Raw})
	issuerPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: issuerDER})

	return append(certPEM, issuerPEM...), nil
}

// Retry request retries the certificate if it was unavailable when calling
// NewCertificate or RenewCertificate.
//
// Note: If you are renewing a certificate, LetsEncrypt may return the same
// certificate. You should load your current x509.Certificate and use the
// Equal method to compare to the "new" certificate. If it's identical,
// you'll need to request a new certificate using NewCertificate, or if your
// chalenges have expired, start a new certificate flow entirely.
func (c *Client) Retry(certResp *CertificateResponse) error {
	if certResp.IsAvailable() {
		return errors.New("Aborting retry request. Certificate is already available")
	}

	if certResp.URI == "" {
		return errors.New("Could not make retry request. No URI available")
	}

	resp, err := c.client.Get(certResp.URI)
	if err != nil {
		return fmt.Errorf("Error retrying certificate request: %s", err)
	}
	defer resp.Body.Close()

	// Certificate is available
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %s", err)
		}

		x509Cert, err := x509.ParseCertificate(body)
		if err != nil {
			return fmt.Errorf("Error parsing x509 certificate: %s", err)
		}

		certResp.Certificate = x509Cert
		certResp.RetryAfter = 0

		if stableURI := resp.Header.Get("Content-Location"); stableURI != "" {
			certResp.StableURI = stableURI
		}

		links := parseLinks(resp.Header["Link"])
		certResp.Issuer = links["up"]

		return nil
	}

	// Certificate still isn't ready.
	if resp.StatusCode == http.StatusAccepted {
		retryAfter, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			return fmt.Errorf("Error parsing retry-after header: %s", err)
		}

		certResp.RetryAfter = retryAfter

		return nil
	}

	return fmt.Errorf("Retry expected status code of %d or %d, given %d", http.StatusOK, http.StatusAccepted, resp.StatusCode)
}

// IsAvailable returns bool true if CertificateResponse has a certificate
// available. It's a convenience function, but it helps with readability.
func (c *CertificateResponse) IsAvailable() bool {
	return c.Certificate != nil
}
