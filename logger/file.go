package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/dropbox/godropbox/errors"
	"os"
	"path"
	"sync"
)

var fileLock = sync.Mutex{}

type fileSender struct{}

func (s *fileSender) Init() {}

func (s *fileSender) Parse(entry *logrus.Entry) {
	msg := formatPlain(entry)

	fileLock.Lock()
	defer fileLock.Unlock()

	pth := path.Join("pritunl-zero.log")

	file, err := os.OpenFile(pth, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "logger: Failed to write entry"),
		}
		return
	}
	defer file.Close()

	file.Write(msg)
}

func init() {
	senders = append(senders, &fileSender{})
}
