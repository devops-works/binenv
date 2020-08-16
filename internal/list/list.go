package list

import (
	"context"
	"sync"
)

// Lister should return a list of available release versions
type Lister interface {
	Get(context.Context, *sync.WaitGroup) ([]string, error)
}

// List contains list definition
type List struct {
	Type string `yaml:"type"`
	URL  string `yaml:"url"`
}

// Factory returns instances that comply to Lister interface
func (l List) Factory() Lister {
	switch l.Type {
	case "github":
		return Github{
			url: l.URL,
		}
	}

	return nil
}
