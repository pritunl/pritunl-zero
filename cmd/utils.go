package cmd

import (
	"github.com/pritunl/pritunl-zero/logger"
	"github.com/pritunl/pritunl-zero/requires"
)

func Init() {
	logger.Init()
	requires.Init()
}
