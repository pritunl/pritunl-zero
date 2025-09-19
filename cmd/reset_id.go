package cmd

import (
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ResetIdCmd)
}

var ResetIdCmd = &cobra.Command{
	Use:   "reset-id",
	Short: "Reset node ID",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		err := config.Load()
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		config.Config.NodeId = bson.NewObjectID().Hex()

		err = config.Save()
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		logrus.WithFields(logrus.Fields{
			"node_id": config.Config.NodeId,
		}).Info("cmd: Reset node ID")
	},
}
