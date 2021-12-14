package search

var (
	mappingsRegistry = map[string][]*Mapping{}
)

func AddMappings(index string, mappings []*Mapping) {
	mappingsRegistry[index] = mappings
}
