package settings

var Endpoint *endpoint

type endpoint struct {
	Id                string `bson:"_id"`
	Name              string `bson:"name"`
	KmsgDisplayLimit  int64  `bson:"kmsg_display_limit" default:"5000"`
	CheckDisplayLimit int64  `bson:"check_display_limit" default:"5000"`
}

func newEndpoint() interface{} {
	return &endpoint{
		Id: "endpoint",
	}
}

func updateEndpoint(data interface{}) {
	Endpoint = data.(*endpoint)
}

func init() {
	register("endpoint", newEndpoint, updateEndpoint)
}
