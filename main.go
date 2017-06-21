package main

import (
	"flag"
	"fmt"
	"github.com/pritunl/pritunl-zero/cmd"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/logger"
	"github.com/pritunl/pritunl-zero/requires"
)

const help = `
Usage: pritunl-zero COMMAND

Commands:
  version  Show version
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
	case "set":
		Init()
		cmd.SettingsSet()
		break
	default:
		fmt.Println(help)
	}
}
