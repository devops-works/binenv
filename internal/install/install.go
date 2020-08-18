package install

// Install defines the install config struct
type Install struct {
	Type     string   `yaml:"type"`
	Binaries []string `yaml:"binaries"`
}

// Installer should implement installation
type Installer interface {
	Install(src, dst, version string) error
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
