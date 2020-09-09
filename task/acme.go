package task

import (
	"github.com/sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/acme"
	"github.com/pritunl/pritunl-zero/certificate"
	"github.com/pritunl/pritunl-zero/database"
)

var acmeRenew = &Task{
	Name:    "acme_renew",
	Hours:   []int{7},
	Mins:    []int{45},
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

		err = acme.Update(db, cert)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"certificate_id":   cert.Id.Hex(),
				"certificate_name": cert.Name,
			}).Warning("task: Failed to update certificate")
			continue
		}

		err = acme.Renew(db, cert)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"certificate_id":   cert.Id.Hex(),
				"certificate_name": cert.Name,
			}).Warning("task: Failed to renew certificate")
			continue
		}
	}

	return
}

func init() {
	register(acmeRenew)
}
