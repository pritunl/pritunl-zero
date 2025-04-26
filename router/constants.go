package router

import (
	"text/template"
)

const redirectConfTempl = `# pritunl-zero redirect server environment
WEB_PORT={{.WebPort}}
PUBLIC_KEY={{.PublicKey}}
KEY={{.Key}}
SECRET={{.Secret}}
`

var (
	redirectConf = template.Must(
		template.New("redirect").Parse(redirectConfTempl))
)

type redirectConfData struct {
	WebPort   int
	PublicKey string
	Key       string
	Secret    string
}
