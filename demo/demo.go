package demo

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/log"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/settings"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func IsDemo() bool {
	return settings.System.Demo
}

func Blocked(c *gin.Context) bool {
	if !IsDemo() {
		return false
	}

	errData := &errortypes.ErrorData{
		Error:   "demo_unavailable",
		Message: "Not available in demo mode",
	}
	c.JSON(400, errData)

	return true
}

var Agent = &agent.Agent{
	OperatingSystem: agent.Linux,
	Browser:         agent.Chrome,
	Ip:              "8.8.8.8",
	Isp:             "Google",
	Continent:       "North America",
	ContinentCode:   "NA",
	Country:         "United States",
	CountryCode:     "US",
	Region:          "Washington",
	RegionCode:      "WA",
	City:            "Seattle",
	Latitude:        47.611,
	Longitude:       -122.337,
}

var Audits = []*audit.Audit{
	&audit.Audit{
		Id:        bson.NewObjectId(),
		Timestamp: time.Unix(1498018860, 0),
		Type:      "admin_login",
		Fields: audit.Fields{
			"method": "local",
		},
		Agent: Agent,
	},
}

var Sessions = []*session.Session{
	&session.Session{
		Id:         "jhgRu4n3oY0iXRYmLb77Ql5jNs2o7uWM",
		Timestamp:  time.Unix(1498018860, 0),
		LastActive: time.Unix(1498018860, 0),
		Removed:    false,
		Agent:      Agent,
	},
}

var Logs = []*log.Entry{
	&log.Entry{
		Id:        bson.NewObjectId(),
		Level:     log.Info,
		Timestamp: time.Unix(1498018860, 0),
		Message:   "router: Starting redirect server",
		Stack:     "",
		Fields: map[string]interface{}{
			"port":       80,
			"production": true,
			"protocol":   "http",
		},
	},
	&log.Entry{
		Id:        bson.NewObjectId(),
		Level:     log.Info,
		Timestamp: time.Unix(1498018860, 0),
		Message:   "router: Starting web server",
		Stack:     "",
		Fields: map[string]interface{}{
			"port":       443,
			"production": true,
			"protocol":   "https",
		},
	},
}
