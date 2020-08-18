package list

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type ghReleaseResponse []struct {
	TagName string `json:"tag_name"`
}

// GithubRelease contains what is required to get a list of release from Github
type GithubRelease struct {
	url string
}

// Get returns a list of available versions
func (g GithubRelease) Get(ctx context.Context, wg *sync.WaitGroup) ([]string, error) {
	if wg != nil {
		defer wg.Done()
	}

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
		versions = append(versions, v.TagName)
	}

	return versions, nil
}
