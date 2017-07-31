package logger

import (
	"github.com/Sirupsen/logrus"
)

type ErrorWriter struct {
	Prefix string
	Fields logrus.Fields
}

func (w *ErrorWriter) Write(input []byte) (n int, err error) {
	n = len(input)

	logrus.WithFields(w.Fields).Error(w.Prefix + string(input))

	return
}
