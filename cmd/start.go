package cmd

import (
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
	"github.com/pritunl/pritunl-zero/task"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	startDebug    = false
	startDebugWeb = false
	startFastExit = false
)

func init() {
	AddCmd.PersistentFlags().BoolVar(
		&startDebug,
		"debug",
		false,
		"Debug mode",
	)
	AddCmd.PersistentFlags().BoolVar(
		&startDebugWeb,
		"debug-web",
		false,
		"Web server debug mode",
	)
	AddCmd.PersistentFlags().BoolVar(
		&startFastExit,
		"fast-exit",
		false,
		"Exit without delay",
	)
	RootCmd.AddCommand(AddCmd)
}

var AddCmd = &cobra.Command{
	Use:   "start",
	Short: "Start server",
	Run: func(cmd *cobra.Command, args []string) {
		constants.Production = !startDebug
		constants.DebugWeb = startDebugWeb
		constants.FastExit = startFastExit

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

		logrus.WithFields(logrus.Fields{
			"production": constants.Production,
		}).Info("router: Starting node")

		routr := &router.Router{}
		routr.Init()

		err = task.Init()
		if err != nil {
			return
		}

		go func() {
			err = routr.Run()
			if err != nil && !constants.Interrupt {
				panic(err)
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
