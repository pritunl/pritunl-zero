package cmd

import (
	"fmt"
	"os"

	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/certificate"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/secret"
	"github.com/spf13/cobra"
)

func init() {
	UpsertCertificateCmd.PersistentFlags().String(
		"name",
		"",
		"Certificate name",
	)
	UpsertCertificateCmd.PersistentFlags().String(
		"comment",
		"",
		"Certificate comment",
	)
	UpsertCertificateCmd.PersistentFlags().String(
		"type",
		"",
		"Certificate type (text, lets_encrypt)",
	)
	UpsertCertificateCmd.PersistentFlags().String(
		"key",
		"",
		"Certificate key",
	)
	UpsertCertificateCmd.PersistentFlags().String(
		"certificate",
		"",
		"Certificate data",
	)
	UpsertCertificateCmd.PersistentFlags().StringSlice(
		"acme-domain",
		[]string{},
		"ACME certificate domain",
	)
	UpsertCertificateCmd.PersistentFlags().String(
		"acme-type",
		"",
		"ACME vertification method (http, dns)",
	)
	UpsertCertificateCmd.PersistentFlags().String(
		"acme-api",
		"",
		"ACME DNS provider (aws, cloudflare, oracle_cloud)",
	)
	UpsertCertificateCmd.PersistentFlags().String(
		"acme-secret",
		"",
		"Secret by name with API key for ACME DNS",
	)
	UpsertCmd.AddCommand(UpsertCertificateCmd)
}

var UpsertCertificateCmd = &cobra.Command{
	Use:   "certificate",
	Short: "Insert or update certificate",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		InitMinimal()

		db := database.GetDatabase()
		defer db.Close()

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fmt.Fprintln(os.Stderr, "Certificate name required")
			os.Exit(1)
		}

		cert, err := certificate.GetOne(db, &bson.M{
			"name": name,
		})
		if err != nil {
			if _, ok := err.(*database.NotFoundError); ok {
				cert = nil
				err = nil
			} else {
				return
			}
		}

		isNew := false
		fields := set.NewSet()
		if cert == nil {
			isNew = true
			cert = &certificate.Certificate{
				Name: name,
			}
		}

		if cmd.Flags().Changed("comment") {
			fields.Add("comment")
			cert.Comment, _ = cmd.Flags().GetString("comment")
		}

		if cmd.Flags().Changed("type") {
			certType, _ := cmd.Flags().GetString("type")
			if certType != "text" && certType != "lets_encrypt" {
				fmt.Fprintf(
					os.Stderr,
					"Invalid certificate type '%s', must be 'text' or 'lets_encrypt'\n",
					certType,
				)
				os.Exit(1)
			}
			fields.Add("type")
			cert.Type = certType
		}

		if cmd.Flags().Changed("key") {
			fields.Add("key")
			cert.Key, _ = cmd.Flags().GetString("key")
		}

		if cmd.Flags().Changed("certificate") {
			fields.Add("certificate")
			cert.Certificate, _ = cmd.Flags().GetString("certificate")
		}

		if cmd.Flags().Changed("acme-domain") {
			fields.Add("acme_domains")
			cert.AcmeDomains, _ = cmd.Flags().GetStringSlice("acme-domain")
		}

		if cmd.Flags().Changed("acme-type") {
			acmeType, _ := cmd.Flags().GetString("acme-type")

			fields.Add("acme_type")
			switch acmeType {
			case "http":
				cert.AcmeType = certificate.AcmeHTTP
			case "dns":
				cert.AcmeType = certificate.AcmeDNS
			case "":
				cert.AcmeAuth = ""
			default:
				fmt.Fprintf(
					os.Stderr,
					"Invalid ACME type '%s'\n",
					acmeType,
				)
				os.Exit(1)
			}
		}

		if cmd.Flags().Changed("acme-api") {
			acmeApi, _ := cmd.Flags().GetString("acme-api")

			fields.Add("acme_auth")
			switch acmeApi {
			case "aws":
				cert.AcmeAuth = certificate.AcmeAWS
			case "cloudflare":
				cert.AcmeAuth = certificate.AcmeCloudflare
			case "oracle_cloud":
				cert.AcmeAuth = certificate.AcmeOracleCloud
			case "":
				cert.AcmeAuth = ""
			default:
				fmt.Fprintf(
					os.Stderr,
					"Invalid ACME API '%s'\n",
					acmeApi,
				)
				os.Exit(1)
			}
		}

		if cmd.Flags().Changed("acme-secret") {
			acmeSecret, _ := cmd.Flags().GetString("acme-secret")

			fields.Add("acme_secret")
			if acmeSecret == "" {
				cert.AcmeSecret = primitive.NilObjectID
			} else {
				secr, e := secret.GetOne(db, &bson.M{
					"name": acmeSecret,
				})
				if e != nil {
					err = e
					if _, ok := err.(*database.NotFoundError); ok {
						fmt.Fprintf(os.Stderr,
							"Failed to find secret '%s'\n", acmeSecret)
						os.Exit(1)
					}
					return
				}

				cert.AcmeSecret = secr.Id
			}
		}

		errData, err := cert.Validate(db)
		if err != nil {
			return
		}

		if errData != nil {
			err = errData.GetError()
			return
		}

		if isNew {
			err = cert.Insert(db)
			if err != nil {
				return
			}
		} else {
			err = cert.CommitFields(db, fields)
			if err != nil {
				return
			}
		}

		_ = event.PublishDispatch(db, "certificate.change")

		return
	},
}
