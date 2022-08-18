package alertevent

import (
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
)

type Resource struct {
	Resource string `json:"resource" bson:"resource"`
	Level    int    `json:"level" bson:"level"`
	Value    int    `json:"value" bson:"value"`
}

func (s *Resource) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	switch s.Resource {
	case SystemHighMemory:
		if s.Value < 1 || s.Value > 100 {
			errData = &errortypes.ErrorData{
				Error:   "alert_resource_value_invalid",
				Message: "Alert resource value is invalid",
			}
			return
		}
		break
	case SystemHighSwap:
		if s.Value < 1 || s.Value > 100 {
			errData = &errortypes.ErrorData{
				Error:   "alert_resource_value_invalid",
				Message: "Alert resource value is invalid",
			}
			return
		}
		break
	case SystemHighHugePages:
		if s.Value < 1 || s.Value > 100 {
			errData = &errortypes.ErrorData{
				Error:   "alert_resource_value_invalid",
				Message: "Alert resource value is invalid",
			}
			return
		}
		break
	case DiskHighUsage:
		if s.Value < 1 || s.Value > 100 {
			errData = &errortypes.ErrorData{
				Error:   "alert_resource_value_invalid",
				Message: "Alert resource value is invalid",
			}
			return
		}
		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "alert_resource_name_invalid",
			Message: "Alert resource name is invalid",
		}
		return
	}

	switch s.Level {
	case Low, Medium, High:
		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "alert_resource_level_invalid",
			Message: "Alert resource level is invalid",
		}
		return
	}

	return
}
