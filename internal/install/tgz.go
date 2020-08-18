package install

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"

	"gitlab.com/devopsworks/tools/binenv/internal/tpl"
)

// Zip handles zip files
type Tgz struct {
	filters []string
}

// Install files from tgz file
func (t Tgz) Install(src, dst, version string) error {
	// var filenames []string

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzf)
	args := tpl.New(version)

	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		fmt.Printf("found file %s in tar\n", header.Name)

		switch header.Typeflag {
		case tar.TypeReg: // regular file
			ok, err := args.MatchFilters(header.Name, t.filters)
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			fmt.Printf("found matching file %s size %d\n", header.Name, header.Size)

			out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
			if err != nil {
				return err
			}
			defer out.Close()
			fmt.Printf("writing file to %s\n", dst)
			if _, err := io.Copy(out, tarReader); err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Printf("finished parsing tar")

	return nil
}
