package auth

import (
	"net/http"
	"time"
)

var (
	client = &http.Client{
		Timeout: 20 * time.Second,
	}
)

type authData struct {
	Url string `json:"url"`
}

type Token struct {
	Id        string    `bson:"_id"`
	Type      string    `bson:"type"`
	Secret    string    `bson:"secret"`
	Timestamp time.Time `bson:"timestamp"`
}
