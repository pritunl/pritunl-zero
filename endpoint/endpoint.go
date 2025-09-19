package endpoint

import (
	"bytes"
	"context"
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
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/alert"
	"github.com/pritunl/pritunl-zero/alertevent"
	"github.com/pritunl/pritunl-zero/check"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/endpoints"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/nonce"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/nacl/box"
)

type Endpoint struct {
	Id            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string        `bson:"name" json:"name"`
	Username      string        `bson:"username" json:"username"`
	Roles         []string      `bson:"roles" json:"roles"`
	ClientKey     *ClientKey    `bson:"client_key" json:"client_key"`
	ServerKey     *ServerKey    `bson:"server_key" json:"-"`
	Info          *EndpointInfo `bson:"-" json:"info"`
	HasClientKey  bool          `bson:"-" json:"has_client_key"`
	Data          *Data         `bson:"data" json:"data"`
	keyLoaded     bool          `bson:"-" json:"-"`
	clientPubKey  [32]byte      `bson:"-" json:"-"`
	serverPrivKey [32]byte      `bson:"-" json:"-"`
}

type EndpointInfo struct {
	Domain string            `json:"domain"`
	Alerts map[string]string `json:"alerts"`
	Checks map[string]string `json:"checks"`
}

type EndpointConf struct {
	Checks []*check.Check `json:"checks"`
}

type Data struct {
	Version        string     `bson:"version" json:"version"`
	Hostname       string     `bson:"hostname" json:"hostname"`
	Uptime         uint64     `bson:"uptime" json:"uptime"`
	Platform       string     `bson:"platform" json:"platform"`
	PackageUpdates int        `bson:"package_updates" json:"package_updates"`
	Virtualization string     `bson:"virtualization" json:"virtualization"`
	CpuCores       int        `bson:"cpu_cores" json:"cpu_cores"`
	MemTotal       int        `bson:"mem_total" json:"mem_total"`
	SwapTotal      int        `bson:"swap_total" json:"swap_total"`
	HugeTotal      int        `bson:"huge_total" json:"huge_total"`
	Mdadm          []*MdState `bson:"md_stat" json:"md_stat"`
}

type MdState struct {
	Name   string `bson:"name" json:"name"`
	State  string `bson:"state" json:"state"`
	Level  string `bson:"level" json:"level"`
	Failed int    `bson:"failed" json:"failed"`
	Spare  int    `bson:"spare" json:"spare"`
	Total  int    `bson:"total" json:"total"`
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

func (e *Endpoint) GetKeys() (clientPubKey, serverPrivKey *[32]byte,
	err error) {

	if !e.keyLoaded {
		clientPubKeySl, er := base64.StdEncoding.DecodeString(
			e.ClientKey.PublicKey)
		if er != nil {
			err = &errortypes.ParseError{
				errors.Wrap(er,
					"endpoint: Failed to decode client private key"),
			}
			return
		}
		copy(e.clientPubKey[:], clientPubKeySl)

		serverPubKeySl, er := base64.StdEncoding.DecodeString(
			e.ServerKey.PrivateKey)
		if er != nil {
			err = &errortypes.ParseError{
				errors.Wrap(er,
					"endpoint: Failed to decode server public key"),
			}
			return
		}
		copy(e.serverPrivKey[:], serverPubKeySl)

		e.keyLoaded = true
	}

	clientPubKey = &e.clientPubKey
	serverPrivKey = &e.serverPrivKey

	return
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
		PublicKey:  base64.StdEncoding.EncodeToString(pubKey[:]),
		PrivateKey: base64.StdEncoding.EncodeToString(privKey[:]),
	}

	return
}

func (e *Endpoint) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	e.Name = utils.FilterName(e.Name)

	if e.Id.IsZero() {
		e.Id, err = utils.RandObjectId()
		if err != nil {
			return
		}
	}

	if e.Roles == nil {
		e.Roles = []string{}
	}

	if e.ClientKey == nil || e.ServerKey == nil {
		err = e.GenerateKey()
		if err != nil {
			return
		}
	}

	if e.Data == nil {
		e.Data = &Data{}
	}

	e.Format()

	return
}

func (e *Endpoint) Format() {
	sort.Strings(e.Roles)
}

func (e *Endpoint) Json(alertsMap map[string][]*alert.Alert,
	checksMap map[string][]*check.Check) {

	if e.ClientKey != nil && e.ClientKey.PublicKey != "" {
		e.ClientKey = nil
		e.HasClientKey = true
	}

	if alertsMap != nil {
		alerts := map[string]string{}
		checks := map[string]string{}

		for _, role := range e.Roles {
			roleAlrts := alertsMap[role]
			if roleAlrts != nil {
				for _, alrt := range roleAlrts {
					alerts[alrt.Id.Hex()] = alrt.Name
				}
			}

			roleChcks := checksMap[role]
			if roleChcks != nil {
				for _, chck := range roleChcks {
					checks[chck.Id.Hex()] = chck.Name
				}
			}
		}

		e.Info = &EndpointInfo{
			Domain: node.Self.EndpointDomain,
			Alerts: alerts,
			Checks: checks,
		}
	}
}

func (e *Endpoint) GetConf(db *database.Database) (
	encData []byte, err error) {

	clientPubKey, serverPrivKey, err := e.GetKeys()
	if err != nil {
		return
	}

	if clientPubKey == nil || serverPrivKey == nil {
		return
	}

	checks, err := check.GetRoles(db, e.Roles)
	if err != nil {
		return
	}

	conf := &EndpointConf{
		Checks: checks,
	}

	confData, err := json.Marshal(conf)
	if err != nil {
		err = &errortypes.ParseError{
			errors.New("endpoint: Failed to marshal endpoint conf"),
		}
		return
	}

	encData, err = utils.RandBytes(24)
	if err != nil {
		return
	}
	var nonceAr [24]byte
	copy(nonceAr[:], encData)

	encData = box.Seal(encData, confData, &nonceAr,
		clientPubKey, serverPrivKey)

	return
}

func (e *Endpoint) ValidateSignature(db *database.Database,
	timestampStr, nonc, sig, method string) (errData *errortypes.ErrorData,
	err error) {

	if e.ClientKey == nil || e.ServerKey == nil {
		errData = &errortypes.ErrorData{
			Error:   "not_initialized",
			Message: "Endpoint key not initialized",
		}
		return
	}

	if len(nonc) < 16 || len(nonc) > 128 {
		err = &errortypes.AuthenticationError{
			errors.New("endpoint: Invalid authentication nonce"),
		}
		return
	}

	timestampInt, _ := strconv.ParseInt(timestampStr, 10, 64)
	if timestampInt == 0 {
		err = &errortypes.AuthenticationError{
			errors.New("endpoint: Invalid authentication timestamp"),
		}
		return
	}

	timestamp := time.Unix(timestampInt, 0)
	if utils.SinceAbs(timestamp) > time.Duration(
		settings.Auth.WindowLong)*time.Second {

		err = &errortypes.AuthenticationError{
			errors.New("endpoint: Authentication timestamp outside window"),
		}
		return
	}

	authString := strings.Join([]string{
		timestampStr,
		nonc,
		method,
	}, "&")

	err = nonce.Validate(db, nonc)
	if err != nil {
		return
	}

	if e.ClientKey.Secret == "" {
		err = &errortypes.ReadError{
			errors.Wrap(err, "session: Empty secret"),
		}
		return
	}

	hashFunc := hmac.New(sha512.New, []byte(e.ClientKey.Secret))
	hashFunc.Write([]byte(authString))
	rawSignature := hashFunc.Sum(nil)
	testSig := base64.URLEncoding.EncodeToString(rawSignature)

	if subtle.ConstantTimeCompare([]byte(sig), []byte(testSig)) != 1 {
		errData = &errortypes.ErrorData{
			Error:   "authentication_error",
			Message: "Register signature does not match",
		}
		return
	}

	return
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

	if len(reqData.Nonce) < 16 || len(reqData.Nonce) > 128 {
		err = &errortypes.AuthenticationError{
			errors.New("endpoint: Invalid authentication nonce"),
		}
		return
	}

	if len(reqData.PublicKey) < 16 || len(reqData.PublicKey) > 512 {
		err = &errortypes.AuthenticationError{
			errors.New("endpoint: Invalid public key"),
		}
		return
	}

	timestamp := time.Unix(reqData.Timestamp, 0)
	if utils.SinceAbs(timestamp) > time.Duration(
		settings.Auth.WindowLong)*time.Second {

		err = &errortypes.AuthenticationError{
			errors.New("endpoint: Authentication timestamp outside window"),
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

	err = nonce.Validate(db, reqData.Nonce)
	if err != nil {
		return
	}

	if e.ClientKey.Secret == "" {
		err = &errortypes.ReadError{
			errors.Wrap(err, "session: Empty secret"),
		}
		return
	}

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

	if e.ClientKey.Secret == "" {
		err = &errortypes.ReadError{
			errors.Wrap(err, "session: Empty secret"),
		}
		return
	}

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

func (e *Endpoint) GetAlerts(db *database.Database) (
	alerts []*alert.Alert, err error) {

	alerts, err = alert.GetRoles(db, e.Roles)
	if err != nil {
		return
	}

	return
}

func (e *Endpoint) InsertDoc(db *database.Database,
	msgData []byte) (err error) {

	if len(msgData) < 32 {
		logrus.WithFields(logrus.Fields{
			"data_len": len(msgData),
		}).Error("endpoint: Data too short")
		return
	}

	clientPubKey, serverPrivKey, err := e.GetKeys()
	if err != nil {
		return
	}

	var nonceAr [24]byte
	copy(nonceAr[:], msgData[:24])

	docData, valid := box.Open([]byte{}, msgData[24:],
		&nonceAr, clientPubKey, serverPrivKey)
	if !valid {
		logrus.WithFields(logrus.Fields{
			"data_len": len(docData),
		}).Error("endpoint: Failed to decrypt doc")
		return
	}

	sepIndex := bytes.Index(docData, []byte(":"))
	if sepIndex == -1 {
		logrus.WithFields(logrus.Fields{
			"data_len": len(docData),
		}).Error("endpoint: Failed to parse doc type")
		return
	}

	docType := docData[:sepIndex]

	doc, err := UnmarshalDoc(string(docType), docData[sepIndex+1:])
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("endpoint: Failed to unmarshal doc")
		err = nil
		return
	}

	timestamp := doc.Format(e.Id)

	staticData := doc.StaticData()
	if staticData != nil && (!constants.Production ||
		timestamp == timestamp.Truncate(5*time.Minute)) {

		coll := db.Endpoints()
		err = coll.UpdateId(e.Id, &bson.M{
			"$set": staticData,
		})
	}

	handled, checkAlerts, err := doc.Handle(db)
	if err != nil {
		return
	}

	if !handled {
		checkAlerts = true
		coll := doc.GetCollection(db)

		_, err = coll.InsertOne(db, doc)
		if err != nil {
			err = database.ParseError(err)
			if _, ok := err.(*database.DuplicateKeyError); ok {
				err = nil
				checkAlerts = false
			} else {
				return
			}
		}
	}

	if checkAlerts {
		alerts, er := e.GetAlerts(db)
		if er != nil {
			err = er
			return
		}

		actAlrts := doc.CheckAlerts(alerts)
		if actAlrts != nil && len(actAlrts) > 0 {
			for _, alrt := range actAlrts {
				go alertevent.New(e.Roles, e.Id, alrt.Name, e.Name, alrt.Resource,
					alrt.Message, alrt.Level, alrt.Frequency)
			}
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

	_, err = coll.InsertOne(db, e)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (e *Endpoint) GetChart(c context.Context, db *database.Database,
	resource string, start, end time.Time, interval time.Duration) (
	data endpoints.ChartData, err error) {

	data, err = endpoints.GetChart(c, db, e.Id, resource,
		start, end, interval)
	if err != nil {
		return
	}

	return
}

func (e *Endpoint) GetLog(c context.Context, db *database.Database,
	resource string) (data endpoints.LogData, err error) {

	data, err = endpoints.GetLog(c, db, e.Id, resource)
	if err != nil {
		return
	}

	return
}

func UnmarshalDoc(docType string, docData []byte) (
	doc endpoints.Doc, err error) {

	doc = endpoints.GetObj(docType)

	err = json.Unmarshal(docData, doc)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "endpoints: Failed to parse doc"),
		}
		return
	}

	return
}
