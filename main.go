package main

import (
	"github.com/devops-works/binenv/cmd"
)

func main() {
	root := cmd.RootCmd()
	root.Execute()
}
