package database

import (
	"github.com/pritunl/mongo-go-driver/v2/bson"
)

type Named struct {
	Id   bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string        `bson:"name" json:"name"`
}
