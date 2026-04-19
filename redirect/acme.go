package main

import (
	"sync"
	"time"
)

type Challenge struct {
	Token    string `json:"token"`
	Response string `json:"response"`
}

var (
	challenges     = map[string]*Challenge{}
	challengesLock = sync.Mutex{}
	clearTimer     *time.Timer
	timerLock      = sync.Mutex{}
)

func AddChallenge(chal *Challenge) {
	timerLock.Lock()
	if clearTimer != nil {
		clearTimer.Stop()
		clearTimer = nil
	}

	clearTimer = time.AfterFunc(60*time.Second, func() {
		challengesLock.Lock()
		challenges = map[string]*Challenge{}
		challengesLock.Unlock()

		timerLock.Lock()
		clearTimer = nil
		timerLock.Unlock()
	})
	timerLock.Unlock()

	challengesLock.Lock()
	challenges[chal.Token] = chal
	challengesLock.Unlock()
}

func GetChallenge(token string) (chal *Challenge) {
	// Retry with backoff for race condition handling
	delays := []time.Duration{
		0 * time.Millisecond,
		150 * time.Millisecond,
		300 * time.Millisecond,
		800 * time.Millisecond,
		1600 * time.Millisecond,
	}

	for i, delay := range delays {
		if i > 0 {
			time.Sleep(delay)
		}
		challengesLock.Lock()
		chal = challenges[token]
		challengesLock.Unlock()
		if chal != nil {
			return
		}
	}
	return
}
