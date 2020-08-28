package fetch

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/schollz/progressbar/v3"

	"github.com/devops-works/binenv/internal/mapping"
	"github.com/devops-works/binenv/internal/tpl"
)

// Download handles direct binary releases
type Download struct {
	url string
}

// Fetch gets the package and returns location of downloaded file
func (d Download) Fetch(ctx context.Context, dist, v string, mapper mapping.Mapper) (string, error) {
	logger := zerolog.Ctx(ctx).With().Str("func", "GithubRelease.Get").Logger()

	args := tpl.New(v, mapper)

	url, err := args.Render(d.url)
	if err != nil {
		return "", err
	}

	logger.Debug().Msgf("fetching version %q for arch %q and OS %q at %s", v, runtime.GOARCH, runtime.GOOS, url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to download binary at %s: %s", url, resp.Status)
	}

	tmpfile, err := ioutil.TempFile("", v)
	if err != nil {
		logger.Fatal().Err(err)
	}

	defer tmpfile.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fmt.Sprintf("fetching %s version %s", dist, v),
	)
	io.Copy(io.MultiWriter(tmpfile, bar), resp.Body)

	// Write the body to file
	_, err = io.Copy(tmpfile, resp.Body)

	return tmpfile.Name(), err
}
