package main

import (
	"github.com/devops-works/binenv/cmd"
)

var ()

func main() {
	root := cmd.RootCmd()
	root.Execute()
}
