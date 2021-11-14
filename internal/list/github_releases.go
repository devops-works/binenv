package list

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
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
	exclude     string
	versionFrom string
}

// Get returns a list of available versions
func (g GithubRelease) Get(ctx context.Context) ([]string, error) {
	// logger := zerolog.Ctx(ctx).With().Str("func", "GithubRelease.Get").Logger()

	var (
		next     = 1
		versions = []string{}
		v        = []string{}
		err      error
	)

	for next > 0 {
		v, next, err = g.doGet(ctx, next)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v...)
	}

	return versions, err
}

func (g GithubRelease) doGet(ctx context.Context, page int) ([]string, int, error) {
	logger := zerolog.Ctx(ctx).With().Str("func", "GithubRelease.doGet").Logger()

	next := 0
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?page=%d", g.url, page), nil)
	if err != nil {
		return nil, 0, err
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	// Check if we are already rate limited
	if isRateLimited(resp) {
		return nil, 0, handleRatelimit(resp)
	}

	gr := ghReleaseResponse{}
	err = json.Unmarshal([]byte(body), &gr)
	if err != nil {
		logger.Error().Err(err).Msgf("error unmarshalling github response for %s", g.url)
		return nil, 0, err
	}

	var re *regexp.Regexp
	if g.exclude != "" {
		re, err = regexp.Compile(g.exclude)
		if err != nil {
			logger.Error().Err(err).Msgf("error compiling regular expression %q", g.exclude)
			return nil, 0, err
		}
	}

	versions := []string{}

	for _, v := range gr {
		sv := v.TagName
		switch g.versionFrom {
		case "name":
			sv = v.Name
		}

		if re != nil && re.Match([]byte(sv)) {
			logger.Debug().Msgf("skipping version %q excluded by exclude regexp %q", sv, g.exclude)
			continue
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

	if len(resp.Header["Link"]) > 0 && strings.Contains(resp.Header["Link"][0], "rel=\"next\"") {
		re := regexp.MustCompile(`page=(\d*)>; rel="next"`)
		match := re.FindStringSubmatch(resp.Header["Link"][0])
		next, err = strconv.Atoi(match[1])

		if err != nil {
			return nil, 0, err
		}
	}

	if isRateLimitClose(resp) {
		return versions, next, handleRatelimit(resp)
	}
	return versions, next, nil
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
	return rateLimit(resp) == 0
}

func isRateLimitClose(resp *http.Response) bool {
	return rateLimit(resp) <= GithubLowRateLimit
}

func handleRatelimit(resp *http.Response) error {
	// We received rate limit information
	// X-Ratelimit-Limit: 60
	// X-Ratelimit-Remaining: 50
	// X-Ratelimit-Reset: 1598169973
	rr, ok := resp.Header["X-Ratelimit-Remaining"]
	if !ok {
		return nil
	}
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
