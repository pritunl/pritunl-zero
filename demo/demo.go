package demo

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/log"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/session"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/ssh"
	"github.com/pritunl/pritunl-zero/subscription"
	"github.com/pritunl/pritunl-zero/user"
	"github.com/pritunl/pritunl-zero/useragent"
	"github.com/pritunl/pritunl-zero/utils"
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

func BlockedSilent(c *gin.Context) bool {
	if !IsDemo() {
		return false
	}

	c.JSON(200, nil)
	return true
}

// Users
var Users = []*user.User{
	&user.User{
		Id:            utils.ObjectIdHex("5b6cd11857e4a9a88cbf072e"),
		Type:          "local",
		Provider:      primitive.ObjectID{},
		Username:      "demo",
		Token:         "",
		Secret:        "",
		LastActive:    time.Now(),
		LastSync:      time.Now(),
		Roles:         []string{"demo", "gitlab"},
		Administrator: "super",
		Disabled:      false,
		ActiveUntil:   time.Time{},
		Permissions:   []string{},
	},
	&user.User{
		Id:            utils.ObjectIdHex("5a7542190accad1a8a53b568"),
		Type:          "local",
		Provider:      primitive.ObjectID{},
		Username:      "pritunl",
		Token:         "",
		Secret:        "",
		LastActive:    time.Time{},
		LastSync:      time.Time{},
		Roles:         []string{},
		Administrator: "super",
		Disabled:      false,
		ActiveUntil:   time.Time{},
		Permissions:   []string{},
	},
}

var Agent = &useragent.Agent{
	OperatingSystem: useragent.Linux,
	Browser:         useragent.Chrome,
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
		Id:        utils.ObjectIdHex("5a17f9bf051a45ffacf2b352"),
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

var Sshcerts = []*ssh.Certificate{
	&ssh.Certificate{
		Id: utils.ObjectIdHex("5a180207051a45ffacf3b846"),
		AuthorityIds: []primitive.ObjectID{
			utils.ObjectIdHex("5a191ca03745632d533cf597"),
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

var Logs = []*log.Entry{
	&log.Entry{
		Id:        utils.ObjectIdHex("5a18e6ae051a45ffac0e5b67"),
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
		Id:        utils.ObjectIdHex("5a190b42051a45ffac129bbc"),
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

// Nodes
var Nodes = []*node.Node{
	&node.Node{
		Id:               utils.ObjectIdHex("5c74b2974ad0407c1ba1ab6e"),
		Name:             "pritunl-east0",
		Type:             "management_proxy_user",
		Timestamp:        time.Now(),
		Port:             80,
		NoRedirectServer: true,
		Protocol:         "http",
		Certificate:      primitive.ObjectID{},
		Certificates: []primitive.ObjectID{
			utils.ObjectIdHex("5a7544ae0accad1a8a53ba3d"),
		},
		ManagementDomain: "demo.zero.pritunl.com",
		UserDomain:       "user.demo.zero.pritunl.com",
		WebauthnDomain:   "zero.pritunl.com",
		EndpointDomain:   "demo.zero.pritunl.com",
		Services: []primitive.ObjectID{
			utils.ObjectIdHex("5b6cd0eb57e4a9a88cbf0678"),
		},
		Authorities:          []primitive.ObjectID{},
		RequestsMin:          32,
		ForwardedForHeader:   "X-Forwarded-For",
		ForwardedProtoHeader: "X-Forwarded-Proto",
		Memory:               25,
		Load1:                10,
		Load5:                15,
		Load15:               20,
		SoftwareVersion:      constants.Version,
		Hostname:             "pritunl-east0",
	},
}

// Policies
var Policies = []*policy.Policy{
	{
		Id:       utils.ObjectIdHex("67b8a03e4866ba90e6c45a8c"),
		Name:     "policy",
		Disabled: false,
		Roles: []string{
			"pritunl",
		},
		Services: []primitive.ObjectID{
			utils.ObjectIdHex("5b6cd0eb57e4a9a88cbf0678"),
		},
		Rules: map[string]*policy.Rule{
			"location": {
				Type:    "location",
				Disable: false,
				Values: []string{
					"US",
				},
			},
			"whitelist_networks": {
				Type:    "whitelist_networks",
				Disable: false,
				Values: []string{
					"10.0.0.0/8",
				},
			},
		},
		AdminSecondary:       primitive.ObjectID{},
		UserSecondary:        primitive.ObjectID{},
		AdminDeviceSecondary: true,
		UserDeviceSecondary:  true,
	},
}
