package twilio

import (
	"github.com/dropbox/godropbox/container/set"
)

var safeCharsPhone = set.NewSet(
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
	'+',
	'=',
	'_',
	'/',
	',',
	'.',
	':',
	'%',
	'@',
	'!',
	' ',
)

var safeCharsMessage = set.NewSet(
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
	'+',
	'=',
	'_',
	'/',
	',',
	'.',
	':',
	'#',
	'@',
	'%',
	'!',
	'[',
	']',
	'(',
	')',
	' ',
)

func filterStr(s string, n int, safe set.Set) string {
	if len(s) == 0 {
		return ""
	}

	if len(s) > n {
		s = s[:n]
	}

	ns := ""
	for _, c := range s {
		if safe.Contains(c) {
			ns += string(c)
		}
	}

	return ns
}

func FilterStrPhone(s string, n int) string {
	return filterStr(s, n, safeCharsPhone)
}

func FilterStrMessage(s string, n int) string {
	return filterStr(s, n, safeCharsMessage)
}
