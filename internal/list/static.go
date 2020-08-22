package list

import (
	"context"
)

// Static contains what is required to get a list of release from Github
type Static struct {
	versions []string
}

// Get returns a list of available versions
func (s Static) Get(ctx context.Context) ([]string, error) {
	return s.versions, nil
}
