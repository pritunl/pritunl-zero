package node

import (
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/container/set"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/event"
	"github.com/pritunl/pritunl-zero/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var Self *Node

type Node struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Type      string        `bson:"type" json:"type"`
	Timestamp time.Time     `bson:"timestamp" json:"timestamp"`
	Port      int           `bson:"port" json:"port"`
	Protocol  string        `bson:"protocol" json:"protocol"`
	Memory    float64       `bson:"memory" json:"memory"`
	Load1     float64       `bson:"load1" json:"load1"`
	Load5     float64       `bson:"load5" json:"load5"`
	Load15    float64       `bson:"load15" json:"load15"`
}

func (n *Node) Commit(db *database.Database) (err error) {
	coll := db.Nodes()

	err = coll.Commit(n.Id, n)
	if err != nil {
		return
	}

	return
}

func (n *Node) CommitFields(db *database.Database, fields set.Set) (
	err error) {

	coll := db.Nodes()

	err = coll.CommitFields(n.Id, n, fields)
	if err != nil {
		return
	}

	return
}

func (n *Node) update(db *database.Database) (err error) {
	coll := db.Nodes()

	change := mgo.Change{
		Update: &bson.M{
			"$set": &bson.M{
				"timestamp": n.Timestamp,
				"memory":    n.Memory,
				"load1":     n.Load1,
				"load5":     n.Load5,
				"load15":    n.Load15,
			},
		},
		Upsert:    false,
		ReturnNew: true,
	}

	_, err = coll.Find(&bson.M{
		"_id": n.Id,
	}).Apply(change, n)
	if err != nil {
		return
	}

	return
}

func (n *Node) keepalive() {
	db := database.GetDatabase()
	defer db.Close()

	for {
		n.Timestamp = time.Now()

		mem, err := utils.MemoryUsed()
		if err != nil {
			n.Memory = 0

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("node: Failed to get memory")
		} else {
			n.Memory = mem
		}

		load, err := utils.LoadAverage()
		if err != nil {
			n.Load1 = 0
			n.Load5 = 0
			n.Load15 = 0

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("node: Failed to get load")
		} else {
			n.Load1 = load.Load1
			n.Load5 = load.Load5
			n.Load15 = load.Load15
		}

		err = n.update(db)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("node: Failed to update node")
		}

		time.Sleep(1 * time.Second)
	}
}

func (n *Node) Init() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	coll := db.Nodes()

	if n.Name == "" {
		n.Name = utils.RandName()
	}

	change := mgo.Change{
		Update: &bson.M{
			"$set": &bson.M{
				"_id":       n.Id,
				"type":      n.Type,
				"name":      n.Name,
				"timestamp": time.Now(),
			},
		},
		Upsert:    true,
		ReturnNew: true,
	}

	_, err = coll.Find(&bson.M{
		"_id": n.Id,
	}).Apply(change, n)
	if err != nil {
		return
	}

	event.PublishDispatch(db, "node.change")

	Self = n

	go n.keepalive()

	return
}
