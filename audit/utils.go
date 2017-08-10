package audit

import (
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/database"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

func Get(db *database.Database, adtId string) (
	adt *Audit, err error) {

	coll := db.Audits()
	adt = &Audit{}

	err = coll.FindOneId(adtId, adt)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database, userId bson.ObjectId) (
	audits []*Audit, err error) {

	coll := db.Audits()
	audits = []*Audit{}

	cursor := coll.Find(bson.M{
		"u": userId,
	}).Sort("-$natural").Iter()

	adt := &Audit{}
	for cursor.Next(adt) {
		audits = append(audits, adt)
		adt = &Audit{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func New(db *database.Database, r *http.Request,
	userId bson.ObjectId, typ string, fields Fields) (
	err error) {

	agnt, err := agent.Parse(db, r)
	if err != nil {
		return
	}

	adt := &Audit{
		User:      userId,
		Timestamp: time.Now(),
		Type:      typ,
		Fields:    fields,
		Agent:     agnt,
	}

	err = adt.Insert(db)
	if err != nil {
		return
	}

	return
}
