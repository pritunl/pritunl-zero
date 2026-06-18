package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pritunl/pritunl-zero/settings"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(GetCmd)
}

var GetCmd = &cobra.Command{
	Use:   "get [group] [key]",
	Short: "Print current internal setting",
	Run: func(cmd *cobra.Command, args []string) {
		Init()

		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Missing required args")
			os.Exit(1)
		}

		group := args[0]
		key := args[1]

		val, err := settings.Value(group, key)
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		output, err := json.Marshal(val)
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		fmt.Println(string(output))
	},
}
