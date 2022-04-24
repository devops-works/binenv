package app

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// GetDefaultBinDir returns the bin directory in usermode
func GetDefaultBinDir() string {
	d, err := homedir.Dir()
	if err != nil {
		d = "~"
	}
	d += "/.binenv/"

	return d
}

// GetDefaultLinkDir returns the bin directory in usermode
func GetDefaultLinkDir() string {
	return GetDefaultBinDir()
}

// GetDefaultConfDir returns the config directory in usermode
func GetDefaultConfDir() string {
	var err error

	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		dir, err = homedir.Dir()
		if err != nil {
			dir = "~"
		}
		dir += "/.config/binenv"
	}

	return dir
}

// GetDefaultCacheDir returns the cache directory in usermode
func GetDefaultCacheDir() string {
	var err error

	dir := os.Getenv("XDG_CACHE_HOME")
	if dir == "" {
		dir, err = homedir.Dir()
		if err != nil {
			dir = "~"
		}
		dir += "/.cache/binenv"
	}

	return dir
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
