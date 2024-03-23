package fetch

import (
	"context"
	"fmt"
	"os"

	"github.com/devops-works/binenv/internal/mapping"
)

// Fetcher should implement fetching a release from a version
// and return a path where the release has been downloaded
type Fetcher interface {
	Fetch(ctx context.Context, dist, version string, mapper mapping.Mapper) (string, error)
}

// Fetch contains fetch configuration
type Fetch struct {
	Type     string   `yaml:"type"`
	URL      string   `yaml:"url"`
	URLs     []string `yaml:"urls"`
	TokenEnv string   `yaml:"token_env"`
}

// Factory returns instances that comply to Fetcher interface
func (r Fetch) Factory() (Fetcher, error) {
	switch r.Type {
	// case "download":
	// 	return Download{
	// 		url: r.URL,
	// 	}
	default:
		headers := map[string]string{}
		if r.TokenEnv != "" {
			token := os.Getenv(r.TokenEnv)
			if token == "" {
				return nil, fmt.Errorf("token env var %s is not defined; did you export it ?", r.TokenEnv)
			}
			headers["PRIVATE-TOKEN"] = token
		}
		return Download{
			urls:    []string{r.URL},
			headers: headers,
		}, nil
	case "download_list":
		return Download{
			urls: r.URLs,
		}, nil
	}
}
