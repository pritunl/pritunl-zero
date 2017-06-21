package settings

var System = &system{
	Id: "system",
}

func init() {
	register("system", System)
}

type system struct {
	Id              string `bson:"_id"`
	DatabaseVersion string `bson:"database_version"`
}
