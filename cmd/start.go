package cmd

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/router"
	"github.com/pritunl/pritunl-zero/sync"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	startDebug    = false
	startDebugWeb = false
	startFastExit = false
)

func init() {
	AddCmd.PersistentFlags().BoolVarP(
		&startDebug,
		"debug",
		"",
		false,
		"Debug mode",
	)
	AddCmd.PersistentFlags().BoolVarP(
		&startDebugWeb,
		"debug-web",
		"",
		false,
		"Web server debug mode",
	)
	AddCmd.PersistentFlags().BoolVarP(
		&startDebug,
		"fast-exit",
		"",
		false,
		"Exit without delay",
	)
	RootCmd.AddCommand(AddCmd)
}

var AddCmd = &cobra.Command{
	Use:   "start",
	Short: "Start server",
	Run: func(cmd *cobra.Command, args []string) {
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

		objId, err := primitive.ObjectIDFromHex(config.Config.NodeId)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "cmd: Failed to parse ObjectId"),
			}
			cobra.CheckErr(err)
			return
		}

		nde := &node.Node{
			Id: objId,
		}
		err = nde.Init()
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		sync.Init()

		routr := &router.Router{}

		routr.Init()

		go func() {
			err = routr.Run()
			if err != nil && !constants.Interrupt {
				cobra.CheckErr(err)
				return
			}
		}()

		sig := make(chan os.Signal, 2)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig

		constants.Interrupt = true

		logrus.Info("cmd.node: Shutting down")
		go routr.Shutdown()
		if !constants.Production || constants.FastExit {
			time.Sleep(300 * time.Millisecond)
		} else {
			time.Sleep(3 * time.Second)
		}
	},
}
