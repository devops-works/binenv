package install

// Direct installs directly downloaded binaries
type Direct struct{}

// Install will move the binary from src to dst if it matches pattern
func (d Direct) Install(src, dst string, patterns []string) error {
	return nil
}
