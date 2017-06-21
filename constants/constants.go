package constants

import (
	"time"
)

const (
	Version           = "1.0.0"
	ConfPath          = "/etc/pritunl_zero.json"
	Production        = false
	BuildTest         = false
	StaticRoot        = "www/dist"
	StaticTestingRoot = "www"
	StaticCache       = false
	RetryDelay        = 3 * time.Second
)
