package app

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

func getCacheDir() (string, error) {
	var err error

	dir := os.Getenv("XDG_CACHE_HOME")
	if dir == "" {
		dir, err = homedir.Dir()
		if err != nil {
			return "", err
		}
		dir += "/.cache/binenv"
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

func stringToEnvVarName(st string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	if err != nil {
		log.Fatal(err)
	}
	result := reg.ReplaceAllString(st, "_")
	return strings.ToUpper(result)
}
