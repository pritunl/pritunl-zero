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
  node     Start node
`

func Init() {
	logger.Init()
	requires.Init()
}

func main() {
	defer time.Sleep(1 * time.Second)

	flag.Parse()

	switch flag.Arg(0) {
	case "node":
		switch flag.Arg(1) {
		case "management":
			Init()
			cmd.ManagementNode()
			return
		}
	case "version":
		fmt.Printf("pritunl-zero v%s\n", constants.Version)
		return
	case "mongo":
		logger.Init()
		err := cmd.Mongo()
		if err != nil {
			panic(err)
		}
		return
	case "set":
		Init()
		err := cmd.SettingsSet()
		if err != nil {
			panic(err)
		}
		return
	}

	fmt.Println(help)
}
