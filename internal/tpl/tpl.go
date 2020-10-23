package tpl

import (
	"bytes"
	"runtime"
	"strings"
	"sync"
	"text/template"

	"github.com/devops-works/binenv/internal/mapping"
	gov "github.com/hashicorp/go-version"
)

// Args holds templating args
type Args struct {
	OS           string
	Arch         string
	Version      string
	VersionMajor string
	VersionMinor string
	VersionPatch string
	NakedVersion string
	ExeExtension string
}

// New returns populated template Args
func New(v string, mapper mapping.Mapper) Args {
	rarch := runtime.GOARCH
	ros := runtime.GOOS

	if mapper != nil {
		rarch = mapper.MustInterpolate(runtime.GOARCH)
		ros = mapper.MustInterpolate(runtime.GOOS)
		// fmt.Printf("remapping arch %s to %s\n", runtime.GOARCH, rarch)
		// fmt.Printf("remapping os %s to %s\n", runtime.GOOS, ros)
	}
	a := Args{
		Arch:         rarch,
		OS:           ros,
		Version:      v,
		VersionMajor: strings.Split(v, ".")[0],
		VersionMinor: strings.Split(v, ".")[1],
		VersionPatch: strings.Split(v, ".")[2],
		NakedVersion: gov.Must(gov.NewVersion(v)).String(),
	}

	if a.OS == "windows" {
		a.ExeExtension = ".exe"
	}

	return a
}

// Interpolate changes some values by others
func (a Args) Interpolate(m map[string]string) {
	if v, ok := m[a.Arch]; ok {
		a.Arch = v
	}
	if v, ok := m[a.OS]; ok {
		a.OS = v
	}
}

// MatchFilters matches a file against a list of template filters
func (a Args) MatchFilters(file string, filters []string) (bool, error) {
	var once sync.Once

	tpls := []*template.Template{}

	var onceErr error

	onceBody := func() {
		for _, v := range filters {
			tpl, err := template.New("matcher").Parse(v)
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
		err := t.Execute(&buf, a)
		if err != nil {
			return false, err
		}
		// fmt.Printf("trying to match %s against %s\n", file, buf.String())

		if strings.HasSuffix(file, buf.String()) {
			// fmt.Printf("file %s matches filters\n", file)
			return true, nil
		}
	}

	return false, nil
}

// Render a passed-in template agains args
func (a Args) Render(t string) (string, error) {
	tmpl, err := template.New("download").Parse(t)
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, a)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
