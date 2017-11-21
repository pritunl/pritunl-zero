package authority

import (
	"gopkg.in/mgo.v2/bson"
)

type Authority struct {
	Id    bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name  string        `bson:"name" json:"name"`
	Roles []string      `bson:"roles" json:"roles"`
}
