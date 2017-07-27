package logger

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/log"
)

type databaseSender struct{}

func (s *databaseSender) Init() {}

func (s *databaseSender) Parse(entry *logrus.Entry) {
	err := s.send(entry)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("logger: Database send error")
	}
}

func (s *databaseSender) send(entry *logrus.Entry) (err error) {
	level := ""

	db := database.GetDatabase()
	if db == nil {
		// TODO Defer to file
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

func init() {
	senders = append(senders, &databaseSender{})
}
