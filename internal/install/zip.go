package install

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	"gitlab.com/devopsworks/tools/binenv/internal/tpl"
)

// Zip handles zip files
type Zip struct {
	filters []string
}

// Install files from zip file
func (z Zip) Install(src, dst, version string) error {
	// var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	args := tpl.New(version)
	for _, f := range r.File {
		fmt.Printf("found file %s in zip\n", f.Name)

		// fpath := filepath.Join(dst, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		// if !strings.HasPrefix(dst, filepath.Clean(dst)+string(os.PathSeparator)) {
		// 	return fmt.Errorf("%s: illegal file path", dst)
		// }

		ok, err := args.MatchFilters(f.Name, z.filters)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		fmt.Printf("installing in %s\n", dst)

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
