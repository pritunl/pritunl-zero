package dns

import (
	"net"
	"strings"
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

var multiLevelTlds = []string{
	"co.za", "co.uk", "com.au", "co.nz", "co.jp", "or.jp",
	"com.br", "co.in", "com.sg", "co.kr", "com.mx",
	"com.cn", "org.cn", "net.cn", "ac.uk", "gov.uk",
}

func extractDomain(domain string) string {
	domain = strings.Trim(domain, ".")
	parts := strings.Split(domain, ".")
	// Check for multi-level TLDs (e.g., .co.za, .co.uk)
	if len(parts) >= 3 {
		lastTwo := parts[len(parts)-2] + "." + parts[len(parts)-1]
		for _, tld := range multiLevelTlds {
			if lastTwo == tld {
				// Return last 3 parts for multi-level TLDs
				return parts[len(parts)-3] + "." + lastTwo
			}
		}
	}
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "." + parts[len(parts)-1]
	}
	return domain
}

func cleanDomain(domain string) string {
	return strings.Trim(domain, ".")
}
