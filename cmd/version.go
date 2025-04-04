package cmd

import (
	"fmt"

	"github.com/pritunl/pritunl-zero/constants"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(VersionCmd)
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show server version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("pritunl-zero v%s\n", constants.Version)
	},
}
