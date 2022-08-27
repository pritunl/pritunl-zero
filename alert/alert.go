package alert

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/utils"
)

type Alert struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Roles     []string           `bson:"roles" json:"roles"`
	Resource  string             `bson:"resource" json:"resource"`
	Level     int                `bson:"level" json:"level"`
	Frequency int                `bson:"frequency" json:"frequency"`
	Ignores   []string           `bson:"ignores" json:"ignores"`
	ValueInt  int                `bson:"value_int" json:"value_int"`
	ValueStr  string             `bson:"value_str" json:"value_str"`
}

func (a *Alert) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if a.Id.IsZero() {
		a.Id, err = utils.RandObjectId()
		if err != nil {
			return
		}
	}

	if a.Roles == nil {
		a.Roles = []string{}
	}

	if a.Frequency == 0 {
		a.Frequency = 300
	}

	if a.Frequency < 300 {
		errData = &errortypes.ErrorData{
			Error:   "alert_frequency_invalid",
			Message: "Alert frequency cannot be less then 300 seconds",
		}
		return
	}

	if a.Frequency > 604800 {
		errData = &errortypes.ErrorData{
			Error:   "alert_frequency_invalid",
			Message: "Alert frequency too large",
		}
		return
	}

	if a.Ignores != nil {
		a.Ignores = []string{}
	}

	switch a.Resource {
	case SystemCpuLevel:
		if a.ValueInt < 1 || a.ValueInt > 100 {
			errData = &errortypes.ErrorData{
				Error:   "alert_value_invalid",
				Message: "Alert value is invalid",
			}
			return
		}
		a.ValueStr = ""
		break
	case SystemMemoryLevel:
		if a.ValueInt < 1 || a.ValueInt > 100 {
			errData = &errortypes.ErrorData{
				Error:   "alert_value_invalid",
				Message: "Alert value is invalid",
			}
			return
		}
		a.ValueStr = ""
		break
	case SystemSwapLevel:
		if a.ValueInt < 1 || a.ValueInt > 100 {
			errData = &errortypes.ErrorData{
				Error:   "alert_value_invalid",
				Message: "Alert value is invalid",
			}
			return
		}
		a.ValueStr = ""
		break
	case SystemHugePagesLevel:
		if a.ValueInt < 1 || a.ValueInt > 100 {
			errData = &errortypes.ErrorData{
				Error:   "alert_value_invalid",
				Message: "Alert value is invalid",
			}
			return
		}
		a.ValueStr = ""
		break
	case DiskUsageLevel:
		if a.ValueInt < 1 || a.ValueInt > 100 {
			errData = &errortypes.ErrorData{
				Error:   "alert_value_invalid",
				Message: "Alert value is invalid",
			}
			return
		}
		a.ValueStr = ""
		break
	case KmsgKeyword:
		if a.ValueStr == "" {
			errData = &errortypes.ErrorData{
				Error:   "alert_value_invalid",
				Message: "Alert value is invalid",
			}
			return
		}
		a.ValueInt = 0
		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "alert_resource_name_invalid",
			Message: "Alert resource name is invalid",
		}
		return
	}

	switch a.Level {
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

func (a *Alert) Commit(db *database.Database) (err error) {
	coll := db.Alerts()

	err = coll.Commit(a.Id, a)
	if err != nil {
		return
	}

	return
}

func (a *Alert) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Alerts()

	err = coll.CommitFields(a.Id, a, fields)
	if err != nil {
		return
	}

	return
}

func (a *Alert) Insert(db *database.Database) (err error) {
	coll := db.Alerts()

	_, err = coll.InsertOne(db, a)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
