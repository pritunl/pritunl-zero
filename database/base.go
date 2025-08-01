package database

import (
	"github.com/pritunl/mongo-go-driver/bson/primitive"
)

type Named struct {
	Id   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
}
