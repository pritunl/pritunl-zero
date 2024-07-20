package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/pritunl/pritunl-zero/cmd"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/logger"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/task"
)

const help = `
Usage: pritunl-zero COMMAND

Commands:
  version           Show version
  mongo             Set MongoDB URI
  set               Set a setting
  unset             Unset a setting
  start             Start node
  clear-logs        Clear logs
  default-password  Get default administrator password
  reset-password    Reset administrator password
  disable-policies  Disable all policies
  export-ssh        Export SSH authorities for emergency client
`

func Init() {
	logger.Init()
	requires.Init()
	task.Init()
}

func main() {
	defer time.Sleep(1 * time.Second)

	flag.Usage = func() {
		fmt.Printf(help)
	}

	flag.Parse()

	switch flag.Arg(0) {
	case "start":
		for _, arg := range flag.Args() {
			switch arg {
			case "--debug":
				constants.Production = false
				break
			case "--debug-web":
				constants.DebugWeb = true
				break
			case "--fast-exit":
				constants.FastExit = true
				break
			}
		}

		Init()
		err := cmd.Node(false)
		if err != nil {
			panic(err)
		}
		return
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
	case "reset-id":
		logger.Init()
		err := cmd.ResetId()
		if err != nil {
			panic(err)
		}
		return
	case "default-password":
		Init()
		err := cmd.DefaultPassword()
		if err != nil {
			panic(err)
		}
		return
	case "reset-password":
		Init()
		err := cmd.ResetPassword()
		if err != nil {
			panic(err)
		}
		return
	case "disable-policies":
		Init()
		err := cmd.DisablePolicies()
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
	case "unset":
		Init()
		err := cmd.SettingsUnset()
		if err != nil {
			panic(err)
		}
		return
	case "export-ssh":
		Init()
		err := cmd.ExportSsh()
		if err != nil {
			panic(err)
		}
		return
	case "clear-logs":
		Init()
		err := cmd.ClearLogs()
		if err != nil {
			panic(err)
		}
		return
	}

	fmt.Printf(help)
}
