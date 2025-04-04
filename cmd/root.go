package cmd

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "pritunl-zero",
	Short: "Pritunl Zero Command Line Tool",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "cmd: Failed to execute help command"),
			}
			cobra.CheckErr(err)
		}
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "cmd: Failed to execute root command"),
		}
		cobra.CheckErr(err)
	}
}
