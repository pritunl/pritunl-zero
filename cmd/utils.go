package cmd

import (
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/logger"
	"github.com/pritunl/pritunl-zero/requires"
)

func Init() {
	logger.Init()
	requires.Init()
}

func InitMinimal() {
	logger.Init()

	err := config.Load()
	if err != nil {
		panic(err)
	}

	err = database.Connect()
	if err != nil {
		panic(err)
	}
}
