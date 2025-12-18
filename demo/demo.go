package demo

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/audit"
	"github.com/pritunl/pritunl-zero/certificate"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/log"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/policy"
	"github.com/pritunl/pritunl-zero/secret"
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
		Provider:      bson.ObjectID{},
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
		Provider:      bson.ObjectID{},
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
		AuthorityIds: []bson.ObjectID{
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
		Fields: map[string]any{
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
		Fields: map[string]any{
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
		Certificate:      bson.ObjectID{},
		Certificates: []bson.ObjectID{
			utils.ObjectIdHex("5a7544ae0accad1a8a53ba3d"),
		},
		ManagementDomain: "demo.zero.pritunl.com",
		UserDomain:       "user.demo.zero.pritunl.com",
		WebauthnDomain:   "zero.pritunl.com",
		EndpointDomain:   "demo.zero.pritunl.com",
		Services: []bson.ObjectID{
			utils.ObjectIdHex("5b6cd0eb57e4a9a88cbf0678"),
		},
		Authorities:          []bson.ObjectID{},
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
		Services: []bson.ObjectID{
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
		AdminSecondary:       bson.ObjectID{},
		UserSecondary:        bson.ObjectID{},
		AdminDeviceSecondary: true,
		UserDeviceSecondary:  true,
	},
}

// Certificates
var Certificates = []*certificate.Certificate{
	{
		Id:      utils.ObjectIdHex("67b89ef24866ba90e6c459e8"),
		Name:    "zero-pritunl-com",
		Comment: "",
		Type:    "lets_encrypt",
		Key: `-----BEGIN RSA PRIVATE KEY-----
MIIJKQIBAAKCAgEAx9Y3Lk2AwV6ap7L/Sx9XC5mXaUf8hvMmDbLBqDZ1Y7xKJM2h
zQ8Xm1rK9q0wzQC6qiL6xHmTpKWTzNVzGsQdM3/qNPLNA7W8PIYCzjkSe5X1YktY
vxldBxYxPRJxXk5S9P8dFYVmFFKF2bvJ5pSMLq9w3z3nTm3TQtRPqWx2Vk3DqV2D
QKmNtqJnhVqYvVKa3QpLLwz8xKqB1sPXLr4XqQ3bz3fLjLxPmYV5WxLhgdKLYZTv
YxQPLPTJkX3Pw4XD4Qs4CrKLW5bYsqYKQ7kKDXgJmTxYzZLjZKf4vSqLxqV5bDPY
rR2YxQ9TKLkYKVMpNtY5J9X2fWzyPSvXqXZfVx7D8xJzDY8YKPLXmvxKQZxLJxSx
zxHQzYKJpX3YmVfqYYmfYxXYzLmYxDzSxXqLvKxVqXxQDsPxQVKfKqQx5KvxsVqD
-----END RSA PRIVATE KEY-----`,
		Certificate: `-----BEGIN CERTIFICATE-----
MIIGGTCCBQGgAwIBAgISBXx9YmN2KQm9g3Y5XmKbvx9YMA0GCSqGSIb3DQEBCwUA
MDMxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQwwCgYDVQQD
EwNSMTEwHhcNMjUwODA4MDY0NzI3WhcNMjUxMTA2MDY0NzI2WjAcMRowGAYDVQQD
ExFjbG91ZC5wcml0dW5sLnJlZDCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoC
ggIBAMfWNy5NgMFemqey/0sfVwuZl2lH/IbzJg2ywag2dWO8SiTNoc0PF5tayvat
MM0AuKoi+sR5k6Slk8zVcxrEHTN/6jTyzQO1vDyGAs45EnuV9WJLWL8ZXQcWMT0S
cV5OUvT/HRWFZhRShdn5iQ2Sry6vcN8950Dt00LUT6lsdlZNw6ldg0CpjbaiZ4Va
mL1Smt0KSy8M/MSqgdbD1y6+F6kN2893y4y8T5mFeVsS4YHSi2GU72MUDyz0yZF9
z8OFw+ELOAqyi1uW2LKmCkO5Cg14CZk8WM2S42Sn+L0qi8aleWwz2K0dmMUPUyi5
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFBjCCAu6gAwIBAgIRAIp9PhPWLzDvI4a9KQdrNPgwDQYJKoZIhvcNAQELBQAw
TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh
cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMjQwMzEzMDAwMDAw
WhcNMjcwMzEyMjM1OTU5WjAzMQswCQYDVQQGEwJVUzEWMBQGA1UEChMNTGV0J3Mg
RW5jcnlwdDEMMAoGA1UEAxMDUjExMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAuoe8XBsAOcvKCs3UZxD5ATylTqVhyybKUvsVAbe5KPUoHu0nsyQYOWcJ
DAjs4DqwO3cOvfPlOVRBDE6uQdaZdN5R2+97/1i9qLcT9t4x1fJyyXJqC4N0lZxG
AGQUmfOx2SLZzaiSqhwmej/+71gFewiVgdtxD4774zEJuwm+UE1fj5F2PVqdnoPy
-----END CERTIFICATE-----`,
		Info: &certificate.Info{
			Hash:         "bba8a3941280c8466a6a2a723cc06f26",
			SignatureAlg: "SHA256-RSA",
			PublicKeyAlg: "RSA",
			Issuer:       "R11",
			IssuedOn:     time.Now(),
			ExpiresOn:    time.Now().Add(2160 * time.Hour),
			DnsNames: []string{
				"zero.pritunl.com",
				"user.zero.pritunl.com",
				"service.zero.pritunl.com",
			},
		},
		AcmeDomains: []string{
			"zero.pritunl.com",
			"user.zero.pritunl.com",
			"service.zero.pritunl.com",
		},
		AcmeType:   "acme_dns",
		AcmeAuth:   "acme_cloudflare",
		AcmeSecret: utils.ObjectIdHex("67b89e8d4866ba90e6c459ba"),
	},
}

// Secrets
var Secrets = []*secret.Secret{
	{
		Id:      utils.ObjectIdHex("67b89e8d4866ba90e6c459ba"),
		Name:    "cloudflare-pritunl-com",
		Comment: "",
		Type:    "cloudflare",
		Key:     "a7kX9mN2vP8Q-4jL6wS3tR5Y-uH1gF7dZ0xC-vB8nM",
		Value:   "",
		Region:  "",
		PublicKey: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAz4K8Lm3QvR7WxN5YdE2P
jX9TpQ6HgM1wV0nS4KaF3ZcB8LrY5UvO2JmN7XsPqI1AgK8EoH3RdWzM9LfY2VtN
kP4QxGsJ7YnR8LwVmT3AqZ5HvK2NdP1XoS8JgR4LmW7YxQ3VnH5TsK9PpL2MdX8Rg
vJ3KqN5WxT1LsM4HgY7RdP8NqV2JmK5XwL3TsR8YgN4HxP1LdK9VwQ2MsT3XpR7Y
nL8KgJ5WdH3TmR9XsL2PqN7VxK4MgT3HdJ8YwP2LsK5RxT1NqM4JgY7PxR8WsL3T
mK9XwN2HgJ5YdL3RsP8VqT2MxK4NhR3JdY8WwL2TsM5QxN1PqK4YgJ7RxP8VsT3M
PwIDAQAB
-----END PUBLIC KEY-----`,
	},
}
