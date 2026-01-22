package secret

import (
	"strings"

	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/errortypes"
	"github.com/pritunl/pritunl-zero/utils"
)

type Secret struct {
	Id         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string        `bson:"name" json:"name"`
	Comment    string        `bson:"comment" json:"comment"`
	Type       string        `bson:"type" json:"type"`
	Key        string        `bson:"key" json:"key"`
	Value      string        `bson:"value" json:"value"`
	Region     string        `bson:"region" json:"region"`
	PublicKey  string        `bson:"public_key" json:"public_key"`
	PrivateKey string        `bson:"private_key" json:"-"`
}

func (c *Secret) Validate(db *database.Database) (
	errData *errortypes.ErrorData, err error) {

	c.Name = utils.FilterName(c.Name)

	switch c.Type {
	case AWS, "":
		c.Type = AWS

		if c.Region == "" {
			c.Region = "us-east-1"
		}

		break
	case Cloudflare:
		c.Value = ""
		c.Region = ""

		break
	case OracleCloud:
		break
	case GoogleCloud:
		c.Value = ""
		c.Region = ""

		break
	default:
		errData = &errortypes.ErrorData{
			Error:   "invalid_secret_type",
			Message: "Secret type invalid",
		}
		return
	}

	if c.PrivateKey == "" {
		privKey, pubKey, e := utils.GenerateRsaKey()
		if e != nil {
			err = e
			return
		}

		c.PublicKey = strings.TrimSpace(string(pubKey))
		c.PrivateKey = strings.TrimSpace(string(privKey))
	}

	return
}

func (c *Secret) GetOracleProvider() (prov *OracleProvider, err error) {
	prov, err = NewOracleProvider(c)
	if err != nil {
		return
	}

	return
}

func (c *Secret) Commit(db *database.Database) (err error) {
	coll := db.Secrets()

	err = coll.Commit(c.Id, c)
	if err != nil {
		return
	}

	return
}

func (c *Secret) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Secrets()

	err = coll.CommitFields(c.Id, c, fields)
	if err != nil {
		return
	}

	return
}

func (c *Secret) Insert(db *database.Database) (err error) {
	coll := db.Secrets()

	if !c.Id.IsZero() {
		err = &errortypes.DatabaseError{
			errors.New("secret: Secret already exists"),
		}
		return
	}

	_, err = coll.InsertOne(db, c)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
