package csrf

import (
	"time"

	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/utils"
)

type CsrfToken struct {
	Id        string    `bson:"_id"`
	Session   string    `bson:"session"`
	Timestamp time.Time `bson:"timestamp"`
}

func NewToken(db *database.Database, sessionId string) (
	token string, err error) {

	coll := db.CsrfTokens()

	tkn, err := utils.RandStr(48)
	if err != nil {
		return
	}

	doc := &CsrfToken{
		Id:        tkn,
		Session:   sessionId,
		Timestamp: time.Now(),
	}

	_, err = coll.InsertOne(db, doc)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	token = tkn
	return
}

func ValidateToken(db *database.Database, sessionId, token string) (
	valid bool, err error) {

	coll := db.CsrfTokens()

	doc := &CsrfToken{}

	err = coll.FindOneId(token, doc)
	if err != nil {
		return
	}

	if doc.Session == sessionId {
		valid = true
		return
	}

	return
}
