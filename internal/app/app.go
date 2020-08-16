package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.com/devopsworks/tools/binenv/internal/install"
	"gitlab.com/devopsworks/tools/binenv/internal/list"
	"gitlab.com/devopsworks/tools/binenv/internal/release"
	"gopkg.in/yaml.v2"
)

const definitionsURL string = "https://gitlab.com/devopsworks/tools/binenv/-/raw/master/definitions/definitions.yaml"

// App implements the core logic
type App struct {
	def        *Definitions
	installers map[string]install.Installer
	listers    map[string]list.Lister
	fetchers   map[string]release.Fetcher
	cache      map[string][]string
}

// New create a new app instance
func New() *App {
	a := &App{
		installers: make(map[string]install.Installer),
		listers:    make(map[string]list.Lister),
		fetchers:   make(map[string]release.Fetcher),
		cache:      make(map[string][]string),
	}

	err := a.readDefinitions()
	if err != nil {
		fmt.Printf("error reading definitions: %v", err)
		os.Exit(1)
	}

	a.loadCache()

	a.createInstallers()
	a.createListers()
	a.createFetchers()

	// fmt.Printf("%+v\n", a.def)
	return a
}

// GetPackagesListWithPrefix returns a list of packages that starts with prefix
func (a *App) GetPackagesListWithPrefix(pfix string) []string {
	res := []string{}
	for k := range a.listers {
		if strings.HasPrefix(k, pfix) {
			res = append(res, k)
		}
	}

	return res
}

// GetVersionsFromCacheFor returns a list of packages that starts with prefix
func (a *App) GetVersionsFromCacheFor(soft string) []string {
	if val, ok := a.cache[soft]; ok {
		return val
	}

	return []string{}
}

// Distributions list or update available distributions
func (a *App) Distributions() error {
	fmt.Fprintf(os.Stderr, "not implemented yet")
	return nil
}

// Install installs or update a distribution
func (a *App) Install() error {
	fmt.Fprintf(os.Stderr, "not implemented yet")
	return nil
}

// Local sets the locally used version for application
func (a *App) Local() error {
	fmt.Fprintf(os.Stderr, "not implemented yet")
	return nil
}

// Update fetches catalog of applications and updates available versions
func (a *App) Update(which string) error {
	for k, v := range a.listers {
		if which == k || which == "" {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			fmt.Printf("feching available versions for %s...", k)
			versions, err := v.Get(ctx, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to fetch versions for %s: %v\n", k, err)
				continue
			}
			a.cache[k] = versions
		}
	}

	a.saveCache()

	return nil
}

// Versions fetches available versions for the application
func (a *App) Versions(which string) error {
	for k, versions := range a.cache {
		if which == k || which == "" {
			fmt.Printf("\n%s:\n", k)
			for _, v := range versions {
				fmt.Printf("\t%s\n", v)
			}
		}
	}
	return nil
}

func (a *App) readDefinitions() error {
	conf, err := getConfigDir()
	if err != nil {
		return err
	}

	conf += "/definitions.yaml"

	if _, err := os.Stat(conf); os.IsNotExist(err) {
		err := a.fetchDefinitions(conf)
		if err != nil {
			return fmt.Errorf("unable to fetch definitions: %w", err)
		}
		return nil
	}

	yml, err := ioutil.ReadFile(conf)
	if err != nil {
		return fmt.Errorf("unable to read file '%s': %w", conf, err)
	}

	err = yaml.Unmarshal([]byte(yml), &a.def)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) fetchDefinitions(conf string) error {
	resp, err := http.Get(definitionsURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	d, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("unable to guess configuration directory: %w", err)
	}

	err = os.MkdirAll(d, 0750)
	if err != nil {
		return fmt.Errorf("unable to create configuration directory '%s': %w", d, err)
	}

	f, err := os.OpenFile(conf, os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(body)

	err = yaml.Unmarshal([]byte(body), &a.def)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) loadCache() {
	conf, err := getConfigDir()
	if err != nil {
		return
	}

	conf += "/cache.json"

	if _, err := os.Stat(conf); os.IsNotExist(err) {
		return
	}

	js, err := ioutil.ReadFile(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read cache %s: %v\n", conf, err)
		return
	}

	err = json.Unmarshal([]byte(js), &a.cache)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to unmarshal cache %s: %v\n", conf, err)
		return
	}
}

func (a *App) saveCache() {
	conf, err := getConfigDir()
	if err != nil {
		return
	}

	conf += "/cache.json"

	js, err := json.Marshal(&a.cache)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to marshal cache: %v\n", err)
		return
	}

	fd, err := os.OpenFile(conf, os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open cache for writing: %v\n", err)
		return
	}
	defer fd.Close()
	fd.Write(js)
}

func (a *App) createInstallers() {
	for k, v := range a.def.Sources {
		i := v.Install.Factory()
		if i == nil {
			fmt.Fprintf(os.Stderr, "warning: '%s' install method for %s is not implemented\n", v.Install.Type, k)
			continue
		}
		a.installers[k] = i
	}
}

func (a *App) createListers() {
	for k, v := range a.def.Sources {
		l := v.List.Factory()
		if l == nil {
			fmt.Fprintf(os.Stderr, "warning: '%s' list method for %s is not implemented\n", v.List.Type, k)
			continue
		}
		a.listers[k] = l
	}
}

func (a *App) createFetchers() {
	for k, v := range a.def.Sources {
		f := v.Release.Factory()
		if f == nil {
			fmt.Fprintf(os.Stderr, "warning: '%s' release method for %s is not implemented\n", v.Release.Type, k)
			continue
		}
		a.fetchers[k] = f
	}
}
