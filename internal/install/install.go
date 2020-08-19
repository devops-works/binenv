package install

import (
	"io"
	"os"

	"github.com/devops-works/binenv/internal/mapping"
)

// Install defines the install config struct
type Install struct {
	Type     string   `yaml:"type"`
	Binaries []string `yaml:"binaries"`
}

// Installer should implement installation
type Installer interface {
	Install(src, dst, version string, mapper mapping.Mapper) error
}

// Factory returns instances that comply to Installer interface
func (i Install) Factory(filters []string) Installer {
	switch i.Type {
	case "direct":
		return Direct{}
	case "zip":
		return Zip{filters: filters}
	case "tgz":
		return Tgz{filters: filters}
	}
	return nil
}

func installFile(src, dst string) error {
	fd, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fd.Close()

	fs, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(fd, fs)
	if err != nil {
		return err
	}

	err = os.Chmod(dst, 0700)
	return err

}
