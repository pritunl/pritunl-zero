package secondary

import (
	"math/rand"
	"time"

	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
)

func New(db *database.Database, userId bson.ObjectID, typ string,
	proivderId bson.ObjectID) (secd *Secondary, err error) {

	token, err := utils.RandStr(64)
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

func NewChallenge(db *database.Database, userId bson.ObjectID,
	typ string, chalId string, proivderId bson.ObjectID) (
	secd *Secondary, err error) {

	token, err := utils.RandStr(64)
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

	err = coll.FindOne(db, &bson.M{
		"_id":  token,
		"type": typ,
		"timestamp": &bson.M{
			"$gte": timestamp,
		},
	}).Decode(secd)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, token string) (err error) {
	coll := db.SecondaryTokens()

	_, err = coll.DeleteMany(db, &bson.M{
		"_id": token,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
