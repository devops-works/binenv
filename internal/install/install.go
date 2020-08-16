package install

// Install defines the install config struct
type Install struct {
	Type     string   `yaml:"type"`
	Binaries []string `yaml:"binaries"`
}

// Installer should implement installation
type Installer interface {
	Install(src string, dst string, patterns []string) error
}

// Factory returns instances that comply to Installer interface
func (i Install) Factory() Installer {
	switch i.Type {
	case "direct":
		return Direct{}
	}

	return nil
}
