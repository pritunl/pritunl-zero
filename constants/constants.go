package constants

import (
	"time"
)

const (
	Version           = "1.0.0"
	Production        = false
	BuildTest         = false
	StaticRoot        = "www/dist"
	StaticTestingRoot = "www"
	StaticCache       = false
	RetryDelay        = 3 * time.Second
)
