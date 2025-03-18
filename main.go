package main

import (
	"github.com/coeeter/cmdhelper/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
