package install

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/xi2/xz"

	"github.com/devops-works/binenv/internal/mapping"
	"github.com/devops-works/binenv/internal/tpl"
)

// TarXZ handles xz files
type TarXZ struct {
	filters []string
}

// Install file from xz file
func (x TarXZ) Install(src, dst, version string, mapper mapping.Mapper) error {
	noMatches := ErrNoMatch

	data, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}

	r, err := xz.NewReader(bytes.NewReader(data), 0)
	if err != nil {
		return err
	}

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
