package check

import (
	"strings"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/mongo-go-driver/bson"
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
	States     []*State           `bson:"states" json:"states"`
}

type State struct {
	Endpoint  primitive.ObjectID `bson:"e" json:"e"`
	Timestamp time.Time          `bson:"t" json:"t"`
	Targets   []string           `bson:"x" json:"x"`
	Latency   []int              `bson:"l" json:"l"`
	Errors    []string           `bson:"r" json:"r"`
}

type Header struct {
	Key   string `bson:"key" json:"key"`
	Value string `bson:"value" json:"value"`
}

func (c *Check) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if c.Id.IsZero() {
		c.Id, err = utils.RandObjectId()
		if err != nil {
			return
		}
	}

	if c.Roles == nil {
		c.Roles = []string{}
	}

	if c.Frequency == 0 {
		c.Frequency = 30
	}

	if c.Frequency < 10 {
		errData = &errortypes.ErrorData{
			Error:   "check_frequency_invalid",
			Message: "Check frequency cannot be less then 10 seconds",
		}
		return
	}

	if c.Frequency > 3600 {
		errData = &errortypes.ErrorData{
			Error:   "check_frequency_invalid",
			Message: "Check frequency too large",
		}
		return
	}

	if c.Targets == nil {
		c.Targets = []string{}
	}

	switch c.Type {
	case Http:
		break
	case Ping:
		c.Method = ""
		c.Headers = []*Header{}
		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "check_type_invalid",
			Message: "Check type is invalid",
		}
		return
	}

	if c.Type == Http {
		switch strings.ToUpper(c.Method) {
		case "":
			c.Method = "GET"
			break
		case "GET":
			c.Method = "GET"
			break
		case "HEAD":
			c.Method = "HEAD"
			break
		default:
			errData = &errortypes.ErrorData{
				Error:   "check_method_invalid",
				Message: "Check method is invalid",
			}
			return
		}
	}

	if c.Headers == nil {
		c.Headers = []*Header{}
	}

	if c.StatusCode <= 0 || c.StatusCode > 900 {
		c.StatusCode = 200
	}

	for _, header := range c.Headers {
		header.Key = utils.FilterStr(header.Key, 256)
		header.Value = utils.FilterStr(header.Value, 2048)
	}

	if c.Timeout < 1 {
		c.Timeout = 5
	} else if c.Timeout > 30 {
		c.Timeout = 30
	}

	if c.States == nil {
		c.States = []*State{}
	}

	return
}

func (c *Check) UpdateState(db *database.Database, state *State) (
	updated bool, err error) {

	coll := db.Checks()

	insert := true
	updated = true

	for _, stat := range c.States {
		if stat.Endpoint == state.Endpoint {
			insert = false
			if stat.Timestamp == state.Timestamp {
				updated = false
			}
		}
	}

	if insert {
		_, err = coll.UpdateOne(db, &bson.M{
			"_id": c.Id,
			"states.e": &bson.M{
				"$ne": state.Endpoint,
			},
		}, &bson.M{
			"$push": &bson.M{
				"states": state,
			},
		})
	} else {
		_, err = coll.UpdateOne(db, &bson.M{
			"_id":      c.Id,
			"states.e": state.Endpoint,
		}, &bson.M{
			"$set": &bson.M{
				"states.$": state,
			},
		})
	}
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (c *Check) Commit(db *database.Database) (err error) {
	coll := db.Checks()

	err = coll.Commit(c.Id, c)
	if err != nil {
		return
	}

	return
}

func (c *Check) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Checks()

	err = coll.CommitFields(c.Id, c, fields)
	if err != nil {
		return
	}

	return
}

func (c *Check) Insert(db *database.Database) (err error) {
	coll := db.Checks()

	_, err = coll.InsertOne(db, c)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
