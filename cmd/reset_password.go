package cmd

import (
	"fmt"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ResetPasswordCmd)
}

var ResetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Reset default administrator password",
	Run: func(cmd *cobra.Command, args []string) {
		Init()

		db := database.GetDatabase()
		defer db.Close()

		coll := db.Users()

		_, err := coll.DeleteOne(db, &bson.M{
			"username": "pritunl",
		})
		if err != nil {
			if _, ok := err.(*database.NotFoundError); ok {
				err = nil
			} else {
				cobra.CheckErr(err)
				return
			}
		}

		usr := user.User{
			Type:          user.Local,
			Username:      "pritunl",
			Administrator: "super",
		}

		_, err = usr.Validate(db)
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		err = usr.GenerateDefaultPassword()
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		err = usr.Insert(db)
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		fmt.Println("Username: pritunl")
		fmt.Println("Password: " + usr.DefaultPassword)
	},
}
