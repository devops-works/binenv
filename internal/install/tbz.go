package install

import (
	"archive/tar"
	"bytes"
	"compress/bzip2"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/devops-works/binenv/internal/mapping"
	"github.com/devops-works/binenv/internal/tpl"
)

// Tbz handles bzip2 files
type Tbz struct {
	filters []string
}

// Install file from bzip2 file
func (x Tbz) Install(src, dst, version string, mapper mapping.Mapper) error {
	noMatches := ErrNoMatch

	data, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}

	r := bzip2.NewReader(bytes.NewReader(data))

	tarReader := tar.NewReader(r)
	args := tpl.New(version, mapper)

	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeReg: // regular file
			ok, err := args.MatchFilters(header.Name, x.filters)
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
			if _, err := io.Copy(out, tarReader); err != nil {
				log.Fatal(err)
			}
			noMatches = nil
		}
	}

	return noMatches

}
