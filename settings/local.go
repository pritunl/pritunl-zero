package settings

var Local *local

type local struct {
	AppId       string
	Facets      []string
	NoLocalAuth bool
	DisableWeb  bool
	DisableMsg  string
}

func init() {
	Local = &local{}
}
