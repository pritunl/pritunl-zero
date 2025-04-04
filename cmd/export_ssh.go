package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/authority"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

type exportData struct {
	Keys []string `json:"keys"`
}

func init() {
	RootCmd.AddCommand(ExportSshCmd)
}

var ExportSshCmd = &cobra.Command{
	Use:   "export-ssh [export_path]",
	Short: "Export SSH authorities for emergency client",
	Run: func(cmd *cobra.Command, args []string) {
		Init()

		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Missing required args")
			os.Exit(1)
		}

		outputPath := args[0]
		if outputPath == "" {
			fmt.Fprintln(os.Stderr, "Missing required args")
			os.Exit(1)
			return
		}

		db := database.GetDatabase()
		defer db.Close()

		fmt.Print("Enter encryption passphrase: ")
		passByt, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			err = &errortypes.ReadError{
				errors.Wrap(err, "cmd.export: Failed to read passphrase"),
			}
			cobra.CheckErr(err)
			return
		}
		pass := string(passByt)
		fmt.Println("")

		authrs, err := authority.GetAll(db)
		if err != nil {
			cobra.CheckErr(err)
			return
		}

		keys := []string{}

		for _, authr := range authrs {
			key, e := authr.Export(pass)
			if e != nil {
				err = e
				cobra.CheckErr(err)
				return
			}

			keys = append(keys, key)
		}

		data := &exportData{
			Keys: keys,
		}

		marhData, err := json.Marshal(data)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "cmd.export: Failed to marshal keys"),
			}
			cobra.CheckErr(err)
			return
		}

		err = ioutil.WriteFile(outputPath, marhData, 0600)
		if err != nil {
			err = &errortypes.WriteError{
				errors.Wrap(err, "cmd.export: Failed to write output file"),
			}
			cobra.CheckErr(err)
			return
		}

		fmt.Printf("Successfully exported keys to %s\n", outputPath)

	},
}
