package main

import (
	"gitlab.com/devopsworks/tools/binenv/cmd"
)

func main() {
	root := cmd.RootCmd()
	root.Execute()
}
