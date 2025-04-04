package cmd

import (
	"fmt"

	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(MongoCmd)
}

var MongoCmd = &cobra.Command{
	Use:   "mongo [mongodb_uri]",
	Short: "Set MongoDB URI",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		err := config.Load()
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		if len(args) == 0 {
			if config.Config.MongoUri == "" {
				fmt.Println(config.DefaultMongoUri)
			} else {
				fmt.Println(config.Config.MongoUri)
			}
			return
		}

		config.Config.MongoUri = args[0]

		err = config.Save()
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		logrus.WithFields(logrus.Fields{
			"mongo_uri": config.Config.MongoUri,
		}).Info("cmd: Set MongoDB URI")
	},
}
