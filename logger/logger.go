package logger

import (
	"github.com/Sirupsen/logrus"
	"os"
)

var (
	buffer  = make(chan *logrus.Entry, 32)
	senders = []sender{}
)

func initSender() {
	for _, sndr := range senders {
		sndr.Init()
	}

	go func() {
		for {
			entry := <-buffer

			if len(entry.Message) > 7 && entry.Message[:7] == "logger:" {
				continue
			}

			for _, sndr := range senders {
				sndr.Parse(entry)
			}
		}
	}()
}

func Init() {
	initSender()

	logrus.SetFormatter(&formatter{})
	logrus.AddHook(&logHook{})
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.InfoLevel)
}
