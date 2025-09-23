package utils

import (
	"fmt"
	"math/rand"
)

var (
	randElm = []string{
		"copper",
		"argon",
		"xenon",
		"radon",
		"cobalt",
		"nickel",
		"carbon",
		"helium",
		"nitrogen",
		"radium",
		"lithium",
		"silicon",
	}
)

func RandName() (name string) {
	name = fmt.Sprintf("%s-%d", randElm[rand.Intn(len(randElm))],
		rand.Intn(8999)+1000)
	return
}
