package utils

import (
	"github.com/dropbox/godropbox/container/set"
	"strings"
)

const nameSafeLimit = 128

var nameSafeChar = set.NewSet(
	'a',
	'b',
	'c',
	'd',
	'e',
	'f',
	'g',
	'h',
	'i',
	'j',
	'k',
	'l',
	'm',
	'n',
	'o',
	'p',
	'q',
	'r',
	's',
	't',
	'u',
	'v',
	'w',
	'x',
	'y',
	'z',
	'A',
	'B',
	'C',
	'D',
	'E',
	'F',
	'G',
	'H',
	'I',
	'J',
	'K',
	'L',
	'M',
	'N',
	'O',
	'P',
	'Q',
	'R',
	'S',
	'T',
	'U',
	'V',
	'W',
	'X',
	'Y',
	'Z',
	'0',
	'1',
	'2',
	'3',
	'4',
	'5',
	'6',
	'7',
	'8',
	'9',
	'-',
)

func FilterName(s string) string {
	if len(s) == 0 {
		return ""
	}

	if len(s) > nameSafeLimit {
		s = s[:nameSafeLimit]
	}

	var ns strings.Builder
	for _, c := range s {
		if safeChars.Contains(c) {
			ns.WriteString(string(c))
		}
	}

	return ns.String()
}
