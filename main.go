package main

import (
	"os"

	"github.com/devops-works/binenv/cmd"
)

func main() {
	err := cmd.RootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
