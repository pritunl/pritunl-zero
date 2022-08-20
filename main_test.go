package main

import (
	"github.com/pritunl/pritunl-zero/cmd"
	"github.com/pritunl/pritunl-zero/constants"
	"testing"
)

func TestServer(t *testing.T) {
	constants.Production = false

	Init()
	err := cmd.Node(true)
	if err != nil {
		panic(err)
	}

	return
}
