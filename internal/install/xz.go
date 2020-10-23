package install

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/xi2/xz"

	"github.com/devops-works/binenv/internal/mapping"
)

// XZ handles xz files
type XZ struct {
}

// Install file from xz file
func (x XZ) Install(src, dst, version string, mapper mapping.Mapper) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}

	r, err := xz.NewReader(bytes.NewReader(data), 0)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}

	return nil
}
