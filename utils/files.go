package utils

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
)

var invalidPaths = set.NewSet("/", "", ".", "./")

const pathSafeLimit = 256

var pathSafeChars = set.NewSet(
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
	'+',
	'=',
	'@',
	'/',
)

func FilterPath(pth string) string {
	if len(pth) > pathSafeLimit {
		pth = pth[:pathSafeLimit]
	}

	cleaned := ""
	for _, c := range pth {
		if pathSafeChars.Contains(c) {
			cleaned += string(c)
		}
	}

	cleaned = filepath.Clean(cleaned)
	cleaned, err := filepath.Abs(cleaned)
	if err != nil {
		return ""
	}
	cleaned = filepath.FromSlash(cleaned)
	cleaned = strings.ReplaceAll(cleaned, "..", "")

	return cleaned
}

func FilterRelPath(pth string) string {
	if len(pth) > pathSafeLimit {
		pth = pth[:pathSafeLimit]
	}

	cleaned := ""
	for _, c := range pth {
		if pathSafeChars.Contains(c) {
			cleaned += string(c)
		}
	}

	cleaned = filepath.Clean(cleaned)
	cleaned = filepath.FromSlash(cleaned)
	cleaned = strings.ReplaceAll(cleaned, "..", "")

	return cleaned
}

func Chmod(pth string, mode os.FileMode) (err error) {
	err = os.Chmod(pth, mode)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "utils: Failed to chmod %s", pth),
		}
		return
	}

	return
}

func Exists(pth string) (exists bool, err error) {
	_, err = os.Stat(pth)
	if err == nil {
		exists = true
		return
	}

	if os.IsNotExist(err) {
		err = nil
		return
	}

	err = &errortypes.ReadError{
		errors.Wrapf(err, "utils: Failed to stat %s", pth),
	}
	return
}

func ExistsDir(pth string) (exists bool, err error) {
	stat, err := os.Stat(pth)
	if err == nil {
		exists = stat.IsDir()
		return
	}

	if os.IsNotExist(err) {
		err = nil
		return
	}

	err = &errortypes.ReadError{
		errors.Wrapf(err, "utils: Failed to stat %s", pth),
	}
	return
}

func ExistsFile(pth string) (exists bool, err error) {
	stat, err := os.Stat(pth)
	if err == nil {
		exists = !stat.IsDir()
		return
	}

	if os.IsNotExist(err) {
		err = nil
		return
	}

	err = &errortypes.ReadError{
		errors.Wrapf(err, "utils: Failed to stat %s", pth),
	}
	return
}

func ExistsMkdir(pth string, perm os.FileMode) (err error) {
	exists, err := ExistsDir(pth)
	if err != nil {
		return
	}

	if !exists {
		err = os.MkdirAll(pth, perm)
		if err != nil {
			err = &errortypes.WriteError{
				errors.Wrapf(err, "utils: Failed to mkdir %s", pth),
			}
			return
		}
	}

	return
}

func ExistsRemove(pth string) (err error) {
	exists, err := Exists(pth)
	if err != nil {
		return
	}

	if exists {
		err = os.RemoveAll(pth)
		if err != nil {
			err = &errortypes.WriteError{
				errors.Wrapf(err, "utils: Failed to rm %s", pth),
			}
			return
		}
	}

	return
}

func ReadExists(path string) (data string, err error) {
	dataByt, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}
		err = &errortypes.ReadError{
			errors.Wrapf(err, "utils: Failed to read '%s'", path),
		}
		return
	}

	data = string(dataByt)
	return
}

func Remove(path string) (err error) {
	if invalidPaths.Contains(path) {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "utils: Invalid remove path '%s'", path),
		}
		return
	}

	err = os.Remove(path)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "utils: Failed to remove '%s'", path),
		}
		return
	}

	return
}

func RemoveAll(path string) (err error) {
	if invalidPaths.Contains(path) {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "utils: Invalid remove path '%s'", path),
		}
		return
	}

	err = os.RemoveAll(path)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "utils: Failed to remove '%s'", path),
		}
		return
	}

	return
}

func RemoveWildcard(matchPath string) (n int, err error) {
	matches, err := filepath.Glob(matchPath)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "utils: Error matching path '%s'", matchPath),
		}
		return
	}

	if len(matches) == 0 {
		return
	}

	delErrors := []string{}
	for _, pth := range matches {
		fileInfo, err := os.Stat(pth)
		if err != nil {
			delErrors = append(delErrors, fmt.Sprintf("%s: %v", pth, err))
			continue
		}

		if fileInfo.IsDir() {
			continue
		}

		err = os.Remove(pth)
		if err != nil {
			delErrors = append(delErrors, fmt.Sprintf("%s: %v", pth, err))
		} else {
			n += 1
		}
	}

	if len(delErrors) > 0 {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "utils: Delete errors '%s'",
				strings.Join(delErrors, ",")),
		}
		return
	}

	return
}

func ContainsDir(pth string) (hasDir bool, err error) {
	exists, err := ExistsDir(pth)
	if !exists {
		return
	}

	entries, err := ioutil.ReadDir(pth)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "queue: Failed to read dir %s", pth),
		}
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			hasDir = true
			return
		}
	}

	return
}

func Open(path string, perm os.FileMode) (file *os.File, err error) {
	file, err = os.OpenFile(path, os.O_RDWR|os.O_TRUNC, perm)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "utils: Failed to open '%s'", path),
		}
		return
	}

	return
}

func Read(path string) (data string, err error) {
	dataByt, err := ioutil.ReadFile(path)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "utils: Failed to read '%s'", path),
		}
		return
	}

	data = string(dataByt)
	return
}

func ReadLines(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "utils: Failed to open '%s'", path),
		}
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			err = &errortypes.ReadError{
				errors.Wrapf(err, "utils: Failed to read '%s'", path),
			}
			return
		}
	}()

	lines = []string{}
	reader := bufio.NewReader(file)
	for {
		line, e := reader.ReadString('\n')
		if e != nil {
			break
		}
		lines = append(lines, strings.Trim(line, "\n"))
	}

	return
}

func Write(path string, data string, perm os.FileMode) (err error) {
	file, err := Open(path, perm)
	if err != nil {
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			err = &errortypes.WriteError{
				errors.Wrapf(err, "utils: Failed to write '%s'", path),
			}
			return
		}
	}()

	_, err = file.WriteString(data)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "utils: Failed to write to file '%s'", path),
		}
		return
	}

	return
}

func Create(path string, perm os.FileMode) (file *os.File, err error) {
	file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "utils: Failed to create '%s'", path),
		}
		return
	}

	return
}

func CreateWrite(path string, data string, perm os.FileMode) (err error) {
	file, err := Create(path, perm)
	if err != nil {
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			err = &errortypes.WriteError{
				errors.Wrapf(err, "utils: Failed to write '%s'", path),
			}
			return
		}
	}()

	_, err = file.WriteString(data)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "utils: Failed to write to file '%s'", path),
		}
		return
	}

	return
}

func FileSha256(pth string) (hash string, err error) {
	file, err := os.Open(pth)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "utils: Failed to read '%s'", pth),
		}
		return
	}
	defer file.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "utils: Failed to read '%s'", pth),
		}
		return
	}

	hash = fmt.Sprintf("%x", hasher.Sum(nil))
	return
}
