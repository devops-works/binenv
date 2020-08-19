package install

import (
	"archive/zip"
	"io"
	"os"

	"github.com/devops-works/binenv/internal/mapping"
	"github.com/devops-works/binenv/internal/tpl"
)

// Zip handles zip files
type Zip struct {
	filters []string
}

// Install files from zip file
func (z Zip) Install(src, dst, version string, mapper mapping.Mapper) error {
	// var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	args := tpl.New(version, mapper)
	for _, f := range r.File {
		ok, err := args.MatchFilters(f.Name, z.filters)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
		if err != nil {
			return err
		}
		defer out.Close()

		in, err := f.Open()
		if err != nil {
			return err
		}
		defer in.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
	}
	return nil
}
