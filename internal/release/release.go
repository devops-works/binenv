package release

import "context"

// // A Version holds a download path for a specific version
// type Version struct {
// 	Version string
// 	URL     string
// }

// Release contains release configuration
type Release struct {
	Type string `yaml:"type"`
	URL  string `yaml:"url"`
}

// Fetcher should implement fetching a release from a version
// and return a path where the release has been downloaded
type Fetcher interface {
	Fetch(context.Context, string) (string, error)
}

// Factory returns instances that comply to Installer interface
func (r Release) Factory() Fetcher {
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

	return nil
}
