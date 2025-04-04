package cmd

import (
	"github.com/pritunl/pritunl-zero/logger"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/task"
)

func Init() {
	logger.Init()
	requires.Init()
	task.Init()
}
