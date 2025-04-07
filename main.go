package main

import (
	"time"

	"github.com/pritunl/pritunl-zero/cmd"
)

func main() {
	defer time.Sleep(300 * time.Millisecond)
	cmd.Execute()
}
