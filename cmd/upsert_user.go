package cmd

import (
	"fmt"
	"os"

	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/spf13/cobra"
)

func init() {
	UpsertUserCmd.PersistentFlags().String(
		"name",
		"",
		"User name",
	)
	UpsertUserCmd.PersistentFlags().StringSlice(
		"role",
		[]string{},
		"User role",
	)
	UpsertCmd.AddCommand(UpsertUserCmd)
}

var UpsertUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Update user",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		InitMinimal()

		db := database.GetDatabase()
		defer db.Close()

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fmt.Fprintln(os.Stderr, "User name required")
			os.Exit(1)
		}

		usr, err := user.GetUsername(db, user.Local, name)
		if err != nil {
			if _, ok := err.(*database.NotFoundError); ok {
				usr = nil
				err = nil
			} else {
				return
			}
		}

		if usr == nil {
			fmt.Fprintln(os.Stderr, "Failed to find user")
			os.Exit(1)
		}

		fields := set.NewSet()

		roles, _ := cmd.Flags().GetStringSlice("role")
		if len(roles) > 0 {
			fields.Add("roles")
			usr.Roles = roles
		}

		errData, err := usr.Validate(db)
		if err != nil {
			return
		}

		if errData != nil {
			err = errData.GetError()
			return
		}

		err = usr.CommitFields(db, fields)
		if err != nil {
			return
		}

		_ = event.PublishDispatch(db, "user.change")

		return
	},
}
