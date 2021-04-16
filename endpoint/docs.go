package endpoint

import (
	"time"

	"github.com/pritunl/mongo-go-driver/bson/primitive"
)

type System struct {
	Id        primitive.ObjectID `json:"i"`
	Timestamp time.Time          `json:"t"`
	Type      string             `json:"x"`

	CpuUsage  float64 `json:"cu"`
	MemTotal  int     `json:"mt"`
	MemUsage  float64 `json:"mu"`
	SwapTotal int     `json:"st"`
	SwapUsage float64 `json:"su"`
}
