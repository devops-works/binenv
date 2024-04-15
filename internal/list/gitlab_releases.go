package list

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var (
	// ErrGitlabRateLimitClose is returned when we are close to a ratelimit
	ErrGitlabRateLimitClose = errors.New("gitlab rate limit is close")

	// ErrGitlabRateLimited is returned when we are close to a ratelimit
	ErrGitlabRateLimited = errors.New("rate limited")

	// GitlabLowRateLimit is the low boundary for rate limits
	GitlabLowRateLimit = 4
)

type glReleaseResponse []struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

// GitlabRelease contains what is required to get a list of release from Gitlab
type GitlabRelease struct {
	url         string
	prefix      string
	exclude     string
	versionFrom string
}

// Get returns a list of available versions
func (g GitlabRelease) Get(ctx context.Context) ([]string, error) {
	// logger := zerolog.Ctx(ctx).With().Str("func", "GitlabRelease.Get").Logger()

	var (
		next     = 1
		versions []string
		v        []string
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

func (g GitlabRelease) doGet(ctx context.Context, page int) ([]string, int, error) {
	logger := zerolog.Ctx(ctx).With().Str("func", "GitlabRelease.doGet").Logger()

	next := 0
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?page=%d", g.url, page), nil)
	if err != nil {
		return nil, 0, err
	}

	if token := os.Getenv("GITLAB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	// // Check if we are already rate limited
	if isGitlabRateLimited(resp) {
		return nil, 0, handleGitlabRatelimit(resp)
	}

	// fmt.Println(string(body))
	gr := glReleaseResponse{}
	err = json.Unmarshal([]byte(body), &gr)
	if err != nil {
		logger.Error().Err(err).Msgf("error unmarshalling gitlab response for %s", g.url)
		return nil, 0, err
	}
	// fmt.Printf("%v\n", gr)

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

	if isGitlabRateLimitClose(resp) {
		return versions, next, handleGitlabRatelimit(resp)
	}
	return versions, next, nil
}

func gitlabRateLimit(resp *http.Response) int {
	if v, ok := resp.Header["Ratelimit-Remaining"]; ok {
		remain, err := strconv.Atoi(v[0])
		if err != nil {
			// Could not find current usage
			return 0
		}
		return remain
	}
	// No header found, assume no limit
	return -1
}

func isGitlabRateLimited(resp *http.Response) bool {
	return gitlabRateLimit(resp) == 0
}

func isGitlabRateLimitClose(resp *http.Response) bool {
	return gitlabRateLimit(resp) <= GitlabLowRateLimit
}

func handleGitlabRatelimit(resp *http.Response) error {
	// We received rate limit information
	// Ratelimit-Limit: 60
	// Ratelimit-Remaining: 50
	// Ratelimit-Reset: 1598169973
	// see https://docs.gitlab.com/ee/user/admin_area/settings/user_and_ip_rate_limits.html#use-a-custom-rate-limit-response
	rr, ok := resp.Header["Ratelimit-Remaining"]
	if !ok {
		return nil
	}
	remain, err := strconv.Atoi(rr[0])
	if err != nil {
		// Could not find current usage
		// But return the ate liit error anyway
		return ErrGitlabRateLimited
	}

	rl := resp.Header["Ratelimit-Limit"]
	limit, err := strconv.Atoi(rl[0])
	if err != nil {
		// Could not find current usage
		// But return the ate liit error anyway
		return fmt.Errorf("error checking ratelimit limit: %v (%w)", err, ErrGitlabRateLimited)
	}

	rs := resp.Header["Ratelimit-Reset"]
	reset, err := strconv.ParseInt(rs[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error checking ratelimit reset: %v (%w)", err, ErrGitlabRateLimited)
	}
	okdate := time.Unix(reset, 0)

	if remain == 0 {
		return fmt.Errorf("%w: rate limited by gitlab; please retry after %s",
			ErrGitlabRateLimited,
			okdate,
		)
	}
	return fmt.Errorf("%w: remaining %d of %d; please retry after %s",
		ErrGitlabRateLimitClose,
		remain,
		limit,
		okdate,
	)
}
