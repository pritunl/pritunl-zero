package cmd

import (
	"fmt"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(DefaultPasswordCmd)
}

var DefaultPasswordCmd = &cobra.Command{
	Use:   "default-password",
	Short: "Get default administrator password",
	Run: func(cmd *cobra.Command, args []string) {
		Init()

		db := database.GetDatabase()
		defer db.Close()

		usr, err := user.GetUsername(db, user.Local, "pritunl")
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		if usr.DefaultPassword == "" {
			err = &errortypes.NotFoundError{
				errors.New("cmd: No default password available"),
			}
			cobra.CheckErr(err)
			return
		}

		fmt.Println("Username: pritunl")
		fmt.Println("Password: " + usr.DefaultPassword)
	},
}
