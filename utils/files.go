package utils

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
)

var invalidPaths = set.NewSet("/", "", ".", "./")

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
