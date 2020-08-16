package install

import "os"

// Direct installs directly downloaded binaries
type Direct struct{}

// Install will move the binary from src to dst
func (d Direct) Install(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return err
	}

	err = os.Chmod(dst, 0700)
	return err
}
