package utils

import (
	"fmt"
	"net"
	"strings"
	"unicode/utf8"

	"github.com/dropbox/godropbox/container/set"
	"golang.org/x/crypto/ssh"
)

const (
	idSafeLimit         = 256
	nameSafeLimit       = 128
	ifaceSafeLimit      = 15
	domainSafeLimit     = 1024
	base64SafeLimit     = 32768
	pemSafeLimit        = 32768
	sshKeySafeLimit     = 32768
	sshCommentSafeLimit = 256
)

var nameSafeChar = set.NewSet(
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'-', '.',
)

var nameCmdSafeChar = set.NewSet(
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'-', '.',
)

var unitSafeChar = set.NewSet(
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	':', '-', '_', '.', '\\', '@',
)

var domainSafeChar = set.NewSet(
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'-', '.',
)

var idSafeChar = set.NewSet(
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', '-', '_',
)

var base64SafeChar = set.NewSet(
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'+', '/', '=', '-', '_',
)

var pemSafeChar = set.NewSet(
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'+', '/', '=', '-', ' ', '\n',
)

var ifaceSafeChar = set.NewSet(
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'-', '_', '.',
)

var sshCommentSafeChar = set.NewSet(
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'+', '/', '=', '-', '_', '.', '@', ' ',
)

func normalizeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

func FilterName(s string) string {
	if len(s) == 0 {
		return ""
	}

	if s == "self" {
		s = "invalid-name"
	}

	if len(s) > nameSafeLimit {
		s = s[:nameSafeLimit]
	}

	var ns strings.Builder
	for _, c := range s {
		if nameSafeChar.Contains(c) {
			ns.WriteString(string(c))
		}
	}

	return ns.String()
}

func FilterNameCmd(s string) string {
	if len(s) == 0 {
		return ""
	}

	if s == "self" {
		s = "invalid-name"
	}

	if len(s) > nameSafeLimit {
		s = s[:nameSafeLimit]
	}

	var ns strings.Builder
	for _, c := range s {
		if nameCmdSafeChar.Contains(c) {
			ns.WriteString(string(c))
		}
	}

	return strings.ToLower(ns.String())
}

func FilterUnit(s string) string {
	if len(s) == 0 {
		return ""
	}

	if len(s) > nameSafeLimit {
		s = s[:nameSafeLimit]
	}

	var ns strings.Builder
	for _, c := range s {
		if unitSafeChar.Contains(c) {
			ns.WriteString(string(c))
		}
	}

	return ns.String()
}

func FilterId(s string) string {
	if len(s) == 0 {
		return ""
	}

	if len(s) > nameSafeLimit {
		s = s[:nameSafeLimit]
	}

	var ns strings.Builder
	for _, c := range s {
		if idSafeChar.Contains(c) {
			ns.WriteString(string(c))
		}
	}

	return ns.String()
}

func FilterIface(s string) string {
	if len(s) == 0 {
		return ""
	}

	if len(s) > ifaceSafeLimit {
		s = s[:ifaceSafeLimit]
	}

	var ns strings.Builder
	for _, c := range s {
		if ifaceSafeChar.Contains(c) {
			ns.WriteString(string(c))
		}
	}

	return ns.String()
}

func FilterDomain(s string) string {
	if len(s) == 0 {
		return ""
	}

	if s == "self" {
		s = "invalid-name"
	}

	if len(s) > nameSafeLimit {
		s = s[:nameSafeLimit]
	}

	var ns strings.Builder
	for _, c := range s {
		if domainSafeChar.Contains(c) {
			ns.WriteString(string(c))
		}
	}

	s = ns.String()
	s = strings.TrimPrefix(s, ".")
	s = strings.TrimSuffix(s, ".")
	return s
}

func FilterAddress(s string) string {
	if len(s) == 0 {
		return ""
	}

	if s == "self" {
		s = "invalid-name"
	}

	if len(s) > domainSafeLimit {
		s = s[:domainSafeLimit]
	}

	s = strings.TrimSpace(s)

	if ip := net.ParseIP(strings.Trim(s, "[]")); ip != nil {
		return ip.String()
	}

	var ns strings.Builder
	for _, c := range s {
		if domainSafeChar.Contains(c) {
			ns.WriteString(string(c))
		}
	}

	s = ns.String()
	s = strings.TrimPrefix(s, ".")
	s = strings.TrimSuffix(s, ".")

	return s
}

func FilterBase64(s string) string {
	if len(s) == 0 {
		return ""
	}

	if len(s) > base64SafeLimit {
		s = s[:base64SafeLimit]
	}

	var ns strings.Builder
	for _, c := range s {
		if base64SafeChar.Contains(c) {
			ns.WriteString(string(c))
		}
	}

	return ns.String()
}

func FilterPem(s string) string {
	if len(s) == 0 {
		return ""
	}

	if len(s) > pemSafeLimit {
		s = s[:pemSafeLimit]
	}

	s = normalizeNewlines(s)

	var ns strings.Builder
	for _, c := range s {
		if pemSafeChar.Contains(c) {
			ns.WriteString(string(c))
		}
	}

	return ns.String()
}

func filterSshPublicKey(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return ""
	}

	key, comment, options, _, err := ssh.ParseAuthorizedKey([]byte(s))
	if err != nil || len(options) > 0 {
		return ""
	}

	if _, ok := key.(*ssh.Certificate); ok {
		return ""
	}

	line := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key)))

	if len(comment) > sshCommentSafeLimit {
		comment = comment[:sshCommentSafeLimit]
	}

	var ns strings.Builder
	for _, c := range comment {
		if sshCommentSafeChar.Contains(c) {
			ns.WriteString(string(c))
		}
	}

	comment = strings.TrimSpace(ns.String())
	if comment != "" {
		line += " " + comment
	}

	return line
}

func FilterSshKey(s string) string {
	if len(s) == 0 {
		return ""
	}

	if len(s) > sshKeySafeLimit {
		s = s[:sshKeySafeLimit]
	}

	s = normalizeNewlines(s)

	lines := []string{}
	for _, line := range strings.Split(s, "\n") {
		line = filterSshPublicKey(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func FilterSystemdEnv(key, value string) string {
	if key == "" || value == "" || !utf8.ValidString(value) {
		return ""
	}

	for i := 0; i < len(key); i++ {
		c := key[i]
		if c == '_' ||
			(c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(i > 0 && c >= '0' && c <= '9') {
			continue
		}
		return ""
	}

	var quoted strings.Builder
	quoted.WriteString(`Environment="`)
	quoted.WriteString(key)
	quoted.WriteByte('=')

	for _, c := range value {
		if c == 0 ||
			(c >= 0xfdd0 && c <= 0xfdef) ||
			(c&0xffff) >= 0xfffe {
			return ""
		}

		switch c {
		case '\\':
			quoted.WriteString(`\\`)
		case '"':
			quoted.WriteString(`\"`)
		case '%':
			quoted.WriteString("%%")
		case '\n':
			quoted.WriteString(`\n`)
		case '\r':
			quoted.WriteString(`\r`)
		case '\t':
			quoted.WriteString(`\t`)
		default:
			if c < 0x20 || c == 0x7f {
				fmt.Fprintf(&quoted, `\x%02x`, c)
			} else {
				quoted.WriteRune(c)
			}
		}
	}

	quoted.WriteByte('"')
	return quoted.String()
}
