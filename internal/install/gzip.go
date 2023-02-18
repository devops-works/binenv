package install

import (
	"compress/gzip"
	"io"
	"os"

	"github.com/devops-works/binenv/internal/mapping"
)

// GZip handles gzip files
type GZip struct {
}

// Install files from gzip file
func (z GZip) Install(src, dst, version string, mapper mapping.Mapper) error {
	fs, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fs.Close()

	in, err := gzip.NewReader(fs)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return nil
}
