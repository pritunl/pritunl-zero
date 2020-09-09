package logger

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/log"
)

var (
	databaseBuffer = make(chan *logrus.Entry, 128)
)

type databaseSender struct{}

func (s *databaseSender) Init() {}

func (s *databaseSender) Parse(entry *logrus.Entry) {
	if len(buffer) <= 32 {
		databaseBuffer <- entry
	}
}

func databaseSend(entry *logrus.Entry) (err error) {
	level := ""

	db := database.GetDatabase()
	if db == nil {
		return
	}
	defer db.Close()

	switch entry.Level {
	case logrus.DebugLevel:
		level = log.Debug
		break
	case logrus.WarnLevel:
		level = log.Warning
		break
	case logrus.InfoLevel:
		level = log.Info
		break
	case logrus.ErrorLevel:
		level = log.Error
		break
	case logrus.FatalLevel:
		level = log.Fatal
		break
	case logrus.PanicLevel:
		level = log.Panic
		break
	default:
		level = log.Unknown
	}

	ent := &log.Entry{
		Level:     level,
		Timestamp: entry.Time,
		Message:   entry.Message,
		Fields:    map[string]interface{}{},
	}

	for key, val := range entry.Data {
		if key == "error" {
			ent.Stack = fmt.Sprintf("%s", val)
		} else {
			ent.Fields[key] = val
		}
	}

	err = ent.Insert(db)
	if err != nil {
		return
	}

	return
}

func initDatabaseSender() {
	go func() {
		for {
			entry := <-databaseBuffer

			if constants.Interrupt {
				return
			}

			if strings.HasPrefix(entry.Message, "logger:") {
				continue
			}

			err := databaseSend(entry)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("logger: Database send error")
			}
		}
	}()
}

func init() {
	senders = append(senders, &databaseSender{})
}
