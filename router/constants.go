package router

import (
	"text/template"
)

const redirectConfTempl = `# pritunl-zero redirect server environment
WEB_PORT={{.WebPort}}
PRIVATE_KEY={{.PrivateKey}}
SECRET={{.Secret}}
`

var (
	redirectConf = template.Must(
		template.New("redirect").Parse(redirectConfTempl))
)

type redirectConfData struct {
	WebPort    int
	PrivateKey string
	Secret     string
}
