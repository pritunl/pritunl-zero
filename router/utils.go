package router

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/utils"
	"os"
)

func generateCert(certPath, keyPath string) (err error) {
	err = utils.Exec("",
		"openssl",
		"ecparam",
		"-name", "prime256v1",
		"-genkey",
		"-noout",
		"-out", keyPath,
	)
	if err != nil {
		return
	}

	err = os.Chmod(keyPath, 0600)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "router: Failed to generate certificate"),
		}
		return
	}

	err = utils.Exec("",
		"openssl",
		"req",
		"-new",
		"-batch",
		"-x509",
		"-days", "3652",
		"-key", keyPath,
		"-out", certPath,
	)
	if err != nil {
		return
	}

	return
}
