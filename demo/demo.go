package demo

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/log"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/ssh"
	"github.com/pritunl/pritunl-zero/subscription"
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

var auditId, _ = primitive.ObjectIDFromHex("5a17f9bf051a45ffacf2b352")
var Audits = []*audit.Audit{
	&audit.Audit{
		Id:        auditId,
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
		Type:       session.User,
		Timestamp:  time.Unix(1498018860, 0),
		LastActive: time.Unix(1498018860, 0),
		Removed:    false,
		Agent:      Agent,
	},
}

var sshId, _ = primitive.ObjectIDFromHex("5a180207051a45ffacf3b846")
var authrId, _ = primitive.ObjectIDFromHex("5a191ca03745632d533cf597")
var Sshcerts = []*ssh.Certificate{
	&ssh.Certificate{
		Id: sshId,
		AuthorityIds: []primitive.ObjectID{
			authrId,
		},
		Timestamp: time.Unix(1498018860, 0),
		CertificatesInfo: []*ssh.Info{
			&ssh.Info{
				Serial:  "2207385157562819502",
				Expires: time.Unix(1498105260, 0),
				Principals: []string{
					"demo",
				},
				Extensions: []string{
					"permit-X11-forwarding",
					"permit-agent-forwarding",
					"permit-port-forwarding",
					"permit-pty",
					"permit-user-rc",
				},
			},
		},
		Agent: Agent,
	},
}

var logId0, _ = primitive.ObjectIDFromHex("5a18e6ae051a45ffac0e5b67")
var logId1, _ = primitive.ObjectIDFromHex("5a190b42051a45ffac129bbc")
var Logs = []*log.Entry{
	&log.Entry{
		Id:        logId0,
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
		Id:        logId1,
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

var Subscription = &subscription.Subscription{
	Active:            true,
	Status:            "active",
	Plan:              "zero",
	Quantity:          1,
	Amount:            5000,
	PeriodEnd:         time.Unix(1893499200, 0),
	TrialEnd:          time.Time{},
	CancelAtPeriodEnd: false,
	Balance:           0,
	UrlKey:            "demo",
}
