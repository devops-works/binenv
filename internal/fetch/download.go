package fetch

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/devops-works/binenv/internal/tpl"
)

// Download handles direct binary releases
type Download struct {
	url string
}

// Fetch gets the package and returns location of downloaded file
func (d Download) Fetch(ctx context.Context, v string) (string, error) {
	args := tpl.New(v)

	url, err := args.Render(d.url)
	if err != nil {
		return "", err
	}

	fmt.Printf("fetching version %q for arch %q and OS %q at %s\n", v, runtime.GOARCH, runtime.GOOS, url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to download binary at %s: %s", url, resp.Status)
	}

	tmpfile, err := ioutil.TempFile("", v)
	if err != nil {
		log.Fatal(err)
	}

	defer tmpfile.Close()

	// Write the body to file
	_, err = io.Copy(tmpfile, resp.Body)

	return tmpfile.Name(), err
}
