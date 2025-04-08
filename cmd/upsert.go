package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(UpsertCmd)
}

var UpsertCmd = &cobra.Command{
	Use:   "upsert",
	Short: "Update or insert a resource",
}
