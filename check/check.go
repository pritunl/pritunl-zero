package check

import (
	"strings"

	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/utils"
)

type Check struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Roles      []string           `bson:"roles" json:"roles"`
	Frequency  int                `bson:"frequency" json:"frequency"`
	Type       string             `bson:"type" json:"type"`
	Targets    []string           `bson:"targets" json:"targets"`
	Timeout    int                `bson:"timeout" json:"timeout"`
	Method     string             `bson:"method" json:"method"`
	StatusCode int                `bson:"status_code" json:"status_code"`
	Headers    []*Header          `bson:"headers" json:"headers"`
}

type Header struct {
	Key   string `bson:"key" json:"key"`
	Value string `bson:"value" json:"value"`
}

func (a *Check) Validate(db *database.Database) (
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
		a.Frequency = 30
	}

	if a.Frequency < 10 {
		errData = &errortypes.ErrorData{
			Error:   "check_frequency_invalid",
			Message: "Check frequency cannot be less then 10 seconds",
		}
		return
	}

	if a.Frequency > 3600 {
		errData = &errortypes.ErrorData{
			Error:   "check_frequency_invalid",
			Message: "Check frequency too large",
		}
		return
	}

	if a.Targets == nil {
		a.Targets = []string{}
	}

	switch a.Type {
	case Http:
		break
	case Ping:
		a.Method = ""
		a.Headers = []*Header{}
		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "check_type_invalid",
			Message: "Check type is invalid",
		}
		return
	}

	switch strings.ToUpper(a.Method) {
	case "":
		a.Method = ""
		break
	case "GET":
		a.Method = "GET"
		break
	case "HEAD":
		a.Method = "HEAD"
		break
	case "POST":
		a.Method = "POST"
		break
	case "PUT":
		a.Method = "PUT"
		break
	case "DELETE":
		a.Method = "DELETE"
		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "check_method_invalid",
			Message: "Check method is invalid",
		}
		return
	}

	if a.Headers == nil {
		a.Headers = []*Header{}
	}

	if a.StatusCode <= 0 || a.StatusCode > 900 {
		a.StatusCode = 200
	}

	for _, header := range a.Headers {
		header.Key = utils.FilterStr(header.Key, 256)
		header.Value = utils.FilterStr(header.Value, 2048)
	}

	if a.Timeout < 1 {
		a.Timeout = 5
	} else if a.Timeout > 30 {
		a.Timeout = 30
	}

	return
}

func (a *Check) Commit(db *database.Database) (err error) {
	coll := db.Checks()

	err = coll.Commit(a.Id, a)
	if err != nil {
		return
	}

	return
}

func (a *Check) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Checks()

	err = coll.CommitFields(a.Id, a, fields)
	if err != nil {
		return
	}

	return
}

func (a *Check) Insert(db *database.Database) (err error) {
	coll := db.Checks()

	_, err = coll.InsertOne(db, a)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
