package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(SetCmd)
}

var SetCmd = &cobra.Command{
	Use:   "set [group] [key] [value]",
	Short: "Change internal setting",
	Run: func(cmd *cobra.Command, args []string) {
		Init()

		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Missing required args")
			os.Exit(1)
		}

		db := database.GetDatabase()
		defer db.Close()

		group := args[0]
		key := args[1]
		val := args[2]

		var valParsed interface{}
		err := json.Unmarshal([]byte(val), &valParsed)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "cmd.settings: Failed to parse value"),
			}
			cobra.CheckErr(err)
			return
		}

		err = settings.Set(db, group, key, valParsed)
		if err != nil {
			cobra.CheckErr(err)
			return
		}
	},
}
