package version

type Version struct {
	Module  string `bson:"_id,omitempty" json:"id"`
	Version int    `bson:"version,omitempty" json:"version"`
}
