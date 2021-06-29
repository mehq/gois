package main

import (
	"os"

	"github.com/mzbaulhaque/gois/cmd"
)

func main() {
	if cmd.Execute() != nil {
		os.Exit(1)
	}
}
