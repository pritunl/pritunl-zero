package demo

import (
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/settings"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func IsDemo() bool {
	_ = settings.System.Demo
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
