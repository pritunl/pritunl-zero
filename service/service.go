package service

import (
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"gopkg.in/mgo.v2/bson"
	"net"
	"sort"
)

type Domain struct {
	Domain string `bson:"domain" json:"domain"`
	Host   string `bson:"host" json:"host"`
}

type Server struct {
	Protocol string `bson:"protocol" json:"protocol"`
	Hostname string `bson:"hostname" json:"hostname"`
	Port     int    `bson:"port" json:"port"`
}

type Service struct {
	Id                bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name              string        `bson:"name" json:"name"`
	Type              string        `bson:"type" json:"type"`
	ShareSession      bool          `bson:"share_session" json:"share_session"`
	LogoutPath        string        `bson:"logout_path" json:"logout_path"`
	WebSockets        bool          `bson:"websockets" json:"websockets"`
	DisableCsrfCheck  bool          `bson:"disable_csrf_check" json:"disable_csrf_check"`
	Domains           []*Domain     `bson:"domains" json:"domains"`
	Roles             []string      `bson:"roles" json:"roles"`
	Servers           []*Server     `bson:"servers" json:"servers"`
	WhitelistNetworks []string      `bson:"whitelist_networks" json:"whitelist_networks"`
}

func (s *Service) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	if s.Type == "" {
		s.Type = Http
	}

	if s.Domains == nil {
		s.Domains = []*Domain{}
	}

	if s.Roles == nil {
		s.Roles = []string{}
	}

	if s.Servers == nil {
		s.Servers = []*Server{}
	}

	if s.WhitelistNetworks == nil {
		s.WhitelistNetworks = []string{}
	}

	for _, server := range s.Servers {
		if server.Protocol != "http" && server.Protocol != "https" {
			errData = &errortypes.ErrorData{
				Error:   "service_protocol_invalid",
				Message: "Invalid service server protocol",
			}
			return
		}

		if server.Hostname == "" {
			errData = &errortypes.ErrorData{
				Error:   "service_hostname_invalid",
				Message: "Invalid service server hostname",
			}
			return
		}

		if server.Port == 0 {
			errData = &errortypes.ErrorData{
				Error:   "service_port_invalid",
				Message: "Invalid service server port",
			}
			return
		}
	}

	for _, cidr := range s.WhitelistNetworks {
		_, _, err = net.ParseCIDR(cidr)
		if err != nil {
			err = nil
			errData = &errortypes.ErrorData{
				Error:   "whitelist_network_invalid",
				Message: "Whitelist network not a valid subnet",
			}
			return
		}
	}

	s.Format()

	return
}

func (s *Service) Format() {
	sort.Strings(s.Roles)
	sort.Strings(s.WhitelistNetworks)
}

func (s *Service) Commit(db *database.Database) (err error) {
	coll := db.Services()

	err = coll.Commit(s.Id, s)
	if err != nil {
		return
	}

	return
}

func (s *Service) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Services()

	err = coll.CommitFields(s.Id, s, fields)
	if err != nil {
		return
	}

	return
}

func (s *Service) Insert(db *database.Database) (err error) {
	coll := db.Services()

	if s.Id != "" {
		err = &errortypes.DatabaseError{
			errors.New("service: Service already exists"),
		}
		return
	}

	err = coll.Insert(s)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
