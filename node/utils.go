package node

import (
	"github.com/pritunl/pritunl-zero/database"
	"gopkg.in/mgo.v2/bson"
)

func Get(db *database.Database, nodeId bson.ObjectId) (
	nde *Node, err error) {

	coll := db.Nodes()
	nde = &Node{}

	err = coll.FindOneId(nodeId, nde)
	if err != nil {
		return
	}

	return
}

func GetAll(db *database.Database) (nodes []*Node, err error) {

	coll := db.Nodes()
	nodes = []*Node{}

	cursor := coll.Find(bson.M{}).Iter()

	nde := &Node{}
	for cursor.Next(nde) {
		nde.SetActive()
		nodes = append(nodes, nde)
		nde = &Node{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Remove(db *database.Database, nodeId bson.ObjectId) (err error) {
	coll := db.Nodes()

	_, err = coll.RemoveAll(&bson.M{
		"_id": nodeId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
