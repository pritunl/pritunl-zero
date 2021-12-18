package search

import (
	"time"
)

var (
	lastLog time.Time
)

func dateSuffix() string {
	return time.Now().Format("-2006-01-02")
}

func logLimit() bool {
	if time.Since(lastLog) > 30*time.Second {
		lastLog = time.Now()
		return true
	} else {
		return false
	}
}
