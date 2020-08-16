package app

import (
	"gitlab.com/devopsworks/tools/binenv/internal/install"
	"gitlab.com/devopsworks/tools/binenv/internal/list"
	"gitlab.com/devopsworks/tools/binenv/internal/release"
)

// Definitions holds the liste of available software sources
type Definitions struct {
	Sources map[string]Sources `yaml:"sources"`
}

// Sources contains a software source definition
type Sources struct {
	// Name    string  `yaml:"name"`
	List    list.List       `yaml:"list"`
	Release release.Release `yaml:"release"`
	Install install.Install `yaml:"install"`
}
