package cmd

import (
	"fmt"
	"os"

	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(UnsetCmd)
}

var UnsetCmd = &cobra.Command{
	Use:   "unset [group] [key]",
	Short: "Reset internal setting to default",
	Run: func(cmd *cobra.Command, args []string) {
		Init()

		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Missing required args")
			os.Exit(1)
		}

		db := database.GetDatabase()
		defer db.Close()

		group := args[0]
		key := args[1]

		err := settings.Unset(db, group, key)
		if err != nil {
			cobra.CheckErr(err)
			return
		}
	},
}
