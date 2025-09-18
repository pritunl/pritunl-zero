package service

import (
	"net"
	"sort"
	"strings"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/pritunl/pritunl-zero/utils"
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

type WhitelistPath struct {
	Path     string `bson:"path" json:"path"`
	extMatch int
}

type Service struct {
	Id                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name               string             `bson:"name" json:"name"`
	Type               string             `bson:"type" json:"type"`
	Http2              bool               `bson:"http2" json:"http2"`
	ShareSession       bool               `bson:"share_session" json:"share_session"`
	LogoutPath         string             `bson:"logout_path" json:"logout_path"`
	WebSockets         bool               `bson:"websockets" json:"websockets"`
	DisableCsrfCheck   bool               `bson:"disable_csrf_check" json:"disable_csrf_check"`
	ClientAuthority    primitive.ObjectID `bson:"client_authority,omitempty" json:"client_authority"`
	Domains            []*Domain          `bson:"domains" json:"domains"`
	Roles              []string           `bson:"roles" json:"roles"`
	Servers            []*Server          `bson:"servers" json:"servers"`
	WhitelistNetworks  []string           `bson:"whitelist_networks" json:"whitelist_networks"`
	WhitelistPaths     []*WhitelistPath   `bson:"whitelist_paths" json:"whitelist_paths"`
	WhitelistOptions   bool               `bson:"whitelist_options" json:"whitelist_options"`
	logoutPathExtMatch int
}

func (s *Service) MatchLogoutPath(pth string) bool {
	if s.LogoutPath == "" {
		return false
	}

	if s.logoutPathExtMatch == 0 {
		if strings.Contains(s.LogoutPath, "*") ||
			strings.Contains(s.LogoutPath, "?") {

			s.logoutPathExtMatch = 2
		} else {
			s.logoutPathExtMatch = 1
		}
	}

	if s.logoutPathExtMatch == 2 {
		return utils.Match(s.LogoutPath, pth)
	} else {
		return pth == s.LogoutPath
	}
}

func (s *Service) MatchWhitelistPath(matchPth string) bool {
	if s.WhitelistPaths == nil || len(s.WhitelistPaths) == 0 {
		return false
	}

	for _, pth := range s.WhitelistPaths {
		if pth.Path == "" {
			continue
		}

		if pth.extMatch == 0 {
			if strings.Contains(pth.Path, "*") ||
				strings.Contains(pth.Path, "?") {

				pth.extMatch = 2
			} else {
				pth.extMatch = 1
			}
		}

		if pth.extMatch == 2 {
			if utils.Match(pth.Path, matchPth) {
				return true
			}
		} else {
			if matchPth == pth.Path {
				return true
			}
		}
	}

	return false
}

func (s *Service) RemoveWhitelistNetworks() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	s.WhitelistNetworks = []string{}
	err = s.CommitFields(db, set.NewSet("whitelist_networks"))
	if err != nil {
		return
	}

	return
}

func (s *Service) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	s.Name = utils.FilterName(s.Name)

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

	if s.WhitelistPaths == nil {
		s.WhitelistPaths = []*WhitelistPath{}
	}

	for _, domain := range s.Domains {
		wildcardCount := strings.Count(domain.Domain, "*")
		if wildcardCount > 1 {
			errData = &errortypes.ErrorData{
				Error:   "domain_wildcards_invalid",
				Message: "Only one wildcard supported in external domain",
			}
			return
		} else if wildcardCount == 1 {
			if domain.Domain[0] != '*' {
				errData = &errortypes.ErrorData{
					Error:   "domain_wildcard_invalid",
					Message: "External domain must start with wildcard",
				}
				return
			}
		}
	}

	for _, server := range s.Servers {
		if server.Protocol != Http && server.Protocol != Https {
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

		if server.Port < 1 || server.Port > 65535 {
			errData = &errortypes.ErrorData{
				Error:   "service_port_invalid",
				Message: "Invalid service server port",
			}
			return
		}
	}

	newWhitelistNetworks := []string{}
	for _, cidr := range s.WhitelistNetworks {
		_, ipNet, e := net.ParseCIDR(cidr)
		if e != nil {
			err = nil
			errData = &errortypes.ErrorData{
				Error:   "whitelist_network_invalid",
				Message: "Whitelist network not a valid subnet",
			}
			return
		}
		newWhitelistNetworks = append(newWhitelistNetworks, ipNet.String())
	}
	s.WhitelistNetworks = newWhitelistNetworks

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

	if !s.Id.IsZero() {
		err = &errortypes.DatabaseError{
			errors.New("service: Service already exists"),
		}
		return
	}

	_, err = coll.InsertOne(db, s)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func init() {
	module := requires.New("service")
	module.After("settings")

	// Fix v1.0.1317.75 service creation
	module.Handler = func() (err error) {
		db := database.GetDatabase()
		defer db.Close()

		coll := db.Services()

		_, err = coll.UpdateMany(db, &bson.M{
			"domains":            nil,
			"roles":              nil,
			"servers":            nil,
			"whitelist_networks": nil,
		}, &bson.M{
			"$set": &bson.M{
				"domains":            []interface{}{},
				"roles":              []interface{}{},
				"servers":            []interface{}{},
				"whitelist_networks": []interface{}{},
			},
		})
		if err != nil {
			err = database.ParseError(err)
			return
		}

		return
	}
}
