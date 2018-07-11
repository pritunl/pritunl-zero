package settings

var Local *local

type local struct {
	AppId       string
	Facets      []string
	NoLocalAuth bool
}

func init() {
	Local = &local{}
}
