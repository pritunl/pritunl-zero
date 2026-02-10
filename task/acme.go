package task

import (
	"github.com/pritunl/pritunl-zero/acme"
	"github.com/pritunl/pritunl-zero/certificate"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/sirupsen/logrus"
)

var acmeRenew = &Task{
	Name:    "acme_renew",
	Version: 1,
	Hours:   []int{7},
	Minutes: []int{45},
	Handler: acmeRenewHandler,
}

func acmeRenewHandler(db *database.Database) (err error) {
	certs, err := certificate.GetAll(db)
	if err != nil {
		return
	}

	for _, cert := range certs {
		if cert.Type != certificate.LetsEncrypt {
			continue
		}

		err = acme.Renew(db, cert, false)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"certificate_id":   cert.Id.Hex(),
				"certificate_name": cert.Name,
				"error":            err,
			}).Error("task: Failed to renew certificate")
			continue
		}
	}
	return
}

func init() {
	register(acmeRenew)
}
