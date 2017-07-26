package settings

var Router *router

type router struct {
	Id                  string `bson:"_id"`
	DialTimeout         int    `bson:"dial_timeout" default:"60"`
	DialKeepAlive       int    `bson:"dial_timeout" default:"60"`
	MaxIdleConns        int    `bson:"max_idle_conns" default:"1000"`
	MaxIdleConnsPerHost int    `bson:"max_idle_conns_per_host" default:"100"`
	IdleConnTimeout     int    `bson:"dial_timeout" default:"90"`
	HandshakeTimeout    int    `bson:"handshake_timeout" default:"10"`
	ContinueTimeout     int    `bson:"continue_timeout" default:"10"`
}

func newRouter() interface{} {
	return &router{
		Id: "router",
	}
}

func updateRouter(data interface{}) {
	Router = data.(*router)
}

func init() {
	register("router", newRouter, updateRouter)
}
