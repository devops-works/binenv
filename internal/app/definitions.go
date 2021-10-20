package app

import (
	"github.com/devops-works/binenv/internal/fetch"
	"github.com/devops-works/binenv/internal/install"
	"github.com/devops-works/binenv/internal/list"
	"github.com/devops-works/binenv/internal/mapping"
	"github.com/devops-works/binenv/internal/platform"
)

// Distributions holds the liste of available software sources
type Distributions struct {
	Sources map[string]Sources `yaml:"sources"`
}

// Sources contains a software source definition
type Sources struct {
	// Name    string  `yaml:"name"`
	Description        string              `yaml:"description"`
	URL                string              `yaml:"url"`
	Map                mapping.Remapper    `yaml:"map"`
	List               list.List           `yaml:"list"`
	Fetch              fetch.Fetch         `yaml:"fetch"`
	Install            install.Install     `yaml:"install"`
	PostInstallMessage string              `yaml:"post_install_message"`
	SupportedPlatforms []platform.Platform `yaml:"supported_platforms"`
}
