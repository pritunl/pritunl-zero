package cmd

import (
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(DisablePoliciesCmd)
}

var DisablePoliciesCmd = &cobra.Command{
	Use:   "disable-policies",
	Short: "Disable all active policies",
	Run: func(cmd *cobra.Command, args []string) {
		InitMinimal()

		db := database.GetDatabase()
		defer db.Close()

		coll := db.Policies()

		_, err := coll.UpdateMany(db, &bson.M{}, &bson.M{
			"$set": &bson.M{
				"disabled": true,
			},
		})
		if err != nil {
			if _, ok := err.(*database.NotFoundError); ok {
				err = nil
			} else {
				cobra.CheckErr(err)
				return
			}
		}

		logrus.Info("cmd: Policies disabled")
	},
}
