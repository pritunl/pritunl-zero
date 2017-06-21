package constants

import (
	"time"
)

const (
	Production        = false
	BuildTest         = false
	StaticRoot        = "www/dist"
	StaticTestingRoot = "www"
	StaticCache       = false
	RetryDelay        = 3 * time.Second
)
