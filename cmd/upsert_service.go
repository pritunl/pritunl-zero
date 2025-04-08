package cmd

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/service"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/spf13/cobra"
)

func init() {
	UpsertServiceCmd.PersistentFlags().String(
		"name",
		"",
		"Service name",
	)
	UpsertServiceCmd.PersistentFlags().String(
		"type",
		"",
		"Service type",
	)
	UpsertServiceCmd.PersistentFlags().StringSlice(
		"role",
		[]string{},
		"Service role",
	)
	UpsertServiceCmd.PersistentFlags().StringSlice(
		"domain",
		[]string{},
		"Service external domain",
	)
	UpsertServiceCmd.PersistentFlags().StringSlice(
		"server",
		[]string{},
		"Service internal server",
	)
	UpsertServiceCmd.PersistentFlags().Bool(
		"http2",
		false,
		"Enable HTTP/2 support",
	)
	UpsertServiceCmd.PersistentFlags().Bool(
		"share-session",
		false,
		"Share sessions between services",
	)
	UpsertServiceCmd.PersistentFlags().String(
		"logout-path",
		"",
		"Custom logout path",
	)
	UpsertServiceCmd.PersistentFlags().Bool(
		"websockets",
		false,
		"Enable WebSockets support",
	)
	UpsertServiceCmd.PersistentFlags().Bool(
		"disable-csrf-check",
		false,
		"Disable CSRF check",
	)
	UpsertServiceCmd.PersistentFlags().StringSlice(
		"whitelist-network",
		[]string{},
		"Whitelist networks",
	)
	UpsertServiceCmd.PersistentFlags().Bool(
		"whitelist-options",
		false,
		"Enable whitelist options",
	)
	UpsertCmd.AddCommand(UpsertServiceCmd)
}

var UpsertServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Insert or update service",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		InitMinimal()

		db := database.GetDatabase()
		defer db.Close()

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fmt.Fprintln(os.Stderr, "Service name required")
			os.Exit(1)
		}

		srvc, err := service.GetOne(db, &bson.M{
			"name": name,
		})
		if err != nil {
			if _, ok := err.(*database.NotFoundError); ok {
				srvc = nil
				err = nil
			} else {
				return
			}
		}

		isNew := false
		fields := set.NewSet()
		if srvc == nil {
			isNew = true
			srvc = &service.Service{
				Name: name,
			}
		}

		srvcType, _ := cmd.Flags().GetString("type")
		if srvcType != "" {
			fields.Add("type")
			switch srvcType {
			case service.Http:
				srvc.Type = service.Http
				break
			default:
				fmt.Fprintln(os.Stderr, "Service type invalid")
				os.Exit(1)
			}
		}

		roles, _ := cmd.Flags().GetStringSlice("role")
		if len(roles) > 0 {
			fields.Add("roles")
			srvc.Roles = roles
		}

		domains, _ := cmd.Flags().GetStringSlice("domain")
		if len(domains) > 0 {
			fields.Add("domains")
			srvc.Domains = []*service.Domain{}
			for _, domain := range domains {
				domainSpl := strings.Split(domain, "@")

				domn := &service.Domain{
					Domain: domainSpl[0],
				}

				if len(domainSpl) > 1 {
					domn.Host = domainSpl[1]
				}

				srvc.Domains = append(srvc.Domains, domn)
			}
		}

		servers, _ := cmd.Flags().GetStringSlice("server")
		if len(servers) > 0 {
			fields.Add("servers")
			srvc.Servers = []*service.Server{}
			for _, server := range servers {
				parsedServer, e := url.Parse(server)
				if e != nil {
					fmt.Fprintf(
						os.Stderr,
						"Service server '%s' invalid\n",
						server,
					)
					os.Exit(1)
				}

				serv := &service.Server{}

				switch parsedServer.Scheme {
				case service.Http:
					serv.Protocol = service.Http
					break
				case service.Https:
					serv.Protocol = service.Https
					break
				default:
					fmt.Fprintf(
						os.Stderr,
						"Service server '%s' invalid protocol\n",
						server,
					)
					os.Exit(1)
				}

				if parsedServer.Host == "" {
					fmt.Fprintf(
						os.Stderr,
						"Service server '%s' invalid host\n",
						server,
					)
					os.Exit(1)
				}

				host, port := utils.SplitHostPort(parsedServer.Host)
				if host == "" {
					fmt.Fprintf(
						os.Stderr,
						"Service server '%s' invalid host\n",
						server,
					)
					os.Exit(1)
				}
				serv.Hostname = host

				if port < 1 || port > 65535 {
					fmt.Fprintf(
						os.Stderr,
						"Service server '%s' invalid port\n",
						server,
					)
					os.Exit(1)
				}
				serv.Port = port

				srvc.Servers = append(srvc.Servers, serv)
			}
		}

		if cmd.Flags().Changed("http2") {
			fields.Add("http2")
			srvc.Http2, _ = cmd.Flags().GetBool("http2")
		}

		if cmd.Flags().Changed("share-session") {
			fields.Add("share_session")
			srvc.ShareSession, _ = cmd.Flags().GetBool("share-session")
		}

		if cmd.Flags().Changed("logout-path") {
			fields.Add("logout_path")
			srvc.LogoutPath, _ = cmd.Flags().GetString("logout-path")
		}

		if cmd.Flags().Changed("websockets") {
			fields.Add("websockets")
			srvc.WebSockets, _ = cmd.Flags().GetBool("websockets")
		}

		if cmd.Flags().Changed("disable-csrf-check") {
			fields.Add("disable_csrf_check")
			srvc.DisableCsrfCheck, _ = cmd.Flags().GetBool(
				"disable-csrf-check")
		}

		if cmd.Flags().Changed("whitelist-options") {
			fields.Add("whitelist_options")
			srvc.WhitelistOptions, _ = cmd.Flags().GetBool(
				"whitelist-options")
		}

		whitelistNets, _ := cmd.Flags().GetStringSlice("whitelist-network")
		if len(whitelistNets) > 0 {
			fields.Add("whitelist_networks")
			srvc.WhitelistNetworks = []string{}
			for _, whitelistNet := range whitelistNets {
				_, ipNet, e := net.ParseCIDR(whitelistNet)
				if e != nil {
					fmt.Fprintf(
						os.Stderr,
						"Service server '%s' invalid whitelist network\n",
						whitelistNet,
					)
					os.Exit(1)
				}

				srvc.WhitelistNetworks = append(
					srvc.WhitelistNetworks, ipNet.String())
			}
		}

		errData, err := srvc.Validate(db)
		if err != nil {
			return
		}

		if errData != nil {
			err = errData.GetError()
			return
		}

		if isNew {
			err = srvc.Insert(db)
			if err != nil {
				return
			}
		} else {
			err = srvc.CommitFields(db, fields)
			if err != nil {
				return
			}
		}

		_ = event.PublishDispatch(db, "service.change")

		return
	},
}
