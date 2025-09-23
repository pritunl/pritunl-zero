package cmd

import (
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ResetNodeWeb)
}

var ResetNodeWeb = &cobra.Command{
	Use:   "reset-node-web",
	Short: "Reset node web server settings",
	Run: func(cmd *cobra.Command, args []string) {
		Init()

		db := database.GetDatabase()
		defer db.Close()

		err := config.Load()
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		ndeId, err := bson.ObjectIDFromHex(config.Config.NodeId)
		if err != nil || ndeId.IsZero() {
			err = nil
			logrus.Info("cmd: Node not initialized")
			cobra.CheckErr(err)
			return
		}

		coll := db.Nodes()

		_, err = coll.UpdateOne(db, &bson.M{
			"_id": ndeId,
		}, &bson.M{
			"$set": &bson.M{
				"type":               "management",
				"port":               443,
				"protocol":           "https",
				"no_redirect_server": false,
				"management_domain":  "",
				"user_domain":        "",
				"webauthn_domain":    "",
				"endpoint_domain":    "",
				"services":           []bson.ObjectID{},
			},
		})
		if err != nil {
			err = database.ParseError(err)
			cobra.CheckErr(err)
			return
		}

		logrus.WithFields(logrus.Fields{
			"node_id": config.Config.NodeId,
		}).Info("cmd: Node web server reset")

		logger.Init()
	},
}
