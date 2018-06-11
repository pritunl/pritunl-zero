package settings

var Local *local

type local struct {
	Facets []string
}

func init() {
	Local = &local{}
}
