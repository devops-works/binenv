package app

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"

	gov "github.com/hashicorp/go-version"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v2"

	"github.com/devops-works/binenv/internal/fetch"
	"github.com/devops-works/binenv/internal/install"
	"github.com/devops-works/binenv/internal/list"

	"github.com/logrusorgru/aurora"

	"github.com/devops-works/binenv/internal/mapping"
)

// Where to fetch distribution list and cached version
// Distributions are fetched from master
// Cache is fetched from develop
// Since master is always behind (or at) develop, it ensures we always have
// cache for listed entries
const distributionsURL string = "https://raw.githubusercontent.com/devops-works/binenv/master/distributions/distributions.yaml"
const cacheURL string = "https://raw.githubusercontent.com/devops-works/binenv/develop/distributions/cache.json"

// App implements the core logic
type App struct {
	def        *Distributions
	mappers    map[string]mapping.Remapper
	installers map[string]install.Installer
	listers    map[string]list.Lister
	fetchers   map[string]fetch.Fetcher
	cache      map[string][]string
	flags      flags

	dryrun      bool
	global      bool
	concurrency int

	bindir    string
	linkdir   string
	cachedir  string
	configdir string

	logger zerolog.Logger
}

// flags holds App boolean flags
type flags struct {
	Setters    map[string](func(*flags, bool))
	justExpand bool
}

type jobResult struct {
	distribution string
	versions     []string
}

var (
	// ErrAlreadyInstalled is returned when the requested version is already installed
	ErrAlreadyInstalled = errors.New("version already installed")
)

// Init prepares App for use
func (a *App) Init(o ...func(*App) error) (*App, error) {
	a.mappers = make(map[string]mapping.Remapper)
	a.installers = make(map[string]install.Installer)
	a.listers = make(map[string]list.Lister)
	a.fetchers = make(map[string]fetch.Fetcher)
	a.cache = make(map[string][]string)
	a.logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()

	// Default to warn log level
	a.logger = a.logger.Level(zerolog.InfoLevel)

	// Apply functional options
	for _, f := range o {
		if err := f(a); err != nil {
			return nil, err
		}
	}

	a.DumpConfig()

	err := a.readDistributions()
	if err != nil {
		a.logger.Error().Err(err).Msgf("unable to read distributions")
		os.Exit(1)
	}

	a.loadCache()

	a.createMappers()
	a.createListers()
	a.createFetchers()
	a.createInstallers()

	a.initializeFlags()
	return a, nil
}

// Search show a list returns a list of packages contains string
// in name or description
func (a *App) Search(pfix string, wide bool) []string {
	res := []string{}
	pfix = strings.ToLower(pfix)
	// First add matching distributions
	for k, v := range a.def.Sources {
		if strings.Contains(strings.ToLower(k), pfix) {
			res = append(res, k)
		}
		if strings.Contains(strings.ToLower(v.Description), pfix) {
			res = append(res, k)
		}
	}

	sort.Strings(res)
	previous := ""
	for _, d := range res {
		if d == previous {
			// Skip duplicate
			continue
		}
		previous = d

		if wide {
			fmt.Printf("%s,%s,%q\n",
				aurora.Bold(d),
				strings.TrimSpace(a.def.Sources[d].URL),
				aurora.Faint(strings.TrimSpace(a.def.Sources[d].Description)),
			)
		} else {
			fmt.Printf("%s: %s\n",
				aurora.Bold(d),
				aurora.Faint(strings.TrimSpace(a.def.Sources[d].Description)),
			)
		}
	}

	return res
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

// GetMostRecent returns the most recent stable available version
func (a *App) GetMostRecent(dist string) string {
	availVersions := a.GetAvailableVersionsFor(dist)
	for _, v := range availVersions {
		if gov.Must(gov.NewVersion(v)).Prerelease() == "" {
			return v
		}
	}
	return ""
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
		a.logger.Error().Err(err).Msgf("unable to fetch installed versions for %q", dist)
		return []string{}
	}

	versionsV := make([]*gov.Version, len(versions))
	for i, raw := range versions {
		v, _ := gov.NewVersion(raw)
		versionsV[i] = v
	}

	sort.Sort(sort.Reverse(gov.Collection(versionsV)))
	versions = []string{}
	for _, v := range versionsV {
		versions = append(versions, v.String())
	}

	return versions
}

// GetAvailableVersionsFor returns a list of versions available for distribution
func (a *App) GetAvailableVersionsFor(dist string) []string {
	if _, ok := a.cache[dist]; !ok {
		return []string{}
	}

	versions := a.cache[dist]
	versionsV := make([]*gov.Version, len(versions))
	for i, raw := range versions {
		v, _ := gov.NewVersion(raw)
		versionsV[i] = v
	}

	sort.Sort(sort.Reverse(gov.Collection(versionsV)))
	versions = []string{}
	for _, v := range versionsV {
		versions = append(versions, v.String())
	}

	return versions
}

// InstallFromLock install distributions/versions to match the local
// .binenv.lock file
func (a *App) InstallFromLock() error {
	// Get listed versions from local .binenv.lock
	curdir, err := os.Getwd()
	if err != nil {
		a.logger.Error().Err(err).Msg("unable to determine current directory")
		return fmt.Errorf("unable to determine current directory: %w", err)
	}
	distributions, lines := a.getDistributionsFromLock()

	// Lets loop on each distribution and find the best versions among
	// available versions
	for i, d := range distributions {
		available := a.GetAvailableVersionsFor(d)
		required, _ := a.GuessBestVersionFor(d, curdir, curdir, available)
		installed := a.GetInstalledVersionsFor(d)

		if required == "" {
			a.logger.Warn().Msgf(`no available versions found for %q. Please run "binenv update %s".`, d, d)
			continue
		}
		if !stringInSlice(required, installed) {
			a.logger.Warn().Msgf("installing %q (%s) to satisfy constraint %q", d, required, lines[i])
			a.install(d, required)
		} else {
			a.logger.Debug().Msgf("will use %q (%s) to satisfy constraint %q", d, required, lines[i])
		}
	}

	return nil
}

// Install installs or update a distribution
func (a *App) Install(specs ...string) error {
	if len(specs)%2 != 0 && len(specs) != 1 {
		a.logger.Error().Msg("invalid number of arguments (must have distribution and version pairs")
		os.Exit(1)
	}

	errored := false

	for i := 0; i < len(specs); i += 2 {
		dist := specs[i]
		version := ""
		if len(specs) > 1 {
			version = specs[i+1]
		}

		supportedPlatforms := a.def.Sources[dist].SupportedPlatforms
		if len(supportedPlatforms) > 0 {
			isSupported := false
			for _, platform := range supportedPlatforms {
				if platform.OS == runtime.GOOS && platform.Arch == runtime.GOARCH {
					isSupported = true
					break
				}
			}
			if !isSupported {
				a.logger.Error().Msgf("%q is not available for %s/%s", dist, runtime.GOOS, runtime.GOARCH)
				errored = true
				continue
			}
		}

		v, err := a.install(dist, version)
		if err != nil && !errors.Is(err, ErrAlreadyInstalled) {
			a.logger.Error().Err(err).Msgf("unable to install %q (%s)", dist, v)
			errored = true
			continue
		}
		if err == nil {
			a.logger.Info().Msgf("%q (%s) installed", dist, v)
		}

		postInstallMessage := strings.TrimSpace(a.def.Sources[dist].PostInstallMessage)
		if len(postInstallMessage) > 0 {
			fmt.Printf("===> Post install message for %s <===\n%s\n", aurora.Bold(dist), postInstallMessage)
		}
	}

	if errored {
		os.Exit(1)
	}

	return nil
}

func (a *App) install(dist, version string) (string, error) {
	// Check if distribution is managed by us
	if a.fetchers[dist] == nil {
		return "", fmt.Errorf("no fetcher found for %q", dist)
	}
	if _, ok := a.fetchers[dist]; !ok {
		a.logger.Error().Msgf("no such distribution %q", dist)
		return "", nil
	}

	versions := a.GetInstalledVersionsFor(dist)

	// If version is not specified, install most recent
	if version == "" {
		version = a.GetMostRecent(dist)
		a.logger.Warn().Msgf("version for %q not specified; using %q", dist, version)
	}

	if version == "" {
		return "", fmt.Errorf("unable to select latest stable version for %q: no stable version available. May be run 'binenv update %s' ?", dist, dist)
	}

	// If version is specified, check if it exists, return if yes
	cleanVersion, err := gov.NewSemver(version)
	if err != nil {
		return "", err
	}
	version = cleanVersion.String()
	if stringInSlice(version, versions) {
		a.logger.Warn().Msgf("version %q already installed for %q", version, dist)
		return version, ErrAlreadyInstalled
	}

	var m mapping.Mapper
	{
		if v, ok := a.mappers[dist]; ok {
			m = v
		}
	}

	ctx := a.logger.WithContext(context.TODO())
	if a.dryrun {
		a.logger.Warn().Msgf("dry-run mode: skipping install for %q (%s)", dist, version)
		return version, nil
	}

	// Call fetcher for distribution
	file, err := a.fetchers[dist].Fetch(ctx, dist, version, m)
	if err != nil {
		return version, err
	}

	// Create destination directory
	if _, err := os.Stat(a.getBinDirFor(dist)); os.IsNotExist(err) {
		var mode os.FileMode = 0750
		if a.global {
			mode = 0755
		}
		err := os.MkdirAll(a.getBinDirFor(dist), mode)
		if err != nil {
			return version, err
		}
	}

	if a.installers[dist] == nil {
		return version, fmt.Errorf("no installer found for %s", dist)
	}

	if a.dryrun {
		a.logger.Warn().Msgf("dry-run mode: skipping install fir %q (%s)", dist, version)
		return version, nil
	}
	err = a.installers[dist].Install(
		file,
		filepath.Join(
			a.getBinDirFor(dist),
			gov.Must(gov.NewVersion(version)).String(),
		),
		version,
		m,
	)
	if err != nil {
		return version, err
	}

	// Install new shim version if needed
	if dist == "binenv" {
		a.logger.Info().Msgf("executing self install using bindir %s", a.bindir)
		err = a.selfInstall(version)
		if err != nil {
			a.logger.Error().Err(err).Msg("unable to set-up myself")
			os.Exit(1)
		}
	}

	err = a.CreateShimFor(dist)
	if err != nil {
		return version, err
	}

	return version, nil
}

// Uninstall installs or update a distribution
func (a *App) Uninstall(specs ...string) error {
	// We accept either
	// - a single argument (remove all versions for distributions)
	// - an even count of arguments (distribution / version pairs)

	if len(specs)%2 != 0 && len(specs) > 1 {
		a.logger.Fatal().Msg("invalid number of arguments (must have distribution and version pairs)")
	}

	for i := 0; i < len(specs); i += 2 {
		dist := specs[i]
		version := ""
		if len(specs) > 1 {
			version = specs[i+1]
		}
		err := a.uninstall(dist, version)
		if err != nil {
			a.logger.Error().Err(err).Msgf("unable to uninstall %q (%s)", dist, version)
		}
	}
	return nil
}

// Uninstall removes a distribution version or the complete distribution
func (a *App) uninstall(dist, version string) error {
	// Check if distribution is managed by us
	// If version is specified, check if it exists
	installed := a.GetInstalledVersionsFor(dist)
	if version != "" {
		if !stringInSlice(version, installed) {
			return fmt.Errorf("version %q for %q is not installed", version, dist)
		}
		bd := a.getBinDirFor(dist)
		binary := filepath.Join(bd, version)

		// Check this is a version number, just to be sure
		file := filepath.Base(binary)
		if _, err := gov.NewSemver(file); err != nil {
			a.logger.Fatal().Msgf("%q does not look like a binary file installed by binenv; bailing out", file)
		}

		err := os.Remove(binary)
		if err != nil {
			return err
		}

		a.logger.Warn().Msgf("removed version %q for %q", version, dist)
		return nil
	}

	fmt.Printf("WARNING: this will remove *ALL* versions for %q. Type %q to confirm [oh now I changed my mind]: ", dist, dist)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return nil
	}
	response = strings.TrimSpace(response)
	if response != dist {
		fmt.Printf("Cancelled")
		return nil
	}

	for _, v := range installed {
		err := a.uninstall(dist, v)
		// Bail out immediately if there is an error
		if err != nil {
			return err
		}
	}

	lnk := filepath.Join(a.linkdir, dist)
	if err = os.Remove(lnk); err != nil {
		return err
	}

	return nil
}

// Local sets the locally used version for application
func (a *App) Local(distribution, version string) error {
	// TODO: Check if distribution is managed by us
	// TODO: Check if version is available
	// TODO: Open local .binenv.lock if exists or create
	// TODO: Replace or create entry for distribution
	a.logger.Fatal().Msg("not implemented yet")
	return nil
}

// Update fetches catalog of applications and updates available versions
func (a *App) Update(definitions, all bool, nocache bool, which ...string) error {
	if definitions || all {
		conf := filepath.Join(a.configdir, "/distributions.yaml")
		err := a.fetchDistributions(conf)
		if err != nil {
			a.logger.Error().Err(err).Msgf("unable to fetch distributions: %v", err)
			os.Exit(1)
		}

		// Return if only definitions were requested
		if definitions {
			return nil
		}
	}

	err := a.readDistributions()
	if err != nil {
		a.logger.Error().Err(err).Msg("unable to read distributions")
		os.Exit(1)
	}

	if len(which) == 0 {
		for k := range a.listers {
			which = append(which, k)
		}
	}

	a.logger.Debug().Msgf("updating %d distributions", len(which))
	if nocache {
		err = a.updateLocally(which...)
	} else {
		err = a.updateGithub()
	}
	if err != nil {
		return err
	}

	err = a.saveCache()

	return err
}

func (a *App) updateGithub() error {
	a.logger.Info().Msgf("retrieving distribution cache from %s", cacheURL)
	resp, err := http.Get(cacheURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(body), &a.cache)
	if err != nil {
		a.logger.Error().Err(err).Msg(`unable to unmarshal Github cache; try to "binenv update" locally`)
		return err
	}

	a.logger.Info().Msgf("fetched updates for %d distributions", len(a.cache))

	return nil
}

func (a *App) fetcher(id int, jobs <-chan string, res chan<- jobResult, timeout time.Duration) {
	for d := range jobs {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		var err error

		r := jobResult{
			distribution: d,
		}

		subctx := a.logger.WithContext(ctx)
		r.versions, err = a.listers[d].Get(subctx)

		if errors.Is(err, list.ErrGithubRateLimitClose) || errors.Is(err, list.ErrGithubRateLimited) {
			a.logger.Error().Err(err).Msgf("rate limit prevents fetching versions for %q", d)
			// return err
			cancel()
			continue
		}
		if err != nil {
			a.logger.Error().Err(err).Msgf("unable to fetch versions for %q", d)
			// continue
		}

		if len(r.versions) == 0 {
			a.logger.Warn().Msgf("found no versions for %q", d)
		} else {
			a.logger.Debug().Msgf("found versions %q for %q", strings.Join(r.versions, ","), d)
		}

		res <- r

		cancel()
	}
}

func (a *App) updateLocally(which ...string) error {
	bar := progressbar.Default(int64(len(which)), "updating distributions")

	jobs := make(chan string)
	res := make(chan jobResult, 1000)
	timeout := 1 * time.Second

	for w := 1; w <= a.concurrency; w++ {
		go a.fetcher(w, jobs, res, timeout)
	}

	count := 0
	for _, d := range which {
		bar.Add(1)

		a.logger.Debug().Msgf("feching available versions for %q", d)
		if _, ok := a.listers[d]; !ok {
			a.logger.Error().Msgf("no distribution named %q", d)
			continue
		}

		jobs <- d
		count++
	}

	close(jobs)

	for c := 0; c < count; c++ {
		r := <-res

		// Skip this entry if no versions are provided
		// see #157, #159, #162...
		if len(r.versions) == 0 {
			a.logger.Warn().Msgf("no versions found for %s; keeping previous versions %s", r.distribution, strings.Join(a.cache[r.distribution], ","))
			continue
		}

		// Flush cache entry if
		a.cache[r.distribution] = []string{}

		// Convert versions to canonical form
		for _, v := range r.versions {
			version, err := gov.NewVersion(v)
			if err != nil {
				a.logger.Debug().Err(err).Msgf("ignoring invalid version for %q", r.distribution)
				continue
			}
			a.cache[r.distribution] = append(a.cache[r.distribution], version.String())
		}
	}
	return nil
}

// Versions fetches available versions for the application
func (a *App) Versions(freezemode bool, specs ...string) error {
	if len(specs) == 0 {
		for k := range a.cache {
			specs = append(specs, k)
		}
	}

	sort.Strings(specs)

	if !freezemode {
		fmt.Printf("# Most recent first; legend: %s, %s, %s\n",
			aurora.Reverse("active"),
			aurora.Bold("installed"),
			aurora.Faint("available"),
		)
	}

	var err error
	for _, s := range specs {
		if !freezemode {
			err = a.versions(s)
		} else {
			err = a.freeze(s)
		}
		if err != nil {
			a.logger.Error().Err(err).Msgf("unable to list versions for %q", s)
		}
	}

	return nil
}

func (a *App) freeze(dist string) error {
	curdir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to determine current directory: %v", err)
	}
	installed := a.GetInstalledVersionsFor(dist)
	a.logger.Debug().Strs("versions", installed).Msgf("installed versions for %s", dist)
	guess, why := a.GuessBestVersionFor(dist, curdir, "", installed)
	a.logger.Debug().Str("guessed", guess).Msgf("guessed version for dist %s", dist)

	if guess != "" {
		fmt.Printf("#%s: %s\n%s=%s\n", dist, why, dist, guess)
	}

	return nil
}

func (a *App) versions(dist string) error {
	curdir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to determine current directory: %v", err)
	}
	available := a.GetAvailableVersionsFor(dist)
	a.logger.Debug().Strs("versions", available).Msgf("available versions for %s", dist)
	installed := a.GetInstalledVersionsFor(dist)
	a.logger.Debug().Strs("versions", installed).Msgf("installed versions for %s", dist)
	guess, why := a.GuessBestVersionFor(dist, curdir, "", installed)
	a.logger.Debug().Str("guessed", guess).Msgf("guessed version for dist %s", dist)

	fmt.Printf("%s: ", dist)

	for _, v := range available {
		// var modifier aurora.Value
		modifier := aurora.Faint(v)
		if stringInSlice(v, installed) {
			if v == guess {
				modifier = aurora.Reverse(fmt.Sprintf("%s (%s)", v, why))
			} else {
				modifier = aurora.Bold(v)
			}
		}
		fmt.Printf("%s ", modifier)
	}
	fmt.Println()
	return nil
}

// Upgrade install last version of all locally installed distributions
func (a *App) Upgrade(ignoreInstallErrors bool) error {
	errored := false
	for dist := range a.cache {

		// ignore uninstalled distribution
		installed := a.GetInstalledVersionsFor(dist)
		if len(installed) == 0 {
			continue
		}

		// get last known version
		version := a.GetMostRecent(dist)

		// call install function
		v, err := a.install(dist, version)

		// process install errors logic.
		if err != nil && !errors.Is(err, ErrAlreadyInstalled) {
			a.logger.Error().Err(err).Msgf("unable to install %q (%s)", dist, v)
			errored = true
			continue
		}

		if err == nil {
			a.logger.Info().Msgf("%q (%s) installed", dist, v)
		}

		postInstallMessage := strings.TrimSpace(a.def.Sources[dist].PostInstallMessage)
		if len(postInstallMessage) > 0 {
			fmt.Printf("===> Post upgrade message for %s <===\n%s\n", aurora.Bold(dist), postInstallMessage)
		}
	}

	if errored && !ignoreInstallErrors {
		os.Exit(1)
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

	lnk := filepath.Join(a.linkdir, dist)
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
	dist := filepath.Base(args[0])

	// Check if args[0] is managed by us. If not write an error and exit. This
	// should not happen since, if we are here, we must have used a symlink to
	// the shim.
	versions := a.GetInstalledVersionsFor(dist)
	if len(versions) == 0 {
		a.logger.Error().Msgf("no versions found for distribution %q (from %s). Something is really odd.", dist, os.Args[0])
	}

	// Check version to use, going up to home directory if needeed and
	// try /etc/binenv
	// Take the first match while going up
	curdir, _ := os.Getwd()
	version, why := a.GuessBestVersionFor(dist, curdir, "", versions)

	// If we did not find any proper version to run
	if version == "" {
		a.logger.Fatal().Msgf("binenv: unable to find proper version for %s (%s)", dist, why)
	}

	bd := a.getBinDirFor(dist)
	binary := filepath.Join(bd, version)

	if a.flags.justExpand {
		fmt.Print(binary)
		return
	}

	if err := syscall.Exec(binary, args, os.Environ()); err != nil {
		a.logger.Fatal().Err(err).Msgf("unable to execute %s", dist)
	}
}

func (a *App) selfInstall(version string) error {
	var mode os.FileMode = 0750
	if a.global {
		mode = 0755
	}

	a.logger.Debug().Msgf("creating bindir %s", a.bindir)
	err := os.MkdirAll(a.bindir, mode)
	if err != nil {
		return err
	}

	bd := a.getBinDirFor("binenv")
	self := filepath.Join(bd, version)

	from, err := os.Open(self)
	if err != nil {
		return err
	}
	defer from.Close()

	shim := filepath.Join(a.bindir, "/shim")
	shimnew := shim + ".new"

	a.logger.Debug().Msgf("installing shim in %s", shim)

	if _, err := os.Stat(shim); os.IsExist(err) {
		shimold := shim + ".old"
		rerr := os.Rename(shim, shimold)
		if rerr != nil {
			return rerr
		}
	}

	to, err := os.OpenFile(shimnew, os.O_RDWR|os.O_CREATE, mode)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	err = os.Rename(shimnew, shim)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) readDistributions() error {
	conf := filepath.Join(a.configdir, "/distributions.yaml")

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
	a.logger.Info().Msg("updating distribution list")
	a.logger.Debug().Msgf("retrieving distribution list from %s", distributionsURL)
	resp, err := http.Get(distributionsURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var mode os.FileMode = 0750
	if a.global {
		mode = 0755
	}
	err = os.MkdirAll(a.configdir, mode)
	if err != nil {
		return fmt.Errorf("unable to create configuration directory '%s': %w", a.configdir, err)
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

func (a *App) getBinDirFor(dist string) string {
	a.logger.Debug().Msgf("finding binaries for %s in %s", dist, filepath.Join(a.bindir, "binaries/", dist))
	return filepath.Join(a.bindir, "binaries/", dist)
}

// GuessBestVersionFor returns closest version requirement given a location,
// a distribution and a version list.
// If no match we return the latest version we have
func (a *App) GuessBestVersionFor(dist, dir, stop string, versions []string) (string, string) {
	// Look at the environment variables in search for the dist version
	// The environment variable need to match [a-zA-Z_]+[a-zA-Z0-9_]*.
	// BINENV is use as prefix and _VERSION as suffix

	envVarName := fmt.Sprintf("BINENV_%v_VERSION", stringToEnvVarName(dist))
	if target := os.Getenv(envVarName); len(target) > 0 {
		for _, v := range versions {
			v1, _ := gov.NewVersion(v)
			// Constraints
			constraints, _ := gov.NewConstraint("=" + target)
			if constraints.Check(v1) {
				return v1.String(), dir
			}
		}

		a.logger.Warn().Msgf(`unable to satisfy constraint %q for %q from environment variable. Ignoring`, target, dist)
	}

	// If stop is "", we enforce stopping in home directory
	if stop == "" {
		stop, _ = homedir.Dir()
	}
	stop = filepath.Clean(stop)
	dir = filepath.Clean(dir)

	if len(versions) == 0 {
		return "", ""
	}

	deflt := versions[0]

	// If no .binenv.lock, try parent until we reach 'stop'
	if _, err := os.Stat(filepath.Join(dir, ".binenv.lock")); os.IsNotExist(err) {
		// If in stop dir, we found nothing
		if dir == stop {
			return deflt, "default"
		}
		// Recurse moving up
		dir = filepath.Clean(filepath.Join(dir, ".."))
		return a.GuessBestVersionFor(dist, dir, filepath.Join(stop, ".."), versions)
	}

	// lock file is found
	f, err := os.Open(filepath.Join(dir, ".binenv.lock"))
	if err != nil {
		return "", ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, dist) {
			constraint := strings.TrimPrefix(line, dist)
			for _, v := range versions {
				v1, _ := gov.NewVersion(v)
				// Constraints
				constraints, _ := gov.NewConstraint(constraint)
				if constraints.Check(v1) {
					return v1.String(), dir
				}
			}
			return "", fmt.Sprintf(`unable to satisfy constraint %q for %q. Try "binenv install -l".`, constraint, dist)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", ""
	}

	// We did not match dist, so return default
	return deflt, "default"
}

// GuessBestVersionFor2 returns closest version requirement given a location,
// a distribution and a version list.
// If no match we return the latest version we have
func (a *App) GuessBestVersionFor2(dist, dir string, versions []string) (string, string) {
	home, _ := homedir.Dir()
	home = filepath.Clean(home)
	dir = filepath.Clean(dir)

	if len(versions) == 0 {
		return "", ""
	}

	deflt := versions[0]

	for {
		if _, err := os.Stat(filepath.Join(dir, ".binenv.lock")); os.IsNotExist(err) {
			// If in homedir, we found nothing
			if dir == home {
				return deflt, "default"
			}
			// Move up
			dir = filepath.Clean(filepath.Join(dir, ".."))
			continue
		}

		// lock file is found
		f, err := os.Open(filepath.Join(dir, ".binenv.lock"))
		if err != nil {
			return "", ""
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, dist) {
				constraint := strings.TrimPrefix(line, dist)
				for _, v := range versions {
					v1, _ := gov.NewVersion(v)
					// Constraints
					constraints, _ := gov.NewConstraint(constraint)
					if constraints.Check(v1) {
						return v1.String(), dir
					}
				}
				constversion := strings.Trim(constraint, "!=<>~")
				return "", fmt.Sprintf("unable to satisfy constraint %q for %q. Try 'binenv install %s %s'.", constraint, dist, dist, constversion)
			}
		}

		if err := scanner.Err(); err != nil {
			return "", ""
		}
		return deflt, "default"
	}
}

func (a *App) getDistributionsFromLock() ([]string, []string) {
	var distributions []string
	var lines []string

	curdir, err := os.Getwd()
	if err != nil {
		a.logger.Error().Err(err).Msg("unable to determine current directory")
		return distributions, lines
	}
	lockfile := filepath.Join(curdir, ".binenv.lock")
	if _, err := os.Stat(lockfile); err != nil {
		// If in stop dir, we found nothing
		a.logger.Error().Err(err).Msg("no .binenv.lock in current directory")
		return distributions, lines
	}

	// lock file is found
	f, err := os.Open(lockfile)
	if err != nil {
		a.logger.Error().Err(err).Msg("unanle to open .binenv.lock")
		return distributions, lines
	}
	defer f.Close()

	seps := "=!<>~"

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Get distribution name
		cpos := strings.IndexAny(line, seps)
		if cpos > 0 {
			distributions = append(distributions, line[0:cpos])
			lines = append(lines, line)
		}
	}
	return distributions, lines
}

func (a *App) loadCache() {
	conf := filepath.Join(a.cachedir, "/cache.json")
	if _, err := os.Stat(conf); os.IsNotExist(err) {
		return
	}

	js, err := ioutil.ReadFile(conf)
	if err != nil {
		a.logger.Error().Err(err).Msgf("unable to read cache %s: please check file permissions", conf)
		return
	}

	err = json.Unmarshal([]byte(js), &a.cache)

	if err != nil {
		a.logger.Error().Err(err).Msgf(`unable to unmarshal cache %s; try to "rm %s && binenv update"`, conf, conf)
		return
	}
}

func (a *App) saveCache() error {
	cache := a.cachedir

	var mode os.FileMode = 0750
	if a.global {
		mode = 0755
	}
	err := os.MkdirAll(cache, mode)
	if err != nil {
		return fmt.Errorf("unable to create cache directory '%s': %w", cache, err)
	}

	cache = filepath.Join(cache, "/cache.json")

	js, err := json.Marshal(&a.cache)
	if err != nil {
		a.logger.Error().Err(err).Msgf("unable to marshal cache %q", cache)
		return nil
	}

	fd, err := os.OpenFile(cache, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
	if err != nil {
		a.logger.Error().Err(err).Msgf("unable to write cache %s: please check file permissions", cache)
		return nil
	}
	defer fd.Close()

	fd.Write(js)

	return nil
}

func (a *App) createInstallers() {
	for k, v := range a.def.Sources {
		i := v.Install.Factory(v.Install.Binaries)
		if i == nil {
			a.logger.Warn().Msgf("%q install method for %q is not implemented", v.Install.Type, k)
			continue
		}
		a.installers[k] = i
	}
}

func (a *App) createMappers() {
	for k, v := range a.def.Sources {
		if v.Map.IsZero() {
			continue
		}
		a.mappers[k] = v.Map
	}
}

func (a *App) createListers() {
	for k, v := range a.def.Sources {
		l := v.List.Factory()
		if l == nil {
			a.logger.Warn().Msgf("%q list method for %q is not implemented", v.List.Type, k)
			continue
		}
		a.listers[k] = l
	}
}

func (a *App) createFetchers() {
	for k, v := range a.def.Sources {
		f := v.Fetch.Factory()
		if f == nil {
			a.logger.Warn().Msgf("%q fetch method for %q is not implemented", v.Fetch.Type, k)
			continue
		}
		a.fetchers[k] = f
	}
}

// Functional options

// WithDiscard sets the log output to /dev/null
func WithDiscard() func(*App) error {
	return func(a *App) error {
		return a.setLogOutput(ioutil.Discard)
	}
}

func (a *App) setLogOutput(w io.Writer) error {
	a.logger = zerolog.Nop()

	return nil
}

// SetBinDir sets bin directory to use the first non-empty argument
func (a *App) SetBinDir(d string) error {
	if d != "" {
		a.bindir = d
		a.logger.Debug().
			Str("bindir", a.bindir).
			Msg("setting configuration")
		return nil
	}

	err := fmt.Errorf("internal error: no bindir has been provided")
	a.logger.Error().Err(err).Msg("setting configuration")
	return err
}

// SetLinkDir sets link directory to use
func (a *App) SetLinkDir(d string) error {
	if d != "" {
		a.linkdir = d
		a.logger.Debug().
			Str("linkdir", a.linkdir).
			Msg("setting configuration")
		return nil
	}

	err := fmt.Errorf("internal error: no linkdir has been provided")
	a.logger.Error().Err(err).Msg("setting configuration")
	return err
}

// SetCacheDir sets cache directory to use
func (a *App) SetCacheDir(d string) error {
	if d != "" {
		a.cachedir = d
		a.logger.Debug().
			Str("cachedir", a.cachedir).
			Msg("setting configuration")
		return nil

	}

	err := fmt.Errorf("internal error: no cachedir has been provided")
	a.logger.Error().Err(err).Msg("setting configuration")
	return err
}

// SetConfigDir sets config directory to use
func (a *App) SetConfigDir(d string) error {
	if d != "" {
		a.configdir = d
		a.logger.Debug().
			Str("configdir", a.configdir).
			Msg("setting configuration")
		return nil
	}

	err := fmt.Errorf("internal error: no configdir has been provided")
	a.logger.Error().Err(err).Msg("setting configuration")
	return err
}

// SetLogLevel sets the log level to use
func (a *App) SetLogLevel(l string) error {
	lvl, err := zerolog.ParseLevel(l)
	if err != nil {
		a.logger.Fatal().Err(err).Msgf("invalid log level %q", l)
	}
	a.logger.Level(lvl)

	return nil
}

// SetVerbose sets the log level to debug
func (a *App) SetVerbose(v bool) {
	if v {
		a.logger = a.logger.Level(zerolog.DebugLevel)
	}
}

// SetDryRun sets the operation mode to dry-run
func (a *App) SetDryRun(v bool) {
	if v {
		a.dryrun = true
	}
}

// SetConcurrency sets the number of goroutines for cache update
func (a *App) SetConcurrency(c int) {
	a.concurrency = c
}

// SetGlobal configures binenv to run in system-wide mode
func (a *App) SetGlobal(g bool) {
	if !g {
		return
	}
	a.bindir = "/var/lib/binenv"
	a.cachedir = "/var/cache/binenv"
	a.configdir = "/var/lib/binenv/config"
	a.linkdir = "/usr/local/bin"
	a.global = true

	a.loadCache()
}

// DumpConfig dumps the configuration to stdout
func (a *App) DumpConfig() {
	a.logger.Debug().
		Str("bindir", a.bindir).
		Str("linkdir", a.bindir).
		Str("cachedir", a.cachedir).
		Str("configdir", a.configdir).
		Msg("final configuration")
}

// initializeFlags sets flags holder
func (a *App) initializeFlags() {
	f := flags{}
	f.Setters = make(map[string]func(*flags, bool))

	// Default behavior here
	f.justExpand = false
	f.Setters["justExpand"] = func(obj *flags, b bool) { obj.justExpand = b }

	// bind defaults flags
	a.flags = f
}

// SetFlag call a setter function if exists
func (a *App) SetFlag(name string, value bool) {
	for key, setter := range a.flags.Setters {
		if key == name {
			setter(&a.flags, value)
		}
	}
}
