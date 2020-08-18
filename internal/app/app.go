package app

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	gov "github.com/hashicorp/go-version"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"gitlab.com/devopsworks/tools/binenv/internal/install"
	"gitlab.com/devopsworks/tools/binenv/internal/list"
	"gitlab.com/devopsworks/tools/binenv/internal/release"
	"gopkg.in/yaml.v2"
)

const distributionsURL string = "https://gitlab.com/devopsworks/tools/binenv/-/raw/master/definitions/distributions.yaml"

// App implements the core logic
type App struct {
	def        *Distributions
	installers map[string]install.Installer
	listers    map[string]list.Lister
	fetchers   map[string]release.Fetcher
	cache      map[string][]string

	bindir string
	logger *log.Logger
}

// New create a new app instance
func New(o ...func(*App) error) (*App, error) {
	d, err := homedir.Dir()
	if err != nil {
		d = "~"
	}
	d = filepath.Join(d, "/.binenv/")

	a := &App{
		installers: make(map[string]install.Installer),
		listers:    make(map[string]list.Lister),
		fetchers:   make(map[string]release.Fetcher),
		cache:      make(map[string][]string),
		bindir:     d,
		logger:     log.New(),
	}

	// Apply functional options
	for _, f := range o {
		if err := f(a); err != nil {
			return nil, err
		}
	}

	if strings.HasSuffix(os.Args[0], "binenv") {
		err = a.selfInstall()
		if err != nil {
			a.logger.Errorf("unable to set-up myself: %v", err)
			os.Exit(1)
		}
	}

	err = a.readDistributions()
	if err != nil {
		a.logger.Errorf("unable to read distributions: %v", err)
		os.Exit(1)
	}

	a.loadCache()

	a.createListers()
	a.createFetchers()
	a.createInstallers()

	return a, nil
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

// GetInstalledVersionsFor returns a sorted list of versionsfor distribution
func (a *App) GetInstalledVersionsFor(dist string) []string {
	if _, err := os.Stat(a.getBinDirFor(dist)); os.IsNotExist(err) {
		return []string{}
	}

	versions := []string{}

	err := filepath.Walk(a.getBinDirFor(dist), func(path string, info os.FileInfo, err error) error {
		if a.getBinDirFor(dist) != path {
			versions = append(versions, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		a.logger.Errorf("unable to fetch versions for %q: %v", dist, err)
		return []string{}
	}

	fmt.Printf("%+v\n", versions)
	versionsV := make([]*gov.Version, len(versions))
	for i, raw := range versions {
		v, _ := gov.NewVersion(raw)
		versionsV[i] = v
	}

	fmt.Printf("%+v\n", versionsV)
	sort.Sort(sort.Reverse(gov.Collection(versionsV)))
	versions = []string{}
	for _, v := range versionsV {
		versions = append(versions, v.String())
	}

	return versions
}

// GetVersionsFromCacheFor returns a list of versions available for distribution
func (a *App) GetVersionsFromCacheFor(dist string) []string {
	if val, ok := a.cache[dist]; ok {
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
func (a *App) Install(dist, version string) error {
	// Check if distribution is managed by us
	if a.fetchers[dist] == nil {
		return fmt.Errorf("no fetcher found for %s", dist)
	}
	if _, ok := a.fetchers[dist]; !ok {
		a.logger.Errorf("no such distribution %q", dist)
		return nil
	}

	versions := a.GetInstalledVersionsFor(dist)

	// If version is specified, check if it exists, return if yes
	if version != "" {
		if stringInSlice(gov.Must(gov.NewVersion(version)).String(), versions) {
			a.logger.Warnf("version %q already installed for %q", version, dist)
			return nil
		}
	}

	// Call fetcher for distribution
	file, err := a.fetchers[dist].Fetch(context.Background(), version)
	if err != nil {
		return err
	}

	// Create destination directory
	if _, err := os.Stat(a.getBinDirFor(dist)); os.IsNotExist(err) {
		err := os.MkdirAll(a.getBinDirFor(dist), 0750)
		if err != nil {
			return err
		}
	}

	if a.installers[dist] == nil {
		return fmt.Errorf("no installer found for %s", dist)
	}
	err = a.installers[dist].Install(
		file,
		filepath.Join(
			a.getBinDirFor(dist),
			gov.Must(gov.NewVersion(version)).String(),
		),
		version,
	)
	if err != nil {
		log.Errorf("unable to install %s version %s: %v", dist, version, err)
		return err
	}

	err = a.CreateShimFor(dist)
	if err != nil {
		return err
	}

	fmt.Printf("%s version %s installed\n", dist, version)
	return nil
}

// Uninstall removes a distribution version or the complete distribution
func (a *App) Uninstall(args []string) error {
	// TODO: Check if distribution is managed by us
	// TODO: If version is specified, check if it exists
	// TODO: If version is specified, remove without confirmation
	// TODO: If version is not specified, ask confirmation and remove distribution
	fmt.Fprintf(os.Stderr, "not implemented yet")
	return nil
}

// Local sets the locally used version for application
func (a *App) Local(distribution, version string) error {
	// TODO: Check if distribution is managed by us
	// TODO: Check if version is available
	// TODO: Open local .binenv.lock if exists or create
	// TODO: Replace or create entry for distribution
	fmt.Fprintf(os.Stderr, "not implemented yet")
	return nil
}

// Update fetches catalog of applications and updates available versions
func (a *App) Update(which string) error {
	err := a.readDistributions()
	if err != nil {
		a.logger.Errorf("unable to read distributions: %v", err)
		os.Exit(1)
	}

	for k, v := range a.listers {
		if which == k || which == "" {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			a.logger.Debugf("feching available versions for %q", k)
			versions, err := v.Get(ctx, nil)
			if err != nil {
				a.logger.Errorf("unable to fetch versions for %q: %v", k, err)
				continue
			}
			a.logger.Debugf("found versions %q", strings.Join(versions, ","))
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

// CreateShimFor creates a shim for the distribution
func (a *App) CreateShimFor(dist string) error {
	// Should not happen
	shim := filepath.Join(a.bindir, "/shim")
	if _, err := os.Stat(shim); os.IsNotExist(err) {
		return fmt.Errorf("unable to find shim file: %w", err)
	}

	lnk := filepath.Join(a.bindir, dist)
	if _, err := os.Stat(lnk); os.IsNotExist(err) {
		err := os.Symlink(shim, lnk)
		if err != nil {
			return err
		}
	}
	return nil
}

// Execute runs the shim function that executes real distributions
func (a *App) Execute(args []string) {
	fmt.Println("execute")
	// Check if args[0] is managed by us. If not write an error and exit. This
	// should not happen since, if we are here, we must have used a symlink to
	// the shim.
	versions := a.GetInstalledVersionsFor(args[0])
	if len(versions) == 0 {
		log.Errorf("no versions found for distribution %q. Something is really odd.", os.Args[0])
	}

	// Check version to use, going up to home directory if needeed and
	// try /etc/binenv
	// Take the first match while going up
	curdir, _ := os.Getwd()
	version := a.GuessBestVersionFor(args[0], curdir, versions)

	// If we did not find any this might be:
	// - because we do not have a version that matches contraints
	// - because we did not find any contraint to apply
	// For now, we use the latest version available
	// TODO: ofc this is not good, fix this and distinguish both cases
	if version == "" {
		version = versions[len(versions)-1]
	}

	bd := a.getBinDirFor(args[0])
	binary := filepath.Join(bd, version)

	fmt.Printf("executing %q\n", binary)

	if err := syscall.Exec(binary, args, os.Environ()); err != nil {
		fmt.Println(err)
	}
}

func (a *App) selfInstall() error {
	err := os.MkdirAll(a.bindir, 0750)
	if err != nil {
		return err
	}

	self, err := os.Executable()
	if err != nil {
		return err
	}

	from, err := os.Open(self)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(filepath.Join(a.bindir, "/shim"), os.O_RDWR|os.O_CREATE, 0750)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) readDistributions() error {
	conf, err := getConfigDir()
	if err != nil {
		return err
	}

	conf = filepath.Join(conf, "/distributions.yaml")

	if _, err := os.Stat(conf); os.IsNotExist(err) {
		err := a.fetchDistributions(conf)
		if err != nil {
			return fmt.Errorf("unable to fetch distributions: %w", err)
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

func (a *App) fetchDistributions(conf string) error {
	resp, err := http.Get(distributionsURL)
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

	conf = filepath.Join(conf, "/cache.json")
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

func (a *App) getBinDirFor(dist string) string {
	return filepath.Join(a.bindir, "binaries/", dist)
}

// GuessBestVersionFor returns closest version requirement given a location,
// a distribution and a version list.
// If no match we return the latest version we have
func (a *App) GuessBestVersionFor(dist, dir string, versions []string) string {
	home, _ := homedir.Dir()
	home = filepath.Clean(home)
	dir = filepath.Clean(dir)

	for {
		// fmt.Printf("in directory %s\n", dir)
		if _, err := os.Stat(filepath.Join(dir, ".binenv.lock")); os.IsNotExist(err) {
			// If in homedir, we found nothing
			if dir == home {
				return ""
			}
			// Move up
			dir = filepath.Clean(filepath.Join(dir, ".."))
			// fmt.Printf("new directory %s\n", dir)
			continue
		}

		// lock file is found
		f, err := os.Open(filepath.Join(dir, ".binenv.lock"))
		if err != nil {
			return ""
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			// fmt.Printf("got line %s\n", line)

			if strings.HasPrefix(line, dist) {
				constraint := strings.TrimPrefix(line, dist)
				for _, v := range versions {
					// fmt.Printf("testing version %s for %s\n", v, dist)
					v1, _ := gov.NewVersion(v)
					// Constraints example.
					constraints, _ := gov.NewConstraint(constraint)
					if constraints.Check(v1) {
						// fmt.Printf("%s satisfies constraints %s\n", v1, constraints)
						return v1.String()
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return ""
		}
		return ""
	}
}

func (a *App) saveCache() {
	conf, err := getConfigDir()
	if err != nil {
		return
	}

	conf = filepath.Join(conf, "/cache.json")

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
		i := v.Install.Factory(v.Install.Binaries)
		if i == nil {
			a.logger.Warnf("warning: '%s' install method for %s is not implemented\n", v.Install.Type, k)
			continue
		}
		a.installers[k] = i
	}
}

func (a *App) createListers() {
	for k, v := range a.def.Sources {
		l := v.List.Factory()
		if l == nil {
			a.logger.Warnf("warning: '%s' list method for %s is not implemented\n", v.List.Type, k)
			continue
		}
		a.listers[k] = l
	}
}

func (a *App) createFetchers() {
	for k, v := range a.def.Sources {
		f := v.Release.Factory()
		if f == nil {
			a.logger.Warnf("warning: '%s' release method for %s is not implemented\n", v.Release.Type, k)
			continue
		}
		a.fetchers[k] = f
	}
}

// Functional options

// WithDiscard sets the repository for this service
func WithDiscard() func(*App) error {
	return func(a *App) error {
		return a.setLogOutput(ioutil.Discard)
	}
}

func (a *App) setLogOutput(w io.Writer) error {
	a.logger.Out = w

	return nil
}

// WithBinDir sets the binaries directory
func WithBinDir(dir string) func(*App) error {
	return func(a *App) error {
		return a.SetBinDir(dir)
	}
}

// SetBinDir sets bin directory to use
func (a *App) SetBinDir(d string) error {
	a.bindir = d

	return nil
}

// WithLogLevel sets the binaries directory
func WithLogLevel(l string) func(*App) error {
	return func(a *App) error {
		return a.SetLogLevel(l)
	}
}

// SetLogLevel sets bin directory to use
func (a *App) SetLogLevel(l string) error {
	level, err := log.ParseLevel(l)
	if err != nil {
		a.logger.Fatalf("invalid log level %q: %v", l, err)
	}
	a.logger.SetLevel(level)

	return nil
}
