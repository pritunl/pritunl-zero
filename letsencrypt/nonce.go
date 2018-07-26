package letsencrypt

import (
	"errors"
	"net/http"
	"sync"
)

// All valid responses by the ACME server are required to have
// a Replay-Nonce header.
// https://letsencrypt.github.io/acme-spec/#replay-protection
var errNoNonce = errors.New("no Replay-Nonce header in HTTP response")

// nonceRoundTripper is a round tripper which
type nonceRoundTripper struct {
	rt http.RoundTripper

	// a cache of unused nonce tokens for JWS objects
	mu     *sync.Mutex
	nonces []string
}

func newNonceRoundTripper(rt http.RoundTripper) *nonceRoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &nonceRoundTripper{rt: rt, mu: new(sync.Mutex)}
}

func (nrt *nonceRoundTripper) Nonce() (string, error) {
	nrt.mu.Lock()
	defer nrt.mu.Unlock()

	n := len(nrt.nonces)
	if n == 0 {
		// TODO: determine a good strategy for replenishing the nonce cache.
		return "", errors.New("acme: nonce cache depleted")
	}
	// pop a nonce from the stack
	nonce := nrt.nonces[n-1]
	nrt.nonces = nrt.nonces[:n-1]
	return nonce, nil
}

func (nrt *nonceRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = nrt.rt.RoundTrip(req)
	if err != nil {
		return
	}
	tok := resp.Header.Get("Replay-Nonce")
	if tok == "" {
		return
	}
	nrt.mu.Lock()
	defer nrt.mu.Unlock()

	// cap the size of the nonce cache
	if len(nrt.nonces) < 2048 {
		nrt.nonces = append(nrt.nonces, tok)
	}

	return
}
