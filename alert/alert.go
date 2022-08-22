package alert

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/utils"
)

type Alert struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Roles    []string           `bson:"roles" json:"roles"`
	Resource string             `bson:"resource" json:"resource"`
	Level    int                `bson:"level" json:"level"`
	Value    int                `bson:"value" json:"value"`
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

	switch a.Resource {
	case SystemHighMemory:
		if a.Value < 1 || a.Value > 100 {
			errData = &errortypes.ErrorData{
				Error:   "alert_resource_value_invalid",
				Message: "Alert resource value is invalid",
			}
			return
		}
		break
	case SystemHighSwap:
		if a.Value < 1 || a.Value > 100 {
			errData = &errortypes.ErrorData{
				Error:   "alert_resource_value_invalid",
				Message: "Alert resource value is invalid",
			}
			return
		}
		break
	case SystemHighHugePages:
		if a.Value < 1 || a.Value > 100 {
			errData = &errortypes.ErrorData{
				Error:   "alert_resource_value_invalid",
				Message: "Alert resource value is invalid",
			}
			return
		}
		break
	case DiskHighUsage:
		if a.Value < 1 || a.Value > 100 {
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