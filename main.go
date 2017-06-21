package main

import (
	"flag"
	"fmt"
	"github.com/pritunl/pritunl-zero/cmd"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/logger"
	"github.com/pritunl/pritunl-zero/requires"
	"time"
)

const help = `
Usage: pritunl-zero COMMAND

Commands:
  version  Show version
  mongo    Set MongoDB URI
  set      Set a setting
`

func Init() {
	logger.Init()
	requires.Init()
}

func main() {
	flag.Parse()

	switch flag.Arg(0) {
	case "version":
		fmt.Printf("pritunl-zero v%s\n", constants.Version)
		break
	case "mongo":
		logger.Init()

		err := cmd.Mongo()
		if err != nil {
			panic(err)
		}
		break
	case "set":
		Init()
		err := cmd.SettingsSet()
		if err != nil {
			panic(err)
		}
		break
	default:
		fmt.Println(help)
	}

	time.Sleep(1 * time.Second)
}
