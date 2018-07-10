package settings

var Local *local

type local struct {
	AppId        string
	Facets       []string
	HasLocalAuth bool
}

func init() {
	Local = &local{}
}
