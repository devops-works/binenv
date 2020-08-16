package app

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

func getConfigDir() (string, error) {
	var err error

	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		dir, err = homedir.Dir()
		if err != nil {
			return "", err
		}
		dir += "/.config/binenv"
	}

	return dir, nil
}

func getPackages() {
	
}
