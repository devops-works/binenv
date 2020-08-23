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
	// ErrGithubRateLimitClose is returned when we are close to a ratelimit
	ErrGithubRateLimitClose = errors.New("github rate limit is close")

	// ErrGithubRateLimited is returned when we are close to a ratelimit
	ErrGithubRateLimited = errors.New("rate limited")

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

	// Checkif we are already rate limited
	if isRateLimited(resp) {
		return nil, handleRatelimit(resp)
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

	if isRateLimitClose(resp) {
		return versions, handleRatelimit(resp)
	}
	return versions, nil
}

func rateLimit(resp *http.Response) int {
	if v, ok := resp.Header["X-Ratelimit-Remaining"]; ok {
		remain, err := strconv.Atoi(v[0])
		if err != nil {
			// Could not find current usage
			return 0
		}
		return remain
	}
	return 0
}

func isRateLimited(resp *http.Response) bool {
	if rateLimit(resp) == 0 {
		return true
	}
	return false
}

func isRateLimitClose(resp *http.Response) bool {
	if rateLimit(resp) <= GithubLowRateLimit {
		return true
	}
	return false
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
		return ErrGithubRateLimited
	}

	rl := resp.Header["X-Ratelimit-Limit"]
	limit, err := strconv.Atoi(rl[0])
	if err != nil {
		// Could not find current usage
		// But return the ate liit error anyway
		return fmt.Errorf("error checking ratelimit limit: %v (%w)", err, ErrGithubRateLimited)
	}

	rs := resp.Header["X-Ratelimit-Reset"]
	reset, err := strconv.ParseInt(rs[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error checking ratelimit reset: %v (%w)", err, ErrGithubRateLimited)
	}
	okdate := time.Unix(reset, 0)

	if remain == 0 {
		return fmt.Errorf("%w: rate limited by github; please retry after %s",
			ErrGithubRateLimited,
			okdate,
		)
	}
	return fmt.Errorf("%w: remaining %d of %d; please retry after %s",
		ErrGithubRateLimitClose,
		remain,
		limit,
		okdate,
	)
}
