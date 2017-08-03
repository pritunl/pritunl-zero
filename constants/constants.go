package constants

import (
	"path/filepath"
	"time"
)

const (
	Version           = "1.0.618.80"
	DatabaseVersion   = 1
	ConfPath          = "/etc/pritunl-zero.json"
	LogPath           = "/var/log/pritunl-zero.log"
	TempPath          = "/tmp/pritunl-zero"
	Production        = false
	BuildTest         = false
	StaticRoot        = "/usr/share/pritunl-zero/www"
	StaticTestingRoot = "www"
	StaticCache       = false
	RetryDelay        = 3 * time.Second
)

var (
	Interrupt = false
	CertPath  = filepath.Join(TempPath, "server.crt")
	KeyPath   = filepath.Join(TempPath, "server.key")
)
