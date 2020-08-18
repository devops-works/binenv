package install

import (
	"bytes"
	"fmt"
	"html/template"
	"runtime"
	"sync"

	gov "github.com/hashicorp/go-version"
)

// Install defines the install config struct
type Install struct {
	Type     string   `yaml:"type"`
	Binaries []string `yaml:"binaries"`
}

// Installer should implement installation
type Installer interface {
	Install(src, dst, version string) error
}

type tplArg struct {
	OS           string
	Arch         string
	Version      string
	NakedVersion string
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

func matchFilters(file string, filters []string, version string) (bool, error) {
	arg := tplArg{
		Arch:         runtime.GOARCH,
		OS:           runtime.GOOS,
		Version:      version,
		NakedVersion: gov.Must(gov.NewVersion(version)).String(),
	}

	var once sync.Once

	tpls := []*template.Template{}

	var onceErr error

	onceBody := func() {
		for _, v := range filters {
			tpl, err := template.New("install").Parse(v)
			if err != nil {
				onceErr = err
				return
			}
			tpls = append(tpls, tpl)
		}
	}

	once.Do(onceBody)
	if onceErr != nil {
		return false, onceErr
	}

	for _, t := range tpls {
		buf := bytes.Buffer{}
		err := t.Execute(&buf, arg)
		if err != nil {
			return false, err
		}
		if buf.String() == file {
			fmt.Printf("file %s matches filters\n", file)
			return true, nil
		}
	}

	return false, nil
}
