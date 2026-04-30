package dns

import (
	"net"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func matchDomains(x, y string) bool {
	if strings.Trim(x, ".") == strings.Trim(y, ".") {
		return true
	}
	return false
}

func matchTxt(x, y string) bool {
	if strings.Trim(x, "\"") == strings.Trim(y, "\"") {
		return true
	}
	return false
}

func normalizeIp(addr string) string {
	ip := net.ParseIP(addr)
	if ip == nil {
		return ""
	}

	return strings.ToLower(ip.String())
}

func extractDomain(domain string) string {
	domain = strings.Trim(domain, ".")
	topDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return domain
	}
	return topDomain
}

func cleanDomain(domain string) string {
	return strings.Trim(domain, ".")
}
