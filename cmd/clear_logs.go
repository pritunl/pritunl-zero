package cmd

import (
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ClearLogsCmd)
}

var ClearLogsCmd = &cobra.Command{
	Use:   "clear-logs",
	Short: "Clear logs",
	Run: func(cmd *cobra.Command, args []string) {
		Init()

		db := database.GetDatabase()
		defer db.Close()

		err := log.Clear(db)
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		logrus.Info("cmd: Logs cleared")
	},
}
