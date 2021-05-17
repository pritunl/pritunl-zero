package endpoint

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"sort"
	"time"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/endpoints"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/utils"
	"golang.org/x/crypto/nacl/box"
)

type Endpoint struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User      primitive.ObjectID `bson:"user,omitempty" json:"user"`
	Name      string             `bson:"name" json:"name"`
	Roles     []string           `bson:"roles" json:"roles"`
	ClientKey *ClientKey         `bson:"client_key" json:"client_key"`
	ServerKey *ServerKey         `bson:"server_key" json:"server_key"`
}

type ClientKey struct {
	PublicKey string `bson:"public_key" json:"-"`
	Secret    string `bson:"secret" json:"secret"`
}

type ServerKey struct {
	PrivateKey string `bson:"private_key" json:"-"`
	PublicKey  string `bson:"public_key" json:"-"`
}

func (e *Endpoint) GenerateKey() (err error) {
	pubKey, privKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "endpoint: Failed to generate nacl key"),
		}
		return
	}

	secret, err := utils.RandStr(64)
	if err != nil {
		return
	}

	e.ClientKey = &ClientKey{
		Secret: secret,
	}

	e.ServerKey = &ServerKey{
		PublicKey:  base64.RawStdEncoding.EncodeToString(pubKey[:]),
		PrivateKey: base64.RawStdEncoding.EncodeToString(privKey[:]),
	}

	return
}

func (e *Endpoint) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if e.Roles == nil {
		e.Roles = []string{}
	}

	if e.ClientKey == nil || e.ServerKey == nil {
		err = e.GenerateKey()
		if err != nil {
			return
		}
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

	doc.Format(e.Id)

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
	start, end time.Time, interval time.Duration) (
	data interface{}, err error) {

	data, err = endpoints.GetChart(db, e.Id, resource, start, end, interval)
	if err != nil {
		return
	}

	return
}

func UnmarshalDoc(docType string, docData string) (
	doc endpoints.Doc, err error) {

	doc = endpoints.GetObj(docType)

	err = json.Unmarshal([]byte(docData), doc)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "endpoints: Failed to parse doc"),
		}
		return
	}

	return
}
