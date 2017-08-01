package logger

import (
	"github.com/Sirupsen/logrus"
)

type ErrorWriter struct {
	Message string
	Fields  logrus.Fields
}

func (w *ErrorWriter) Write(input []byte) (n int, err error) {
	n = len(input)

	w.Fields["err"] = string(input)
	logrus.WithFields(w.Fields).Error(w.Message)

	return
}
