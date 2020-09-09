package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type ErrorWriter struct {
	Message string
	Fields  logrus.Fields
	Filters []string
}

func (w *ErrorWriter) Write(input []byte) (n int, err error) {
	n = len(input)

	inputStr := string(input)

	if w.Filters != nil {
		for _, filter := range w.Filters {
			if strings.Contains(inputStr, filter) {
				return
			}
		}
	}

	w.Fields["err"] = inputStr
	logrus.WithFields(w.Fields).Error(w.Message)

	return
}
