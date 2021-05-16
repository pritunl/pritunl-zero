package endpoints

import (
	"fmt"
	"time"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
)

type System struct {
	Id        primitive.Binary   `bson:"_id" json:"id"`
	Endpoint  primitive.ObjectID `bson:"e" json:"e"`
	Timestamp time.Time          `bson:"t" json:"t"`

	CpuUsage  float64 `bson:"cu" json:"cu"`
	MemTotal  int     `bson:"mt" json:"mt"`
	MemUsage  float64 `bson:"mu" json:"mu"`
	SwapTotal int     `bson:"st" json:"st"`
	SwapUsage float64 `bson:"su" json:"su"`
}

func (d *System) GetCollection(db *database.Database) *database.Collection {
	return db.EndpointsSystem()
}

func (d *System) Format(id primitive.ObjectID) {
	d.Endpoint = id
	d.Timestamp = d.Timestamp.UTC().Truncate(1 * time.Minute)
	d.Id = GenerateId(id, d.Timestamp)
}

func (d *System) Print() {
	fmt.Println("***************************************************")
	fmt.Printf("Id: %x\n", d.Id)
	fmt.Println("Timestamp:", d.Timestamp)
	fmt.Println("Type: system")

	fmt.Println("CpuUsage:", d.CpuUsage)
	fmt.Println("MemTotal:", d.MemTotal)
	fmt.Println("MemUsage:", d.MemUsage)
	fmt.Println("SwapTotal:", d.SwapTotal)
	fmt.Println("SwapUsage:", d.SwapUsage)
	fmt.Println("***************************************************")
}
