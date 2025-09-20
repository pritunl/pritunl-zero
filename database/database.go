package database

import (
	"context"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/pritunl/mongo-go-driver/v2/mongo/readconcern"
	"github.com/pritunl/mongo-go-driver/v2/mongo/writeconcern"
	"github.com/pritunl/mongo-go-driver/v2/x/mongo/driver/connstring"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/requires"
	"github.com/sirupsen/logrus"
)

type Database struct {
	ctx      context.Context
	client   *mongo.Client
	database *mongo.Database
}

func (d *Database) Deadline() (time.Time, bool) {
	if d.ctx != nil {
		return d.ctx.Deadline()
	}
	return time.Time{}, false
}

func (d *Database) Done() <-chan struct{} {
	if d.ctx != nil {
		return d.ctx.Done()
	}
	return nil
}

func (d *Database) Err() error {
	if d.ctx != nil {
		return d.ctx.Err()
	}
	return nil
}

func (d *Database) Value(key interface{}) interface{} {
	if d.ctx != nil {
		return d.ctx.Value(key)
	}
	return nil
}

func (d *Database) String() string {
	return "context.database"
}

func (d *Database) Close() {
}

func (d *Database) getCollection(name string) (coll *Collection) {
	coll = &Collection{
		db:         d,
		Collection: d.database.Collection(name),
	}
	return
}

func (d *Database) getCollectionWeak(name string) (coll *Collection) {
	opts := options.Collection()
	opts.SetWriteConcern(writeconcern.W1())
	opts.SetReadConcern(readconcern.Local())

	coll = &Collection{
		db:         d,
		Collection: d.database.Collection(name, opts),
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

func (d *Database) Devices() (coll *Collection) {
	coll = d.getCollection("devices")
	return
}

func (d *Database) Alerts() (coll *Collection) {
	coll = d.getCollection("alerts")
	return
}

func (d *Database) AlertsEvent() (coll *Collection) {
	coll = d.getCollection("alerts_event")
	return
}

func (d *Database) AlertsEventLock() (coll *Collection) {
	coll = d.getCollection("alerts_event_lock")
	return
}

func (d *Database) Checks() (coll *Collection) {
	coll = d.getCollection("checks")
	return
}

func (d *Database) Endpoints() (coll *Collection) {
	coll = d.getCollection("endpoints")
	return
}

func (d *Database) EndpointsSystem() (coll *Collection) {
	coll = d.getCollection("endpoints_system")
	return
}

func (d *Database) EndpointsLoad() (coll *Collection) {
	coll = d.getCollection("endpoints_load")
	return
}

func (d *Database) EndpointsDisk() (coll *Collection) {
	coll = d.getCollection("endpoints_disk")
	return
}

func (d *Database) EndpointsDiskIo() (coll *Collection) {
	coll = d.getCollection("endpoints_diskio")
	return
}

func (d *Database) EndpointsNetwork() (coll *Collection) {
	coll = d.getCollection("endpoints_network")
	return
}

func (d *Database) EndpointsCheck() (coll *Collection) {
	coll = d.getCollection("endpoints_check")
	return
}

func (d *Database) EndpointsCheckLog() (coll *Collection) {
	coll = d.getCollection("endpoints_check_log")
	return
}

func (d *Database) EndpointsKmsg() (coll *Collection) {
	coll = d.getCollection("endpoints_kmsg")
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

func (d *Database) Versions() (coll *Collection) {
	coll = d.getCollection("versions")
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

func (d *Database) SecondaryTokens() (coll *Collection) {
	coll = d.getCollection("secondary_tokens")
	return
}

func (d *Database) Nonces() (coll *Collection) {
	coll = d.getCollection("nonces")
	return
}

func (d *Database) Rokeys() (coll *Collection) {
	coll = d.getCollection("rokeys")
	return
}

func (d *Database) Settings() (coll *Collection) {
	coll = d.getCollection("settings")
	return
}

func (d *Database) Events() (coll *Collection) {
	coll = d.getCollectionWeak("events")
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

func (d *Database) Secrets() (coll *Collection) {
	coll = d.getCollection("secrets")
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
	mongoUrl, err := connstring.ParseAndValidate(config.Config.MongoUri)
	if err != nil {
		err = &ConnectionError{
			errors.Wrap(err, "database: Failed to parse mongo uri"),
		}
		return
	}

	logrus.WithFields(logrus.Fields{
		"mongodb_hosts": mongoUrl.Hosts,
	}).Info("database: Connecting to MongoDB server")

	if mongoUrl.Database != "" {
		DefaultDatabase = mongoUrl.Database
	}

	opts := options.Client().ApplyURI(config.Config.MongoUri)
	opts.SetRetryReads(true)
	opts.SetRetryWrites(true)

	client, err := mongo.Connect(opts)
	if err != nil {
		err = &ConnectionError{
			errors.Wrap(err, "database: Connection error"),
		}
		return
	}

	setClient(client)

	version, err := ValidateDatabase()
	if err != nil {
		return
	}

	logrus.WithFields(logrus.Fields{
		"mongodb_hosts":   mongoUrl.Hosts,
		"mongodb_version": version,
	}).Info("database: Connected to MongoDB server")

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

func ValidateDatabase() (version string, err error) {
	db := GetDatabase()
	defer db.Close()

	buildInfo := bson.M{}
	err = db.database.RunCommand(
		db,
		bson.D{{"buildInfo", 1}},
	).Decode(&buildInfo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("database: Failed to get MongoDB version")
		err = nil
	}

	version, ok := buildInfo["version"].(string)
	if version == "" || !ok {
		version = "unknown"
	}

	cursor, err := db.database.ListCollections(
		db, &bson.M{})
	if err != nil {
		err = ParseError(err)
		return
	}
	defer cursor.Close(db)

	for cursor.Next(db) {
		item := &struct {
			Name string `bson:"name"`
		}{}
		err = cursor.Decode(item)
		if err != nil {
			err = ParseError(err)
			return
		}

		if item.Name == "servers" {
			err = &errortypes.DatabaseError{
				errors.New("database: Cannot connect to pritunl database"),
			}
			return
		}
	}

	err = cursor.Err()
	if err != nil {
		err = ParseError(err)
		return
	}

	return
}

func getDatabase(ctx context.Context, client *mongo.Client) *Database {
	if client == nil {
		return nil
	}

	database := client.Database(DefaultDatabase)

	return &Database{
		ctx:      ctx,
		client:   client,
		database: database,
	}
}

func GetDatabase() *Database {
	return getDatabase(nil, getClient())
}

func GetDatabaseCtx(ctx context.Context) *Database {
	return getDatabase(ctx, getClient())
}

func addIndexes() (err error) {
	db := GetDatabase()
	defer db.Close()

	index := &Index{
		Collection: db.Users(),
		Keys: &bson.D{
			{"username", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}
	index = &Index{
		Collection: db.Users(),
		Keys: &bson.D{
			{"type", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}
	index = &Index{
		Collection: db.Users(),
		Keys: &bson.D{
			{"roles", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}
	index = &Index{
		Collection: db.Users(),
		Keys: &bson.D{
			{"token", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Logs(),
		Keys: &bson.D{
			{"timestamp", 1},
		},
		Expire: 4320 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Audits(),
		Keys: &bson.D{
			{"user", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Policies(),
		Keys: &bson.D{
			{"roles", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}
	index = &Index{
		Collection: db.Policies(),
		Keys: &bson.D{
			{"services", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}
	index = &Index{
		Collection: db.Policies(),
		Keys: &bson.D{
			{"authorities", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.CsrfTokens(),
		Keys: &bson.D{
			{"timestamp", 1},
		},
		Expire: 168 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.SecondaryTokens(),
		Keys: &bson.D{
			{"timestamp", 1},
		},
		Expire: 3 * time.Minute,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Nodes(),
		Keys: &bson.D{
			{"name", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Nonces(),
		Keys: &bson.D{
			{"timestamp", 1},
		},
		Expire: 24 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Rokeys(),
		Keys: &bson.D{
			{"type", 1},
			{"timeblock", 1},
		},
		Unique: true,
	}
	err = index.Create()
	if err != nil {
		return
	}
	index = &Index{
		Collection: db.Rokeys(),
		Keys: &bson.D{
			{"timestamp", 1},
		},
		Expire: 720 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Devices(),
		Keys: &bson.D{
			{"user", 1},
			{"mode", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}
	index = &Index{
		Collection: db.Devices(),
		Keys: &bson.D{
			{"provider", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Alerts(),
		Keys: &bson.D{
			{"name", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Alerts(),
		Keys: &bson.D{
			{"roles", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.AlertsEvent(),
		Keys: &bson.D{
			{"timestamp", 1},
		},
		Expire: 48 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.AlertsEvent(),
		Keys: &bson.D{
			{"source", 1},
			{"resource", 1},
			{"timestamp", -1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.AlertsEventLock(),
		Keys: &bson.D{
			{"timestamp", 1},
		},
		Expire: 72 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Checks(),
		Keys: &bson.D{
			{"name", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Checks(),
		Keys: &bson.D{
			{"roles", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Endpoints(),
		Keys: &bson.D{
			{"roles", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Authorities(),
		Keys: &bson.D{
			{"host_tokens", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}
	index = &Index{
		Collection: db.Authorities(),
		Keys: &bson.D{
			{"hsm_token", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.SshChallenges(),
		Keys: &bson.D{
			{"timestamp", 1},
		},
		Expire: 6 * time.Minute,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.SshCertificates(),
		Keys: &bson.D{
			{"timestamp", 1},
		},
		Expire: 168 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Devices(),
		Keys: &bson.D{
			{"user", 1},
			{"mode", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}
	index = &Index{
		Collection: db.Devices(),
		Keys: &bson.D{
			{"provider", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Sessions(),
		Keys: &bson.D{
			{"user", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Sessions(),
		Keys: &bson.D{
			{"timestamp", 1},
		},
		Expire: 4320 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Tasks(),
		Keys: &bson.D{
			{"timestamp", 1},
		},
		Expire: 48 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Events(),
		Keys: &bson.D{
			{"channel", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.AcmeChallenges(),
		Keys: &bson.D{
			{"timestamp", 1},
		},
		Expire: 3 * time.Minute,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.Geo(),
		Keys: &bson.D{
			{"t", 1},
		},
		Expire: 360 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsSystem(),
		Keys: &bson.D{
			{"t", 1},
			{"e", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsSystem(),
		Keys: &bson.D{
			{"t", 1},
		},
		Expire: 2160 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsLoad(),
		Keys: &bson.D{
			{"t", 1},
			{"e", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsLoad(),
		Keys: &bson.D{
			{"t", 1},
		},
		Expire: 2160 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsDisk(),
		Keys: &bson.D{
			{"t", 1},
			{"e", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsDisk(),
		Keys: &bson.D{
			{"t", 1},
		},
		Expire: 2160 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsDiskIo(),
		Keys: &bson.D{
			{"t", 1},
			{"e", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsDiskIo(),
		Keys: &bson.D{
			{"t", 1},
		},
		Expire: 2160 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsNetwork(),
		Keys: &bson.D{
			{"t", 1},
			{"e", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsNetwork(),
		Keys: &bson.D{
			{"t", 1},
		},
		Expire: 2160 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsCheck(),
		Keys: &bson.D{
			{"t", 1},
			{"c", 1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsCheck(),
		Keys: &bson.D{
			{"t", 1},
		},
		Expire: 4320 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsCheckLog(),
		Keys: &bson.D{
			{"t", -1},
			{"c", -1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsCheckLog(),
		Keys: &bson.D{
			{"t", 1},
		},
		Expire: 4320 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsKmsg(),
		Keys: &bson.D{
			{"e", -1},
			{"b", -1},
			{"s", -1},
		},
	}
	err = index.Create()
	if err != nil {
		return
	}

	index = &Index{
		Collection: db.EndpointsKmsg(),
		Keys: &bson.D{
			{"t", 1},
		},
		Expire: 4320 * time.Hour,
	}
	err = index.Create()
	if err != nil {
		return
	}

	return
}

func addCollections() (err error) {
	db := GetDatabase()
	defer db.Close()

	cursor, err := db.database.ListCollections(
		db, &bson.M{})
	if err != nil {
		err = ParseError(err)
		return
	}
	defer cursor.Close(db)

	eventsExists := false
	isCapped := false

	for cursor.Next(db) {
		item := &struct {
			Name    string `bson:"name"`
			Options bson.M `bson:"options"`
		}{}
		err = cursor.Decode(item)
		if err != nil {
			err = ParseError(err)
			return
		}

		if item.Name == "events" {
			eventsExists = true
			if options, ok := item.Options["capped"]; ok {
				if cappedBool, ok := options.(bool); ok && cappedBool {
					isCapped = true
				}
			}
			break
		}
	}

	err = cursor.Err()
	if err != nil {
		err = ParseError(err)
		return
	}

	if eventsExists && !isCapped {
		logrus.WithFields(logrus.Fields{
			"collection": "events",
		}).Warning("database: Correcting events capped collection")

		err = db.database.Collection("events").Drop(context.Background())
		if err != nil {
			err = ParseError(err)
			return
		}
		eventsExists = false
	}

	if !eventsExists {
		err = db.database.RunCommand(
			context.Background(),
			bson.D{
				{"create", "events"},
				{"capped", true},
				{"max", 1000},
				{"size", 5242880},
			},
		).Err()
		if err != nil {
			err = ParseError(err)
			return
		}
	}

	return
}

func fixData() (err error) {
	db := GetDatabase()
	defer db.Close()

	coll := db.Policies()
	_, err = coll.UpdateMany(db, &bson.M{
		"admin_secondary": nil,
	}, &bson.M{
		"$unset": &bson.M{
			"admin_secondary": 1,
		},
	})
	if err != nil {
		err = ParseError(err)
		return
	}
	_, err = coll.UpdateMany(db, &bson.M{
		"user_secondary": nil,
	}, &bson.M{
		"$unset": &bson.M{
			"user_secondary": 1,
		},
	})
	if err != nil {
		err = ParseError(err)
		return
	}
	_, err = coll.UpdateMany(db, &bson.M{
		"proxy_secondary": nil,
	}, &bson.M{
		"$unset": &bson.M{
			"proxy_secondary": 1,
		},
	})
	if err != nil {
		err = ParseError(err)
		return
	}
	_, err = coll.UpdateMany(db, &bson.M{
		"authority_secondary": nil,
	}, &bson.M{
		"$unset": &bson.M{
			"authority_secondary": 1,
		},
	})
	if err != nil {
		err = ParseError(err)
		return
	}

	coll = db.SshCertificates()
	_, err = coll.UpdateMany(db, &bson.M{
		"user_id": nil,
	}, &bson.M{
		"$unset": &bson.M{
			"user_id": 1,
		},
	})
	if err != nil {
		err = ParseError(err)
		return
	}

	coll = db.SshChallenges()
	_, err = coll.UpdateMany(db, &bson.M{
		"certificate_id": nil,
	}, &bson.M{
		"$unset": &bson.M{
			"certificate_id": 1,
		},
	})
	if err != nil {
		err = ParseError(err)
		return
	}

	return
}

func init() {
	module := requires.New("database")
	module.After("logger")

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

		err = fixData()
		if err != nil {
			return
		}

		return
	}
}
