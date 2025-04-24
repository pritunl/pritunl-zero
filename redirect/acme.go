package main

import (
	"sync"
)

type Challenge struct {
	Token    string `json:"token"`
	Response string `json:"response"`
}

var (
	challenges     = map[string]*Challenge{}
	challengesLock = sync.Mutex{}
)

func AddChallenge(chal *Challenge) {
	challengesLock.Lock()
	challenges[chal.Token] = chal
	challengesLock.Unlock()
}

func GetChallenge(token string) (chal *Challenge) {
	challengesLock.Lock()
	chal = challenges[token]
	challengesLock.Unlock()
	return
}
