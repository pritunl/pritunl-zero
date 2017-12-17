package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/authority"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"syscall"
)

type exportData struct {
	Keys []string `json:"keys"`
}

func ExportSsh() (err error) {
	outputPath := flag.Arg(1)

	if outputPath == "" {
		err = &errortypes.ReadError{
			errors.Wrap(err, "cmd.export: Missing export path"),
		}
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
		return
	}
	pass := string(passByt)
	fmt.Println("")

	authrs, err := authority.GetAll(db)
	if err != nil {
		return
	}

	keys := []string{}

	for _, authr := range authrs {
		key, e := authr.Export(pass)
		if e != nil {
			err = e
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
		return
	}

	err = ioutil.WriteFile(outputPath, marhData, 0600)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "cmd.export: Failed to write output file"),
		}
		return
	}

	fmt.Printf("Successfully exported keys to %s\n", outputPath)

	return
}
