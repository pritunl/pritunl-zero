package settings

var (
	registry = map[string]*group{}
)

type newFunc func() any
type updateFunc func(any)

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
