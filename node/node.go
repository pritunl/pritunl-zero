package node

import (
	"time"
)

type Node struct {
	Id        string    `bson:"_id" json:"id"`
	Type      string    `bson:"type" json:"type"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Memory    float64   `bson:"memory" json:"memory"`
	Load1     float64   `bson:"load1" json:"load1"`
	Load5     float64   `bson:"load5" json:"load5"`
	Load15    float64   `bson:"load15" json:"load15"`
}
