package settings

var Elastic *elastic

type elastic struct {
	Id            string   `bson:"_id"`
	Addresses     []string `bson:"addresses"`
	Username      string   `bson:"username"`
	Password      string   `bson:"password"`
	ProxyRequests bool     `bson:"proxy_requests"`
	BufferLength  int      `bson:"buffer_length" default:"1024"`
	BufferSize    int      `bson:"buffer_size" default:"536870912"`
	GroupLength   int      `bson:"group_length" default:"100"`
	GroupSize     int      `bson:"group_size" default:"10485760"`
}

func newElastic() interface{} {
	return &elastic{
		Id: "elastic",
	}
}

func updateElastic(data interface{}) {
	Elastic = data.(*elastic)
}

func init() {
	register("elastic", newElastic, updateElastic)
}
