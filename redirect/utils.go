package main

import (
	"strings"
)

func StripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}

	n := strings.Count(hostport, ":")
	if n > 1 {
		if i := strings.IndexByte(hostport, ']'); i != -1 {
			return strings.TrimPrefix(hostport[:i], "[")
		}
		return hostport
	}

	return hostport[:colon]
}
