package secondary

import (
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"time"
)

func New(db *database.Database, userId bson.ObjectId, typ string,
	proivderId bson.ObjectId) (secd *Secondary, err error) {

	token, err := utils.RandStr(48)
	if err != nil {
		return
	}

	secd = &Secondary{
		Id:         token,
		UserId:     userId,
		Type:       typ,
		ProviderId: proivderId,
		Timestamp:  time.Now(),
	}

	err = secd.Insert(db)
	if err != nil {
		return
	}

	return
}

func NewChallenge(db *database.Database, userId bson.ObjectId,
	typ string, chalId string, proivderId bson.ObjectId) (
	secd *Secondary, err error) {

	token, err := utils.RandStr(48)
	if err != nil {
		return
	}

	secd = &Secondary{
		Id:          token,
		UserId:      userId,
		Type:        typ,
		ChallengeId: chalId,
		ProviderId:  proivderId,
		Timestamp:   time.Now(),
	}

	err = secd.Insert(db)
	if err != nil {
		return
	}

	return
}

func Get(db *database.Database, token string, typ string) (
	secd *Secondary, err error) {

	coll := db.SecondaryTokens()
	secd = &Secondary{}

	timestamp := time.Now().Add(
		-time.Duration(settings.Auth.SecondaryExpire) * time.Second)

	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)

	err = coll.FindOne(&bson.M{
		"_id":  token,
		"type": typ,
		"timestamp": &bson.M{
			"$gte": timestamp,
		},
	}, secd)
	if err != nil {
		return
	}

	return
}

func Remove(db *database.Database, token string) (err error) {
	coll := db.SecondaryTokens()

	_, err = coll.RemoveAll(&bson.M{
		"_id": token,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
