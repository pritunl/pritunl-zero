package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

var httpErrCodes = map[int]string{
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Payload Too Large",
	414: "URI Too Long",
	415: "Unsupported Media Type",
	416: "Range Not Satisfiable",
	417: "Expectation Failed",
	421: "Misdirected Request",
	422: "Unprocessable Entity",
	423: "Locked",
	424: "Failed Dependency",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	451: "Unavailable For Legal Reasons",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	510: "Not Extended",
	511: "Network Authentication Required",
}

func StripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}

func ParseObjectId(strId string) (objId bson.ObjectId, ok bool) {
	bytId, err := hex.DecodeString(strId)
	if err != nil {
		return
	}

	if len(bytId) != 12 {
		return
	}

	objId = bson.ObjectId(bytId)
	ok = true
	return
}

func AbortWithStatus(c *gin.Context, code int) {
	r := render.String{
		Format: fmt.Sprintf("%d %s", code, httpErrCodes[code]),
	}

	c.Status(code)
	r.WriteContentType(c.Writer)
	c.Writer.WriteHeaderNow()
	r.Render(c.Writer)
	c.Abort()
}

func AbortWithError(c *gin.Context, code int, err error) {
	AbortWithStatus(c, code)
	c.Error(err)
}
