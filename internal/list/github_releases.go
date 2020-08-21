package list

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

	return versions, nil
}
