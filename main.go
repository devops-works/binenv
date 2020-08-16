package main

import (
	"fmt"
	"os"

	"gitlab.com/devopsworks/tools/binenv/cmd"
)

func main() {
	fmt.Println(os.Args[0])
	root := cmd.RootCmd()
	root.Execute()
}
