package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/config"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/router"
	"gopkg.in/mgo.v2/bson"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Node() (err error) {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	nde := &node.Node{
		Id: bson.ObjectIdHex(config.Config.NodeId),
	}
	err = nde.Init()
	if err != nil {
		return
	}

	routr := &router.Router{}

	routr.Init()

	go func() {
		err = routr.Run()
		if err != nil {
			panic(err)
		}
	}()

	<-sig
	logrus.Info("cmd.node: Shutting down")
	go routr.Shutdown()
	if constants.Production {
		time.Sleep(200 * time.Millisecond)
	} else {
		time.Sleep(300 * time.Millisecond)
	}

	return
}
