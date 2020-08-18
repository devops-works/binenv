package fetch

import "context"

// Fetcher should implement fetching a release from a version
// and return a path where the release has been downloaded
type Fetcher interface {
	Fetch(context.Context, string) (string, error)
}

// Fetch contains fetch configuration
type Fetch struct {
	Type string `yaml:"type"`
	URL  string `yaml:"url"`
}

// Factory returns instances that comply to Fecther interface
func (r Fetch) Factory() Fetcher {
	switch r.Type {
	case "download":
		return Download{
			url: r.URL,
		}
	default:
		return Download{
			url: r.URL,
		}
	}
}
