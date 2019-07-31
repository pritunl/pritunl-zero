package rokey

import (
	"time"

	"github.com/pritunl/mongo-go-driver/bson/primitive"
)

type Rokey struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Type      string             `bson:"type" json:"type"`
	Timeblock time.Time          `bson:"timeblock" json:"timeblock"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Secret    string             `bson:"secret" json:"-"`
}
