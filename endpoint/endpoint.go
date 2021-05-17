package endpoint

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
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

type RegisterData struct {
	Timestamp int64  `json:"timestamp"`
	Nonce     string `json:"nonce"`
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
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

func (e *Endpoint) Register(db *database.Database, reqData *RegisterData) (
	resData *RegisterData, errData *errortypes.ErrorData, err error) {

	if e.ClientKey == nil || e.ServerKey == nil {
		errData = &errortypes.ErrorData{
			Error:   "not_initialized",
			Message: "Endpoint key not initialized",
		}
		return
	}

	if e.ClientKey.PublicKey != "" {
		errData = &errortypes.ErrorData{
			Error:   "already_registered",
			Message: "Endpoint is already registered",
		}
		return
	}

	authString := strings.Join([]string{
		strconv.FormatInt(reqData.Timestamp, 10),
		reqData.Nonce,
		reqData.PublicKey,
	}, "&")

	hashFunc := hmac.New(sha512.New, []byte(e.ClientKey.Secret))
	hashFunc.Write([]byte(authString))
	rawSignature := hashFunc.Sum(nil)
	testSig := base64.StdEncoding.EncodeToString(rawSignature)

	if subtle.ConstantTimeCompare([]byte(
		reqData.Signature), []byte(testSig)) != 1 {

		errData = &errortypes.ErrorData{
			Error:   "authentication_error",
			Message: "Register signature does not match",
		}
		return
	}

	clientPubKey, err := base64.StdEncoding.DecodeString(reqData.PublicKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "endpoint: Failed to parse register public key"),
		}
		return
	}

	e.ClientKey.PublicKey = base64.StdEncoding.EncodeToString(clientPubKey)

	resData = &RegisterData{
		Timestamp: time.Now().Unix(),
		Nonce:     reqData.Nonce,
		PublicKey: e.ServerKey.PublicKey,
	}

	authString = strings.Join([]string{
		strconv.FormatInt(resData.Timestamp, 10),
		resData.Nonce,
		resData.PublicKey,
	}, "&")

	hashFunc = hmac.New(sha512.New, []byte(e.ClientKey.Secret))
	hashFunc.Write([]byte(authString))
	rawSignature = hashFunc.Sum(nil)
	resData.Signature = base64.StdEncoding.EncodeToString(rawSignature)

	fields := set.NewSet(
		"client_key",
	)

	errData, err = e.Validate(db)
	if err != nil {
		return
	}
	if errData != nil {
		return
	}

	err = e.CommitFields(db, fields)
	if err != nil {
		return
	}

	return
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
