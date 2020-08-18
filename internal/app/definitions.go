package app

import (
	"gitlab.com/devopsworks/tools/binenv/internal/fetch"
	"gitlab.com/devopsworks/tools/binenv/internal/install"
	"gitlab.com/devopsworks/tools/binenv/internal/list"
	"gitlab.com/devopsworks/tools/binenv/internal/mapping"
)

// Distributions holds the liste of available software sources
type Distributions struct {
	Sources map[string]Sources `yaml:"sources"`
}

// Sources contains a software source definition
type Sources struct {
	// Name    string  `yaml:"name"`
	Map     mapping.Remapper `yaml:"map"`
	List    list.List        `yaml:"list"`
	Fetch   fetch.Fetch      `yaml:"fetch"`
	Install install.Install  `yaml:"install"`
}
