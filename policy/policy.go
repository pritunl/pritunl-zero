package policy

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/agent"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/user"
	"gopkg.in/mgo.v2/bson"
	"net"
	"net/http"
)

type Rule struct {
	Type    string   `bson:"type" json:"type"`
	Disable bool     `bson:"disable" json:"disable"`
	Values  []string `bson:"values" json:"values"`
}

type Policy struct {
	Id                        bson.ObjectId    `bson:"_id,omitempty" json:"id"`
	Name                      string           `bson:"name" json:"name"`
	Services                  []bson.ObjectId  `bson:"services" json:"services"`
	Authorities               []bson.ObjectId  `bson:"authorities" json:"authorities"`
	Roles                     []string         `bson:"roles" json:"roles"`
	Rules                     map[string]*Rule `bson:"rules" json:"rules"`
	AdminSecondary            bson.ObjectId    `bson:"admin_secondary,omitempty" json:"admin_secondary"`
	UserSecondary             bson.ObjectId    `bson:"user_secondary,omitempty" json:"user_secondary"`
	ProxySecondary            bson.ObjectId    `bson:"proxy_secondary,omitempty" json:"proxy_secondary"`
	AuthoritySecondary        bson.ObjectId    `bson:"authority_secondary,omitempty" json:"authority_secondary"`
	AdminDeviceSecondary      bool             `bson:"admin_device_secondary" json:"admin_device_secondary"`
	UserDeviceSecondary       bool             `bson:"user_device_secondary" json:"user_device_secondary"`
	ProxyDeviceSecondary      bool             `bson:"proxy_device_secondary" json:"proxy_device_secondary"`
	AuthorityDeviceSecondary  bool             `bson:"authority_device_secondary" json:"authority_device_secondary"`
	AuthorityRequireSmartCard bool             `bson:"authority_require_smart_card" json:"authority_require_smart_card"`
}

func (p *Policy) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if p.Services == nil {
		p.Services = []bson.ObjectId{}
	}

	services := []bson.ObjectId{}
	coll := db.Services()
	err = coll.Find(&bson.M{
		"_id": &bson.M{
			"$in": p.Services,
		},
	}).Distinct("_id", &services)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	p.Services = services

	if p.Authorities == nil {
		p.Authorities = []bson.ObjectId{}
	}

	authorities := []bson.ObjectId{}
	coll = db.Authorities()
	err = coll.Find(&bson.M{
		"_id": &bson.M{
			"$in": p.Authorities,
		},
	}).Distinct("_id", &authorities)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	p.Authorities = authorities

	if p.AdminSecondary != "" &&
		settings.Auth.GetSecondaryProvider(p.AdminSecondary) == nil {

		p.AdminSecondary = ""
	}
	if p.UserSecondary != "" &&
		settings.Auth.GetSecondaryProvider(p.UserSecondary) == nil {

		p.UserSecondary = ""
	}
	if p.ProxySecondary != "" &&
		settings.Auth.GetSecondaryProvider(p.ProxySecondary) == nil {

		p.ProxySecondary = ""
	}
	if p.AuthoritySecondary != "" &&
		settings.Auth.GetSecondaryProvider(p.AuthoritySecondary) == nil {

		p.AuthoritySecondary = ""
	}

	hasUserNode := false
	nodes, err := node.GetAll(db)
	if err != nil {
		return
	}

	for _, nde := range nodes {
		if nde.UserDomain != "" {
			hasUserNode = true
			break
		}
	}

	if (p.AdminDeviceSecondary || p.UserDeviceSecondary ||
		p.ProxyDeviceSecondary || p.AuthorityDeviceSecondary) &&
		!hasUserNode {

		errData = &errortypes.ErrorData{
			Error: "user_node_unavailable",
			Message: "At least one node must have a user domain configured " +
				"to use secondary device authentication",
		}
		return
	}

	return
}

func (p *Policy) ValidateUser(db *database.Database, usr *user.User,
	r *http.Request) (errData *errortypes.ErrorData, err error) {

	agnt, err := agent.Parse(db, r)
	if err != nil {
		return
	}

	for _, rule := range p.Rules {
		switch rule.Type {
		case OperatingSystem:
			match := false
			for _, value := range rule.Values {
				if value == agnt.OperatingSystem {
					match = true
					break
				}
			}

			if !match {
				if rule.Disable {
					errData = &errortypes.ErrorData{
						Error:   "unauthorized",
						Message: "Not authorized",
					}

					usr.Disabled = true
					err = usr.CommitFields(db, set.NewSet("disabled"))
					if err != nil {
						return
					}
				} else {
					errData = &errortypes.ErrorData{
						Error:   "operating_system_policy",
						Message: "Operating system not permitted",
					}
				}
				return
			}
			break
		case Browser:
			match := false
			for _, value := range rule.Values {
				if value == agnt.Browser {
					match = true
					break
				}
			}

			if !match {
				if rule.Disable {
					errData = &errortypes.ErrorData{
						Error:   "unauthorized",
						Message: "Not authorized",
					}

					usr.Disabled = true
					err = usr.CommitFields(db, set.NewSet("disabled"))
					if err != nil {
						return
					}
				} else {
					errData = &errortypes.ErrorData{
						Error:   "browser_policy",
						Message: "Browser not permitted",
					}
				}
				return
			}
			break
		case Location:
			match := false
			regionKey := fmt.Sprintf("%s_%s",
				agnt.CountryCode, agnt.RegionCode)

			for _, value := range rule.Values {
				if value == agnt.CountryCode || value == regionKey {
					match = true
					break
				}
			}

			if !match {
				if rule.Disable {
					errData = &errortypes.ErrorData{
						Error:   "unauthorized",
						Message: "Not authorized",
					}

					usr.Disabled = true
					err = usr.CommitFields(db, set.NewSet("disabled"))
					if err != nil {
						return
					}
				} else {
					errData = &errortypes.ErrorData{
						Error:   "location_policy",
						Message: "Location not permitted",
					}
				}
				return
			}
			break
		case WhitelistNetworks:
			match := false
			clientIp := net.ParseIP(agnt.Ip)

			for _, value := range rule.Values {
				_, network, e := net.ParseCIDR(value)
				if e != nil {
					err = &errortypes.ParseError{
						errors.Wrap(e, "policy: Failed to parse network"),
					}

					logrus.WithFields(logrus.Fields{
						"network": value,
						"error":   err,
					}).Error("policy: Invalid whitelist network")
					err = nil
					continue
				}

				if network.Contains(clientIp) {
					match = true
					break
				}
			}

			if !match {
				if rule.Disable {
					errData = &errortypes.ErrorData{
						Error:   "unauthorized",
						Message: "Not authorized",
					}

					usr.Disabled = true
					err = usr.CommitFields(db, set.NewSet("disabled"))
					if err != nil {
						return
					}
				} else {
					errData = &errortypes.ErrorData{
						Error:   "whitelist_networks_policy",
						Message: "Network not permitted",
					}
				}
				return
			}
			break
		case BlacklistNetworks:
			match := false
			clientIp := net.ParseIP(agnt.Ip)

			for _, value := range rule.Values {
				_, network, e := net.ParseCIDR(value)
				if e != nil {
					err = &errortypes.ParseError{
						errors.Wrap(e, "policy: Failed to parse network"),
					}

					logrus.WithFields(logrus.Fields{
						"network": value,
						"error":   err,
					}).Error("policy: Invalid blacklist network")
					err = nil
					continue
				}

				if network.Contains(clientIp) {
					match = true
					break
				}
			}

			if match {
				if rule.Disable {
					errData = &errortypes.ErrorData{
						Error:   "unauthorized",
						Message: "Not authorized",
					}

					usr.Disabled = true
					err = usr.CommitFields(db, set.NewSet("disabled"))
					if err != nil {
						return
					}
				} else {
					errData = &errortypes.ErrorData{
						Error:   "blacklist_networks_policy",
						Message: "Network not permitted",
					}
				}
				return
			}
			break
		}
	}

	return
}

func (p *Policy) Commit(db *database.Database) (err error) {
	coll := db.Policies()

	err = coll.Commit(p.Id, p)
	if err != nil {
		return
	}

	return
}

func (p *Policy) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Policies()

	err = coll.CommitFields(p.Id, p, fields)
	if err != nil {
		return
	}

	return
}

func (p *Policy) Insert(db *database.Database) (err error) {
	coll := db.Policies()

	if p.Id != "" {
		err = &errortypes.DatabaseError{
			errors.New("policy: Policy already exists"),
		}
		return
	}

	err = coll.Insert(p)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
