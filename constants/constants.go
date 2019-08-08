package constants

import (
	"go/build"
	"path"
	"time"
)

const (
	Version         = "1.0.1353.86"
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
	Interrupt  = false
	StaticRoot = []string{
		"www/dist",
		"/usr/share/pritunl-zero/www",
		path.Join(
			build.Default.GOPATH,
			"src/github.com/pritunl/pritunl-zero/www/dist",
		),
	}
	StaticTestingRoot = []string{
		"www",
		"/usr/share/pritunl-zero/www",
		path.Join(
			build.Default.GOPATH,
			"src/github.com/pritunl/pritunl-zero/www",
		),
	}
)
