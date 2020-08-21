package app

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// GetDefaultBinDir returns the bin directory
func GetDefaultBinDir() string {
	d, err := homedir.Dir()
	if err != nil {
		d = "~"
	}
	d += "/.binenv/"

	return d
}

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

func getDistributionsFilePath() (string, error) {
	conf, err := getConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(conf, "/distributions.yaml"), nil
}

func stringInSlice(st string, sl []string) bool {
	for _, v := range sl {
		if v == st {
			return true
		}
	}

	return false
}
