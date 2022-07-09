package alert

import (
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
)

type Resource struct {
	Name  string `json:"name" bson:"name"`
	Level int    `json:"level" bson:"level"`
}

func (s *Resource) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if s.Name == "" {
		errData = &errortypes.ErrorData{
			Error:   "alert_source_name_invalid",
			Message: "Alert source name is invalid",
		}
		return
	}

	switch s.Level {
	case Low:
	case Medium:
	case High:
		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "alert_source_level_invalid",
			Message: "Alert source level is invalid",
		}
		return
	}

	return
}
