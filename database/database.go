package database

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/requires"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net"
	"net/url"
	"time"
)

var (
	Session *mgo.Session
)

type Database struct {
	session  *mgo.Session
	database *mgo.Database
}

func (d *Database) Close() {
	d.session.Close()
}

func (d *Database) getCollection(name string) (coll *Collection) {
	coll = &Collection{
		*d.database.C(name),
		d,
	}
	return
}

func (d *Database) Users() (coll *Collection) {
	coll = d.getCollection("users")
	return
}

func (d *Database) Services() (coll *Collection) {
	coll = d.getCollection("services")
	return
}

func (d *Database) Policies() (coll *Collection) {
	coll = d.getCollection("policies")
	return
}

func (d *Database) Sessions() (coll *Collection) {
	coll = d.getCollection("sessions")
	return
}

func (d *Database) Tasks() (coll *Collection) {
	coll = d.getCollection("tasks")
	return
}

func (d *Database) Tokens() (coll *Collection) {
	coll = d.getCollection("tokens")
	return
}

func (d *Database) CsrfTokens() (coll *Collection) {
	coll = d.getCollection("csrf_tokens")
	return
}

func (d *Database) Nonces() (coll *Collection) {
	coll = d.getCollection("nonces")
	return
}

func (d *Database) Settings() (coll *Collection) {
	coll = d.getCollection("settings")
	return
}

func (d *Database) Events() (coll *Collection) {
	coll = d.getCollection("events")
	return
}

func (d *Database) Nodes() (coll *Collection) {
	coll = d.getCollection("nodes")
	return
}

func (d *Database) Certificates() (coll *Collection) {
	coll = d.getCollection("certificates")
	return
}

func (d *Database) Authorities() (coll *Collection) {
	coll = d.getCollection("authorities")
	return
}

func (d *Database) SshChallenges() (coll *Collection) {
	coll = d.getCollection("ssh_challenges")
	return
}

func (d *Database) SshCertificates() (coll *Collection) {
	coll = d.getCollection("ssh_certificates")
	return
}

func (d *Database) KeybaseChallenges() (coll *Collection) {
	coll = d.getCollection("keybase_challenges")
	return
}

func (d *Database) AcmeChallenges() (coll *Collection) {
	coll = d.getCollection("acme_challenges")
	return
}

func (d *Database) Logs() (coll *Collection) {
	coll = d.getCollection("logs")
	return
}

func (d *Database) Audits() (coll *Collection) {
	coll = d.getCollection("audits")
	return
}

func (d *Database) Geo() (coll *Collection) {
	coll = d.getCollection("geo")
	return
}

func Connect() (err error) {
	mgoUrl, err := url.Parse(config.Config.MongoUri)
	if err != nil {
		err = &ConnectionError{
			errors.Wrap(err, "database: Failed to parse mongo uri"),
		}
		return
	}

	vals := mgoUrl.Query()
	mgoSsl := vals.Get("ssl")
	mgoSslCerts := vals.Get("ssl_ca_certs")
	vals.Del("ssl")
	vals.Del("ssl_ca_certs")
	mgoUrl.RawQuery = vals.Encode()
	mgoUri := mgoUrl.String()

	if mgoSsl == "true" {
		info, e := mgo.ParseURL(mgoUri)
		if e != nil {
			err = &ConnectionError{
				errors.Wrap(e, "database: Failed to parse mongo url"),
			}
			return
		}

		info.DialServer = func(addr *mgo.ServerAddr) (
			conn net.Conn, err error) {

			tlsConf := &tls.Config{}

			if mgoSslCerts != "" {
				caData, e := ioutil.ReadFile(mgoSslCerts)
				if e != nil {
					err = &CertificateError{
						errors.Wrap(e, "database: Failed to load certificate"),
					}
					return
				}

				caPool := x509.NewCertPool()
				if ok := caPool.AppendCertsFromPEM(caData); !ok {
					err = &CertificateError{
						errors.Wrap(err,
							"database: Failed to parse certificate"),
					}
					return
				}

				tlsConf.RootCAs = caPool
			}

			conn, err = tls.Dial("tcp", addr.String(), tlsConf)
			return
		}
		Session, err = mgo.DialWithInfo(info)
		if err != nil {
			err = &ConnectionError{
				errors.Wrap(err, "database: Connection error"),
			}
			return
		}
	} else {
		Session, err = mgo.Dial(mgoUri)
		if err != nil {
			err = &ConnectionError{
				errors.Wrap(err, "database: Connection error"),
			}
			return
		}
	}

	Session.SetMode(mgo.Strong, true)

	err = ValidateDatabase()
	if err != nil {
		Session = nil
		return
	}

	return
}

func ValidateDatabase() (err error) {
	db := GetDatabase()

	names, err := db.database.CollectionNames()
	if err != nil {
		err = ParseError(err)
		return
	}

	for _, name := range names {
		if name == "servers" {
			err = &errortypes.DatabaseError{
				errors.New("database: Cannot connect to pritunl database"),
			}
			return
		}
	}

	return
}

func GetDatabase() (db *Database) {
	sess := Session
	if sess == nil {
		return
	}

	session := sess.Copy()
	database := session.DB("")

	db = &Database{
		session:  session,
		database: database,
	}
	return
}

func addIndexes() (err error) {
	db := GetDatabase()
	defer db.Close()

	coll := db.Users()
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"username"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"type"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"roles"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Audits()
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"user"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Policies()
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"roles"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"services"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.CsrfTokens()
	err = coll.EnsureIndex(mgo.Index{
		Key:         []string{"timestamp"},
		ExpireAfter: 168 * time.Hour,
		Background:  true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
		return
	}

	coll = db.Nonces()
	err = coll.EnsureIndex(mgo.Index{
		Key:         []string{"timestamp"},
		ExpireAfter: 24 * time.Hour,
		Background:  true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.SshChallenges()
	err = coll.EnsureIndex(mgo.Index{
		Key:         []string{"timestamp"},
		ExpireAfter: 6 * time.Minute,
		Background:  true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.SshCertificates()
	err = coll.EnsureIndex(mgo.Index{
		Key:         []string{"timestamp"},
		ExpireAfter: 168 * time.Hour,
		Background:  true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.KeybaseChallenges()
	err = coll.EnsureIndex(mgo.Index{
		Key:         []string{"timestamp"},
		ExpireAfter: 6 * time.Minute,
		Background:  true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Sessions()
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"user"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Tasks()
	err = coll.EnsureIndex(mgo.Index{
		Key:         []string{"timestamp"},
		ExpireAfter: 720 * time.Hour,
		Background:  true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.Events()
	err = coll.EnsureIndex(mgo.Index{
		Key:        []string{"channel"},
		Background: true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
	}

	coll = db.AcmeChallenges()
	err = coll.EnsureIndex(mgo.Index{
		Key:         []string{"timestamp"},
		ExpireAfter: 3 * time.Minute,
		Background:  true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
		return
	}

	coll = db.Geo()
	err = coll.EnsureIndex(mgo.Index{
		Key:         []string{"t"},
		ExpireAfter: 360 * time.Hour,
		Background:  true,
	})
	if err != nil {
		err = &IndexError{
			errors.Wrap(err, "database: Index error"),
		}
		return
	}

	return
}

func addCollections() (err error) {
	db := GetDatabase()
	defer db.Close()
	coll := db.Events()

	names, err := db.database.CollectionNames()
	if err != nil {
		err = ParseError(err)
		return
	}

	for _, name := range names {
		if name == "events" {
			return
		}
	}

	err = coll.Create(&mgo.CollectionInfo{
		Capped:   true,
		MaxDocs:  1000,
		MaxBytes: 5242880,
	})
	if err != nil {
		err = ParseError(err)
		return
	}

	return
}

func init() {
	module := requires.New("database")
	module.After("config")

	module.Handler = func() (err error) {
		for {
			e := Connect()
			if e != nil {
				logrus.WithFields(logrus.Fields{
					"error": e,
				}).Error("database: Connection error")
			} else {
				break
			}

			time.Sleep(constants.RetryDelay)
		}

		err = addCollections()
		if err != nil {
			return
		}

		err = addIndexes()
		if err != nil {
			return
		}

		return
	}
}
