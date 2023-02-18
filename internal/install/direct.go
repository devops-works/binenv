package install

import "github.com/devops-works/binenv/internal/mapping"

// Direct installs directly downloaded binaries
type Direct struct {
}

// Install will move the binary from src to dst
func (d Direct) Install(src, dst, version string, mapper mapping.Mapper) error {
	return installFile(src, dst)
}
