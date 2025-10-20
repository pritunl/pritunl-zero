package utils

import (
	"container/list"
	"io/ioutil"
	"os/exec"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/sirupsen/logrus"
)

var isSystemd *bool

var (
	safeChars = set.NewSet(
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
		'~',
		'@',
		'#',
		'!',
		'&',
		' ',
	)
	unixSafeChars = set.NewSet(
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
		'_',
		'.',
	)
)

func FilterStr(s string, n int) string {
	if len(s) == 0 {
		return ""
	}

	if len(s) > n {
		s = s[:n]
	}

	var ns strings.Builder
	for _, c := range s {
		if safeChars.Contains(c) {
			ns.WriteString(string(c))
		}
	}

	return ns.String()
}

func FilterUnixStr(s string, n int) string {
	if len(s) == 0 {
		return ""
	}

	if len(s) > n {
		s = s[:n]
	}

	var ns strings.Builder
	for _, c := range s {
		if unixSafeChars.Contains(c) {
			ns.WriteString(string(c))
		}
	}

	return ns.String()
}

func PointerBool(x bool) *bool {
	return &x
}

func PointerInt(x int) *int {
	return &x
}

func PointerString(x string) *string {
	return &x
}

func SinceAbs(t time.Time) (s time.Duration) {
	s = time.Since(t)
	if s < 0 {
		s = s * -1
	}
	return
}

func IsSystemd() bool {
	if isSystemd != nil {
		return *isSystemd
	}

	data, err := ioutil.ReadFile("/proc/1/cmdline")
	if err == nil {
		parts := strings.Split(string(data), "\x00")
		if len(parts) > 0 && strings.Contains(
			strings.ToLower(parts[0]), "systemd") {

			isSysd := true
			isSystemd = &isSysd
			return true
		}
	}

	data, err = ioutil.ReadFile("/proc/1/comm")
	if err == nil {
		if strings.Contains(strings.ToLower(string(data)), "systemd") {
			isSysd := true
			isSystemd = &isSysd
			return true
		}
	}

	cmd := exec.Command("ps", "-p", "1", "-o", "comm=")
	output, err := cmd.Output()
	if err == nil {
		if strings.Contains(strings.ToLower(string(output)), "systemd") {
			isSysd := true
			isSystemd = &isSysd
			return true
		}
	}

	isSysd := false
	isSystemd = &isSysd
	return false
}

func RecoverLog() {
	panc := recover()
	if panc != nil {
		logrus.WithFields(logrus.Fields{
			"trace": string(debug.Stack()),
			"panic": panc,
		}).Error("sync: Panic in goroutine")
	}
}

func CopyList(src *list.List) *list.List {
	dst := list.New()
	for x := src.Front(); x != nil; x = x.Next() {
		dst.PushBack(x.Value)
	}
	return dst
}

func GetIntVer(version string) int {
	re := regexp.MustCompile(`\d+`)
	ver := re.FindAllString(version, -1)

	if len(ver) == 0 {
		return 0
	}

	lastNum, err := strconv.Atoi(ver[len(ver)-1])
	if err != nil {
		return 0
	}
	ver[len(ver)-1] = strconv.Itoa(lastNum + 4000)

	var builder strings.Builder
	for _, v := range ver {
		num, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		builder.WriteString(strings.Repeat(
			"0", 4-len(strconv.Itoa(num))) + strconv.Itoa(num))
	}

	result, err := strconv.Atoi(builder.String())
	if err != nil {
		return 0
	}

	return result
}
