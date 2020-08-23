package list

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrRateLimit is returned when we are close to a ratelimit
	ErrRateLimit = errors.New("rate limit is close")

	// GithubLowRateLimit is the low boundary for rate limits
	GithubLowRateLimit = 4
)

type ghReleaseResponse []struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

// GithubRelease contains what is required to get a list of release from Github
type GithubRelease struct {
	url         string
	prefix      string
	versionFrom string
}

// Get returns a list of available versions
func (g GithubRelease) Get(ctx context.Context) ([]string, error) {
	resp, err := http.Get(g.url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	gr := ghReleaseResponse{}
	err = json.Unmarshal([]byte(body), &gr)
	if err != nil {
		fmt.Printf("error unmarshalling github response for %s: %v\n", g.url, err)
		return nil, err
	}

	versions := []string{}
	for _, v := range gr {
		sv := v.TagName
		switch g.versionFrom {
		case "name":
			sv = v.Name
		}
		if g.prefix == "" {
			versions = append(versions, sv)
			continue
		}
		if strings.HasPrefix(sv, g.prefix) {
			cleanv := strings.TrimPrefix(sv, g.prefix)
			versions = append(versions, cleanv)
		}
	}

	if v, ok := resp.Header["X-Ratelimit-Remaining"]; ok {
		remain, err := strconv.Atoi(v[0])
		if err != nil {
			// Could not find current usage
			// But return the rate limit error anyway
			return nil, ErrRateLimit
		}
		if remain <= GithubLowRateLimit {
			return versions, handleRatelimit(resp)
		}
	}
	return versions, nil
}

func handleRatelimit(resp *http.Response) error {
	// We received rate limit information
	// X-Ratelimit-Limit: 60
	// X-Ratelimit-Remaining: 50
	// X-Ratelimit-Reset: 1598169973
	rr := resp.Header["X-Ratelimit-Remaining"]
	remain, err := strconv.Atoi(rr[0])
	if err != nil {
		// Could not find current usage
		// But return the ate liit error anyway
		return ErrRateLimit
	}

	rl := resp.Header["X-Ratelimit-Limit"]
	limit, err := strconv.Atoi(rl[0])
	if err != nil {
		// Could not find current usage
		// But return the ate liit error anyway
		return fmt.Errorf("error checking ratelimit limit: %v (%w)", err, ErrRateLimit)
	}

	rs := resp.Header["X-Ratelimit-Reset"]
	reset, err := strconv.ParseInt(rs[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error checking ratelimit reset: %v (%w)", err, ErrRateLimit)
	}
	okdate := time.Unix(reset, 0)

	return fmt.Errorf("%w: remaining %d of %d; please retry after %s",
		ErrRateLimit,
		remain,
		limit,
		okdate,
	)
}
