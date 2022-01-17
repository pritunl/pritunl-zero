package logger

import (
	"os"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/sirupsen/logrus"
)

type fileSender struct{}

func (s *fileSender) Init() {}

func (s *fileSender) Parse(entry *logrus.Entry) {
	err := s.send(entry)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("logger: File send error")
	}
}

func (s *fileSender) send(entry *logrus.Entry) (err error) {
	msg := formatPlain(entry)

	file, err := os.OpenFile(constants.LogPath,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "logger: Failed to open log file"),
		}
		return
	}

	stat, err := file.Stat()
	if err != nil {
		_ = file.Close()
		err = &errortypes.ReadError{
			errors.Wrap(err, "logger: Failed to stat log file"),
		}
		return
	}

	if stat.Size() >= 5000000 {
		_ = os.Remove(constants.LogPath2)
		err = os.Rename(constants.LogPath, constants.LogPath2)
		if err != nil {
			_ = file.Close()
			err = &errortypes.WriteError{
				errors.Wrap(err, "logger: Failed to rotate log file"),
			}
			return
		}

		err = file.Close()
		if err != nil {
			err = &errortypes.WriteError{
				errors.Wrap(err, "logger: Failed to close log file"),
			}
			return
		}

		file, err = os.OpenFile(constants.LogPath,
			os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			err = &errortypes.WriteError{
				errors.Wrap(err, "logger: Failed to open log file"),
			}
			return
		}
	}

	_, err = file.Write(msg)
	if err != nil {
		_ = file.Close()
		err = &errortypes.WriteError{
			errors.Wrap(err, "logger: Failed to write to log file"),
		}
		return
	}

	err = file.Close()
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "logger: Failed to close log file"),
		}
		return
	}

	return
}

func init() {
	senders = append(senders, &fileSender{})
}
