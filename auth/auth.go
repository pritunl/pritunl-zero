package auth

import (
	"github.com/pritunl/pritunl-zero/database"
	"net/http"
	"time"
)

var (
	client = &http.Client{
		Timeout: 20 * time.Second,
	}
)

type authData struct {
	Url string `json:"url"`
}

type Token struct {
	Id        string    `bson:"_id"`
	Type      string    `bson:"type"`
	Secret    string    `bson:"secret"`
	Timestamp time.Time `bson:"timestamp"`
}

func (t *Token) Remove(db *database.Database) (err error) {
	coll := db.Tokens()

	err = coll.RemoveId(t.Id)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
