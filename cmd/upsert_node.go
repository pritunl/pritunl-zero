package cmd

import (
	"fmt"
	"os"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/spf13/cobra"
)

type Node struct {
	Name                 string `bson:"name" json:"name"`
	Port                 int    `bson:"port" json:"port"`
	NoRedirectServer     bool   `bson:"no_redirect_server" json:"no_redirect_server"`
	Protocol             string `bson:"protocol" json:"protocol"`
	ManagementDomain     string `bson:"management_domain" json:"management_domain"`
	UserDomain           string `bson:"user_domain" json:"user_domain"`
	WebauthnDomain       string `bson:"webauthn_domain" json:"webauthn_domain"`
	EndpointDomain       string `bson:"endpoint_domain" json:"endpoint_domain"`
	ForwardedForHeader   string `bson:"forwarded_for_header" json:"forwarded_for_header"`
	ForwardedProtoHeader string `bson:"forwarded_proto_header" json:"forwarded_proto_header"`
	Hostname             string `bson:"hostname" json:"hostname"`
}

func init() {
	UpsertNodeCmd.PersistentFlags().String(
		"name",
		"",
		"Node name",
	)
	UpsertNodeCmd.PersistentFlags().Int(
		"port",
		0,
		"Node port",
	)
	UpsertNodeCmd.PersistentFlags().Bool(
		"no-redirect-server",
		false,
		"Disable redirect server",
	)
	UpsertNodeCmd.PersistentFlags().String(
		"protocol",
		"",
		"Node protocol",
	)
	UpsertNodeCmd.PersistentFlags().String(
		"management-domain",
		"",
		"Management domain",
	)
	UpsertNodeCmd.PersistentFlags().String(
		"user-domain",
		"",
		"User domain",
	)
	UpsertNodeCmd.PersistentFlags().String(
		"webauthn-domain",
		"",
		"WebAuthn domain",
	)
	UpsertNodeCmd.PersistentFlags().String(
		"endpoint-domain",
		"",
		"Endpoint domain",
	)
	UpsertNodeCmd.PersistentFlags().String(
		"forwarded-for-header",
		"",
		"Forwarded for header",
	)
	UpsertNodeCmd.PersistentFlags().String(
		"forwarded-proto-header",
		"",
		"Forwarded proto header",
	)
	UpsertNodeCmd.PersistentFlags().String(
		"hostname",
		"",
		"Node hostname",
	)
	UpsertCmd.AddCommand(UpsertNodeCmd)
}

var UpsertNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Update node",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		InitMinimal()

		db := database.GetDatabase()
		defer db.Close()

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fmt.Fprintln(os.Stderr,
				"Node name required, use self for this node")
			os.Exit(1)
		}

		var nde *node.Node
		if name == "self" {
			objId, e := primitive.ObjectIDFromHex(config.Config.NodeId)
			if e != nil {
				err = &errortypes.ParseError{
					errors.Wrap(e, "cmd: Failed to parse ObjectId"),
				}
				return
			}

			nde, err = node.GetOne(db, &bson.M{
				"_id": objId,
			})
			if err != nil {
				if _, ok := err.(*database.NotFoundError); ok {
					nde = nil
					err = nil
				} else {
					return
				}
			}
		} else {
			nde, err = node.GetOne(db, &bson.M{
				"name": name,
			})
			if err != nil {
				if _, ok := err.(*database.NotFoundError); ok {
					nde = nil
					err = nil
				} else {
					return
				}
			}
		}

		if nde == nil {
			fmt.Fprintln(os.Stderr, "Failed to find node")
			os.Exit(1)
		}

		fields := set.NewSet()

		if cmd.Flags().Changed("port") {
			fields.Add("port")
			port, _ := cmd.Flags().GetInt("port")
			if port < 1 || port > 65535 {
				fmt.Fprintln(os.Stderr, "Invalid port number")
				os.Exit(1)
			}
			nde.Port = port
		}

		if cmd.Flags().Changed("no-redirect-server") {
			fields.Add("no_redirect_server")
			nde.NoRedirectServer, _ = cmd.Flags().GetBool(
				"no-redirect-server")
		}

		if cmd.Flags().Changed("protocol") {
			fields.Add("protocol")
			protocol, _ := cmd.Flags().GetString("protocol")
			nde.Protocol = protocol
		}

		if cmd.Flags().Changed("management-domain") {
			fields.Add("management_domain")
			nde.ManagementDomain, _ = cmd.Flags().GetString(
				"management-domain")
		}

		if cmd.Flags().Changed("user-domain") {
			fields.Add("user_domain")
			nde.UserDomain, _ = cmd.Flags().GetString("user-domain")
		}

		if cmd.Flags().Changed("webauthn-domain") {
			fields.Add("webauthn_domain")
			nde.WebauthnDomain, _ = cmd.Flags().GetString("webauthn-domain")
		}

		if cmd.Flags().Changed("endpoint-domain") {
			fields.Add("endpoint_domain")
			nde.EndpointDomain, _ = cmd.Flags().GetString("endpoint-domain")
		}

		if cmd.Flags().Changed("forwarded-for-header") {
			fields.Add("forwarded_for_header")
			nde.ForwardedForHeader, _ = cmd.Flags().GetString(
				"forwarded-for-header")
		}

		if cmd.Flags().Changed("forwarded-proto-header") {
			fields.Add("forwarded_proto_header")
			nde.ForwardedProtoHeader, _ = cmd.Flags().GetString(
				"forwarded-proto-header")
		}

		if cmd.Flags().Changed("hostname") {
			fields.Add("hostname")
			nde.Hostname, _ = cmd.Flags().GetString("hostname")
		}

		err = nde.CommitFields(db, fields)
		if err != nil {
			return
		}

		_ = event.PublishDispatch(db, "node.change")

		return
	},
}
