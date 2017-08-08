package audit

import (
	"github.com/pritunl/pritunl-zero/agent"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Event struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	User      bson.ObjectId `bson:"user" json:"user"`
	Timestamp time.Time     `bson:"timestamp" json:"timestamp"`
	Type      string        `bson:"type" json:"type"`
	Agent     *agent.Agent  `bson:"agent" json:"agent"`
	Message   string        `bson:"message" json:"message"`
}
