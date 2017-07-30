package utils

import (
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
)

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

func GetRemoteAddr(r *http.Request) (addr string) {
	addr = strings.Split(r.RemoteAddr, ":")[0]
	return
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
