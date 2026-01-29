package constants

import (
	"time"
)

const (
	Version         = "1.0.3719.54"
	DatabaseVersion = 1
	ConfPath        = "/etc/pritunl-zero.json"
	LogPath         = "/var/log/pritunl-zero.log"
	LogPath2        = "/var/log/pritunl-zero.log.1"
	TempPath        = "/tmp/pritunl-zero"
	StaticCache     = true
	RetryDelay      = 3 * time.Second
)

var (
	Production = true
	DebugWeb   = false
	FastExit   = false
	Interrupt  = false
	StaticRoot = []string{
		"www/dist",
		"/usr/share/pritunl-zero/www",
	}
	StaticTestingRoot = []string{
		"/home/cloud/git/pritunl-zero/www/dist-dev",
		"/home/cloud/go/src/github.com/pritunl/pritunl-zero/www/dist-dev",
		"/usr/share/pritunl-zero/www",
	}
)
