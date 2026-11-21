package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/schollz/progressbar/v3"

	"github.com/devops-works/binenv/internal/mapping"
	"github.com/devops-works/binenv/internal/tpl"
)

// Download handles direct binary releases
type Download struct {
	urls    []string
	headers map[string]string
}

// Fetch gets the package and returns location of downloaded file
func (d Download) Fetch(ctx context.Context, dist, v string, mapper mapping.Mapper) (string, error) {
	logger := zerolog.Ctx(ctx).With().Str("func", "Download.Fetch").Logger()

	var resp *http.Response

	for i, u := range d.urls {

		args := tpl.New(v, mapper)

		url, err := args.Render(u)
		if err != nil {
			if len(d.urls)-1 > i {
				continue
			} else {
				return "", err
			}
		}

		logger.Debug().Msgf("fetching version %q for arch %q and OS %q at %s", v, runtime.GOARCH, runtime.GOOS, url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			if len(d.urls)-1 > i {
				continue
			} else {
				return "", err
			}
		}

		for k, v := range d.headers {
			req.Header.Add(k, v)
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			if len(d.urls)-1 > i {
				continue
			} else {
				return "", err
			}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if len(d.urls)-1 > i {
				logger.Debug().Msgf("unable to download binary at %s: %s, %d urls left to try...", url, resp.Status, len(d.urls)-1-i)
				continue
			} else {
				return "", fmt.Errorf("unable to download binary at %s: %s", url, resp.Status)
			}
		}
		// if we reach this point, download was successful, let's move on
		break
	}

	tmpfile, err := os.CreateTemp("", v)
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
