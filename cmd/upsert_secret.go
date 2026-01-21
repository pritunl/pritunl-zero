package cmd

import (
	"fmt"
	"os"

	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/secret"
	"github.com/spf13/cobra"
)

func init() {
	UpsertSecretCmd.PersistentFlags().String(
		"name",
		"",
		"Secret name",
	)
	UpsertSecretCmd.PersistentFlags().String(
		"comment",
		"",
		"Secret comment",
	)
	UpsertSecretCmd.PersistentFlags().String(
		"type",
		"",
		"Secret type",
	)
	UpsertSecretCmd.PersistentFlags().String(
		"aws-key",
		"",
		"AWS key ID",
	)
	UpsertSecretCmd.PersistentFlags().String(
		"aws-secret",
		"",
		"AWS secret ID",
	)
	UpsertSecretCmd.PersistentFlags().String(
		"aws-region",
		"",
		"AWS region",
	)
	UpsertSecretCmd.PersistentFlags().String(
		"cloudflare-token",
		"",
		"Cloudflare token",
	)
	UpsertSecretCmd.PersistentFlags().String(
		"oracle-tenancy",
		"",
		"Oracle Cloud tenancy OCID",
	)
	UpsertSecretCmd.PersistentFlags().String(
		"oracle-user",
		"",
		"Oracle Cloud user OCID",
	)
	UpsertSecretCmd.PersistentFlags().String(
		"oracle-region",
		"",
		"Oracle Cloud region",
	)
	UpsertSecretCmd.PersistentFlags().String(
		"gcp-credentials",
		"",
		"GCP service account JSON credentials",
	)
	UpsertCmd.AddCommand(UpsertSecretCmd)
}

var UpsertSecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Insert or update secret",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		InitMinimal()

		db := database.GetDatabase()
		defer db.Close()

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fmt.Fprintln(os.Stderr, "Secret name required")
			os.Exit(1)
		}

		secr, err := secret.GetOne(db, &bson.M{
			"name": name,
		})
		if err != nil {
			if _, ok := err.(*database.NotFoundError); ok {
				secr = nil
				err = nil
			} else {
				return
			}
		}

		isNew := false
		fields := set.NewSet()
		if secr == nil {
			isNew = true
			secr = &secret.Secret{
				Name: name,
			}
		}

		if cmd.Flags().Changed("comment") {
			fields.Add("comment")
			secr.Comment, _ = cmd.Flags().GetString("comment")
		}

		if cmd.Flags().Changed("type") {
			secrType, _ := cmd.Flags().GetString("type")
			fields.Add("type")
			switch secrType {
			case "aws":
				secr.Type = secret.AWS
			case "cloudflare":
				secr.Type = secret.Cloudflare
			case "oracle_cloud":
				secr.Type = secret.OracleCloud
			case "gcp":
				secr.Type = secret.GoogleCloud
			default:
				fmt.Fprintf(
					os.Stderr,
					"Invalid secret type '%s'\n",
					secrType,
				)
				os.Exit(1)
			}
		}

		if secr.Type == secret.AWS {
			if cmd.Flags().Changed("aws-key") {
				fields.Add("key")
				secr.Key, _ = cmd.Flags().GetString("aws-key")
			}
			if cmd.Flags().Changed("aws-secret") {
				fields.Add("value")
				secr.Value, _ = cmd.Flags().GetString("aws-secret")
			}
			if cmd.Flags().Changed("aws-region") {
				fields.Add("region")
				secr.Region, _ = cmd.Flags().GetString("aws-region")
			}
		}

		if secr.Type == secret.Cloudflare {
			if cmd.Flags().Changed("cloudflare-token") {
				fields.Add("key")
				secr.Key, _ = cmd.Flags().GetString("cloudflare-token")
			}
		}

		if secr.Type == secret.OracleCloud {
			if cmd.Flags().Changed("oracle-tenancy") {
				fields.Add("key")
				secr.Key, _ = cmd.Flags().GetString("oracle-tenancy")
			}
			if cmd.Flags().Changed("oracle-user") {
				fields.Add("value")
				secr.Value, _ = cmd.Flags().GetString("oracle-user")
			}
			if cmd.Flags().Changed("oracle-region") {
				fields.Add("region")
				secr.Region, _ = cmd.Flags().GetString("oracle-region")
			}
		}

		if secr.Type == secret.GoogleCloud {
			if cmd.Flags().Changed("google-service-key") {
				fields.Add("key")
				secr.Key, _ = cmd.Flags().GetString("google-service-key")
			}
		}

		errData, err := secr.Validate(db)
		if err != nil {
			return
		}

		if errData != nil {
			err = errData.GetError()
			return
		}

		if isNew {
			err = secr.Insert(db)
			if err != nil {
				return
			}
		} else {
			err = secr.CommitFields(db, fields)
			if err != nil {
				return
			}
		}

		_ = event.PublishDispatch(db, "secret.change")

		return
	},
}
