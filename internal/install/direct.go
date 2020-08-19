package install

// Direct installs directly downloaded binaries
type Direct struct {
	filter string
}

// Install will move the binary from src to dst
func (d Direct) Install(src, dst, version string) error {
	return installFile(src, dst)
}
