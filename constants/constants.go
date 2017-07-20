package constants

import (
	"path/filepath"
	"time"
)

const (
	Version                = "1.0.0"
	DatabaseVersion        = 1
	ConfPath               = "/etc/pritunl_zero.json"
	LogPath                = "/var/log/pritunl_zero.log"
	TempPath               = "/tmp/pritunl_zero"
	Production             = false
	BuildTest              = false
	StaticRoot             = "www/dist"
	StaticTestingRoot      = "www"
	ProxyStaticRoot        = "www_proxy/dist"
	ProxyStaticTestingRoot = "www_proxy"
	StaticCache            = false
	RetryDelay             = 3 * time.Second
)

var (
	Interrupt = false
	CertPath  = filepath.Join(TempPath, "server.crt")
	KeyPath   = filepath.Join(TempPath, "server.key")
)
