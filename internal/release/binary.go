package release

import "context"

// Binary handles direct binary releases
type Binary struct {
	url string
}

// Fetch gets the package
func (b Binary) Fetch(ctx context.Context, v string) (string, error) {
	return "", nil
}
