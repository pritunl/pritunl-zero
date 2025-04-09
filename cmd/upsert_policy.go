package cmd

import (
	"fmt"
	"os"

	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/spf13/cobra"
)

func init() {
	UpsertPolicyCmd.PersistentFlags().String(
		"name",
		"",
		"Policy name",
	)
	UpsertPolicyCmd.PersistentFlags().Bool(
		"enabled",
		true,
		"Enable policy",
	)
	UpsertPolicyCmd.PersistentFlags().StringSlice(
		"service",
		[]string{},
		"Policy services",
	)
	UpsertPolicyCmd.PersistentFlags().StringSlice(
		"role",
		[]string{},
		"Policy role",
	)
	UpsertPolicyCmd.PersistentFlags().String(
		"admin-secondary",
		"",
		"Admin secondary provider",
	)
	UpsertPolicyCmd.PersistentFlags().String(
		"user-secondary",
		"",
		"User secondary provider",
	)
	UpsertPolicyCmd.PersistentFlags().String(
		"proxy-secondary",
		"",
		"Proxy secondary provider",
	)
	UpsertPolicyCmd.PersistentFlags().String(
		"authority-secondary",
		"",
		"Authority secondary provider",
	)
	UpsertPolicyCmd.PersistentFlags().Bool(
		"admin-device-secondary",
		false,
		"Enable admin device secondary",
	)
	UpsertPolicyCmd.PersistentFlags().Bool(
		"user-device-secondary",
		false,
		"Enable user device secondary",
	)
	UpsertPolicyCmd.PersistentFlags().Bool(
		"proxy-device-secondary",
		false,
		"Enable proxy device secondary",
	)
	UpsertPolicyCmd.PersistentFlags().Bool(
		"authority-device-secondary",
		false,
		"Enable authority device secondary",
	)
	UpsertPolicyCmd.PersistentFlags().Bool(
		"authority-require-smart-card",
		false,
		"Require smart card for authority",
	)
	UpsertCmd.AddCommand(UpsertPolicyCmd)
}

var UpsertPolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Insert or update policy",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		InitMinimal()

		db := database.GetDatabase()
		defer db.Close()

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fmt.Fprintln(os.Stderr, "Policy name required")
			os.Exit(1)
		}

		pol, err := policy.GetOne(db, &bson.M{
			"name": name,
		})
		if err != nil {
			if _, ok := err.(*database.NotFoundError); ok {
				pol = nil
				err = nil
			} else {
				return
			}
		}

		isNew := false
		fields := set.NewSet()
		if pol == nil {
			isNew = true
			pol = &policy.Policy{
				Name: name,
			}
		}

		if cmd.Flags().Changed("enabled") {
			fields.Add("disabled")
			enabled, _ := cmd.Flags().GetBool("enabled")
			pol.Disabled = !enabled
		}

		roles, _ := cmd.Flags().GetStringSlice("role")
		if len(roles) > 0 {
			fields.Add("roles")
			pol.Roles = roles
		}

		adminSecondary, _ := cmd.Flags().GetString("admin-secondary")
		if adminSecondary != "" {
			fields.Add("admin_secondary")
			adminSecondaryId, e := primitive.ObjectIDFromHex(adminSecondary)
			if e != nil {
				fmt.Fprintf(
					os.Stderr,
					"Invalid admin secondary provider ID '%s'\n",
					adminSecondary,
				)
				os.Exit(1)
			}
			pol.AdminSecondary = adminSecondaryId
		}

		userSecondary, _ := cmd.Flags().GetString("user-secondary")
		if userSecondary != "" {
			fields.Add("user_secondary")
			userSecondaryId, e := primitive.ObjectIDFromHex(userSecondary)
			if e != nil {
				fmt.Fprintf(
					os.Stderr,
					"Invalid user secondary provider ID '%s'\n",
					userSecondary,
				)
				os.Exit(1)
			}
			pol.UserSecondary = userSecondaryId
		}

		proxySecondary, _ := cmd.Flags().GetString("proxy-secondary")
		if proxySecondary != "" {
			fields.Add("proxy_secondary")
			proxySecondaryId, e := primitive.ObjectIDFromHex(proxySecondary)
			if e != nil {
				fmt.Fprintf(
					os.Stderr,
					"Invalid proxy secondary provider ID '%s'\n",
					proxySecondary,
				)
				os.Exit(1)
			}
			pol.ProxySecondary = proxySecondaryId
		}

		authoritySecondary, _ := cmd.Flags().GetString("authority-secondary")
		if authoritySecondary != "" {
			fields.Add("authority_secondary")
			authoritySecondaryId, e := primitive.ObjectIDFromHex(
				authoritySecondary)
			if e != nil {
				fmt.Fprintf(
					os.Stderr,
					"Invalid authority secondary provider ID '%s'\n",
					authoritySecondary,
				)
				os.Exit(1)
			}
			pol.AuthoritySecondary = authoritySecondaryId
		}

		if cmd.Flags().Changed("admin-device-secondary") {
			fields.Add("admin_device_secondary")
			pol.AdminDeviceSecondary, _ = cmd.Flags().GetBool(
				"admin-device-secondary")
		}

		if cmd.Flags().Changed("user-device-secondary") {
			fields.Add("user_device_secondary")
			pol.UserDeviceSecondary, _ = cmd.Flags().GetBool(
				"user-device-secondary")
		}

		if cmd.Flags().Changed("proxy-device-secondary") {
			fields.Add("proxy_device_secondary")
			pol.ProxyDeviceSecondary, _ = cmd.Flags().GetBool(
				"proxy-device-secondary")
		}

		if cmd.Flags().Changed("authority-device-secondary") {
			fields.Add("authority_device_secondary")
			pol.AuthorityDeviceSecondary, _ = cmd.Flags().GetBool(
				"authority-device-secondary")
		}

		if cmd.Flags().Changed("authority-require-smart-card") {
			fields.Add("authority_require_smart_card")
			pol.AuthorityRequireSmartCard, _ = cmd.Flags().GetBool(
				"authority-require-smart-card")
		}

		errData, err := pol.Validate(db)
		if err != nil {
			return
		}

		if errData != nil {
			err = errData.GetError()
			return
		}

		if isNew {
			err = pol.Insert(db)
			if err != nil {
				return
			}
		} else {
			err = pol.CommitFields(db, fields)
			if err != nil {
				return
			}
		}

		_ = event.PublishDispatch(db, "policy.change")

		return
	},
}
