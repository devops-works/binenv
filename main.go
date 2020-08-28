package main

import (
	"github.com/devops-works/binenv/cmd"
)

func main() {
	// var verbose bool

	root := cmd.RootCmd()
	// root.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose operation")
	// verbose, _ := root.PersistentFlags().GetBool("verbose")
	// fmt.Printf("verbose is %t\n", verbose)

	root.Execute()
}
