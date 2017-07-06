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
