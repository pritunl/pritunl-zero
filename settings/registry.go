package settings

var (
	registry = map[string]*group{}
)

type newFunc func() interface{}
type updateFunc func(interface{})

type group struct {
	New    newFunc
	Update updateFunc
}

func register(name string, new newFunc, update updateFunc) {
	grp := &group{
		New:    new,
		Update: update,
	}

	registry[name] = grp
}
