package utils

import (
	"strings"

	"github.com/dropbox/godropbox/container/set"
)

const nameSafeLimit = 128

var nameSafeChar = set.NewSet(
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'-', '+', '=', '_', '/', ',', '.', '@', '#',
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

var idSafeChar = set.NewSet(
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', '-', '_',
)

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

func FilterDomain(s string) string {
	s = FilterName(strings.ToLower(s))
	s = strings.TrimPrefix(s, ".")
	s = strings.TrimSuffix(s, ".")
	return s
}
