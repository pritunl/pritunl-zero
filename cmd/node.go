package cmd

import (
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/router"
	"gopkg.in/mgo.v2/bson"
)

func Node() (err error) {
	nde := &node.Node{
		Id:   bson.ObjectIdHex(config.Config.NodeId),
		Type: node.Management,
	}
	err = nde.Init()
	if err != nil {
		return
	}

	routr := &router.Router{}

	routr.Init()

	err = routr.Run()
	if err != nil {
		return
	}

	return
}
