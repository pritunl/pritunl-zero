package search

import (
	"time"
)

func dateSuffix() string {
	return time.Now().Format("-2006-01-02")
}
