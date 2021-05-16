package endpoint

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/endpoints"
	"github.com/pritunl/pritunl-zero/errortypes"
)

type Endpoint struct {
	Id    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User  primitive.ObjectID `bson:"user,omitempty" json:"user"`
	Name  string             `bson:"name" json:"name"`
	Roles []string           `bson:"roles" json:"roles"`
}

func (e *Endpoint) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if e.Roles == nil {
		e.Roles = []string{}
	}

	e.Format()

	return
}

func (e *Endpoint) Format() {
	sort.Strings(e.Roles)
}

func (e *Endpoint) InsertDoc(db *database.Database, doc endpoints.Doc) (
	err error) {

	coll := doc.GetCollection(db)

	doc.SetEndpoint(e.Id)

	_, err = coll.InsertOne(db, doc)
	if err != nil {
		err = database.ParseError(err)
		if _, ok := err.(*database.DuplicateKeyError); ok {
			err = nil
		} else {
			return
		}
	}

	return
}

func (e *Endpoint) Commit(db *database.Database) (err error) {
	coll := db.Endpoints()

	err = coll.Commit(e.Id, e)
	if err != nil {
		return
	}

	return
}

func (e *Endpoint) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Endpoints()

	err = coll.CommitFields(e.Id, e, fields)
	if err != nil {
		return
	}

	return
}

func (e *Endpoint) Insert(db *database.Database) (err error) {
	coll := db.Endpoints()

	if !e.Id.IsZero() {
		err = &errortypes.DatabaseError{
			errors.New("endpoint: Endpoint already exists"),
		}
		return
	}

	_, err = coll.InsertOne(db, e)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (e *Endpoint) GetData(db *database.Database, resource string,
	start, end time.Time) (data interface{}, err error) {

	data, err = endpoints.GetChart(db, e.Id, resource, start, end)
	if err != nil {
		return
	}

	return
}

func ProcessDoc(db *database.Database, endpt *Endpoint,
	docType string, docData string) (err error) {

	doc := endpoints.GetObj(docType)

	err = json.Unmarshal([]byte(docData), doc)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "endpoints: Failed to parse doc"),
		}
		return
	}

	doc.SetEndpoint(endpt.Id)

	err = endpt.InsertDoc(db, doc)
	if err != nil {
		return
	}

	return
}
