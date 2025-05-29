package settings

var Router *router

type router struct {
	Id                     string `bson:"_id"`
	ReadTimeout            int    `bson:"read_timeout" default:"300"`
	HeaderTimeout          int    `bson:"header_timeout" default:"60"`
	WriteTimeout           int    `bson:"write_timeout" default:"300"`
	IdleTimeout            int    `bson:"idle_timeout" default:"60"`
	DialTimeout            int    `bson:"dial_timeout" default:"60"`
	RequestTimeout         int    `bson:"request_timeout" default:"90"`
	DialKeepAlive          int    `bson:"dial_keep_alive" default:"60"`
	MaxIdleConns           int    `bson:"max_idle_conns" default:"1000"`
	MaxIdleConnsPerHost    int    `bson:"max_idle_conns_per_host" default:"100"`
	IdleConnTimeout        int    `bson:"idle_conn_timeout" default:"90"`
	HandshakeTimeout       int    `bson:"handshake_timeout" default:"10"`
	ContinueTimeout        int    `bson:"continue_timeout" default:"10"`
	DisableIdleConnections bool   `bson:"disable_idle_connections"`
	MaxHeaderBytes         int    `bson:"max_header_bytes" default:"10485760"`
	MaxResponseHeaderBytes int    `bson:"max_response_header_bytes" default:"33554432"`
	ForceRedirectSystemd   bool   `bson:"force_redirect_systemd"`
	RootDomainCookie       bool   `bson:"root_domain_cookie"`
	UnsafeRequests         bool   `bson:"unsafe_requests"`
	UnsafeRemoteHeader     bool   `bson:"unsafe_remote_header"`
	DebugWebSocket         bool   `bson:"debug_websocket"`
	SkipVerify             bool   `bson:"skip_verify"`
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
