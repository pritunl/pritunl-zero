package acme

import (
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/secret"
)

type DnsService interface {
	Connect(db *database.Database, secr *secret.Secret) (err error)
	DnsTxtUpsert(db *database.Database, domain, val string) (err error)
	DnsTxtDelete(db *database.Database, domain, val string) (err error)
}
