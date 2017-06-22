package constants

import (
	"time"
)

const (
	Version                = "1.0.0"
	ConfPath               = "/etc/pritunl_zero.json"
	LogPath                = "/var/log/pritunl_zero.log"
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
)
