package utils

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

func GetRemoteAddr(c *gin.Context) (addr string) {
	addr = c.Request.Header.Get("CF-Connecting-IP")
	if addr != "" {
		return
	}

	addr = c.Request.Header.Get("X-Forwarded-For")
	if addr != "" {
		return
	}

	addr = c.Request.Header.Get("X-Real-Ip")
	if addr != "" {
		return
	}

	addr = strings.Split(c.Request.RemoteAddr, ":")[0]
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
