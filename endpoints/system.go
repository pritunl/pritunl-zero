package endpoints

import (
	"time"

	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
)

type System struct {
	Id        []byte             `bson:"_id" json:"id"`
	ClientId  primitive.ObjectID `bson:"-" json:"i"`
	Endpoint  primitive.ObjectID `bson:"e" json:"e"`
	Timestamp time.Time          `bson:"t" json:"t"`
	Type      string             `bson:"x" json:"x"`

	CpuUsage  float64 `bson:"cu" json:"cu"`
	MemTotal  int     `bson:"mt" json:"mt"`
	MemUsage  float64 `bson:"mu" json:"mu"`
	SwapTotal int     `bson:"st" json:"st"`
	SwapUsage float64 `bson:"su" json:"su"`
}

func (d *System) GetCollection(db *database.Database) *database.Collection {
	return db.EndpointsSystem()
}

func (d *System) SetEndpoint(id primitive.ObjectID) {
	d.Id = GenerateId(id, d.ClientId)
	d.Endpoint = id
}
