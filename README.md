# binenv

The last binary you'll ever install.

- [binenv](#binenv)
  - [What](#what)
  - [Quick start](#quick-start)
    - [Linux (bash/zsh)](#linux-bashzsh)
    - [MacOS (with bash)](#macos-with-bash)
    - [Windows](#windows)
    - [FreeBSD (bash/zsh)](#freebsd-bashzsh)
    - [OpenBSD (bash/zsh)](#openbsd-bashzsh)
  - [Install](#install)
    - [User install](#user-install)
  - [Updating binenv](#updating-binenv)
  - [Supported "distributions"](#supported-distributions)
  - [Usage](#usage)
    - [Updating available distributions versions](#updating-available-distributions-versions)
      - [Updating versions using a token](#updating-versions-using-a-token)
      - [Update available distributions](#update-available-distributions)
      - [Examples](#examples)
    - [Searching distributions](#searching-distributions)
    - [Installing new versions](#installing-new-versions)
      - [Examples](#examples-1)
    - [Listing versions](#listing-versions)
      - [Examples](#examples-2)
      - [Freezing versions](#freezing-versions)
    - [Uninstalling versions](#uninstalling-versions)
      - [Examples](#examples-3)
    - [Completion](#completion)
    - [Expanding binary absolute path](#expanding-binary-absolute-path)
      - [Example](#example)
    - [Upgrading all installed distributions](#upgrading-all-installed-distributions)
  - [Selecting versions](#selecting-versions)
    - [Version selection process](#version-selection-process)
    - [Install versions form .binenv.lock](#install-versions-form-binenvlock)
      - [Example](#example-1)
    - [Adding versions to `.binenv.lock`](#adding-versions-to-binenvlock)
  - [Environment variables](#environment-variables)
  - [Removing binenv stuff](#removing-binenv-stuff)
  - [Status](#status)
  - [FAQ](#faq)
    - [I installed a binary but is still see the system (or wrong) version](#i-installed-a-binary-but-is-still-see-the-system-or-wrong-version)
    - [After installing a distribution, I get a "shim: no such file or directory"](#after-installing-a-distribution-i-get-a-shim-no-such-file-or-directory)
    - [It does not work with sudo !](#it-does-not-work-with-sudo-)
    - [I don't like binenv, are there alternatives ?](#i-dont-like-binenv-are-there-alternatives-)
  - [Distributions file format](#distributions-file-format)
    - [Distributions file reference](#distributions-file-reference)
    - [Distributions file example](#distributions-file-example)
  - [Caveats](#caveats)
  - [Contributions](#contributions)
  - [Licence](#licence)

## What

`binenv` will help you download, install and manage the binaries programs (we
call them "distributions") you need in you everyday DevOps life (e.g. kubectl,
helm, ...).

Think of it as a `tfenv` + `tgenv` + `helmenv` + ...

Now you can install your [favorite utility](#supported-distributions) just by
typing `binenv install something`.

## Quick start

See [System-wide installation](./SYSTEM.md) for system-wide installations
(a.k.a. global mode).

### Linux (bash/zsh)

```
wget -q https://github.com/devops-works/binenv/releases/download/v0.19.11/binenv_linux_amd64
wget -q https://github.com/devops-works/binenv/releases/download/v0.19.11/checksums.txt
sha256sum  --check --ignore-missing checksums.txt
mv binenv_linux_amd64 binenv
chmod +x binenv
./binenv update
./binenv install binenv
rm binenv
if [[ -n $BASH ]]; then ZESHELL=bash; fi
if [[ -n $ZSH_NAME ]]; then ZESHELL=zsh; fi
echo $ZESHELL
echo -e '\nexport PATH=~/.binenv:$PATH' >> ~/.${ZESHELL}rc
echo "source <(binenv completion ${ZESHELL})" >> ~/.${ZESHELL}rc
exec $SHELL
```

### MacOS (with bash)

```
wget -q https://github.com/devops-works/binenv/releases/download/v0.19.11/binenv_darwin_amd64
wget -q https://github.com/devops-works/binenv/releases/download/v0.19.11/checksums.txt
sha256sum  --check --ignore-missing checksums.txt
mv binenv_darwin_amd64 binenv
chmod +x binenv
./binenv update
./binenv install binenv
rm binenv
echo -e '\nexport PATH=~/.binenv:$PATH' >> ~/.bashrc
echo 'source <(binenv completion bash)' >> ~/.bashrc
exec $SHELL
```

### Windows

binenv does not support windows.

### FreeBSD (bash/zsh)

```
fetch https://github.com/devops-works/binenv/releases/download/v0.19.11/binenv_freebsd_amd64
fetch https://github.com/devops-works/binenv/releases/download/v0.19.11/checksums.txt
shasum --ignore-missing -a 512 -c checksums.txt
mv binenv_freebsd_amd64 binenv
chmod +x binenv
./binenv update
./binenv install binenv
rm binenv
if [[ -n $BASH ]]; then ZESHELL=bash; fi
if [[ -n $ZSH_NAME ]]; then ZESHELL=zsh; fi
echo $ZESHELL
echo -e '\nexport PATH=~/.binenv:$PATH' >> ~/.${ZESHELL}rc
echo "source <(binenv completion ${ZESHELL})" >> ~/.${ZESHELL}rc
exec $SHELL
```

If you are using a different shell, skip adding completion to your `.${SHELL}rc` file.

To be able to verify checksums, you have to install the `p5-Digest-SHA` package.

### OpenBSD (bash/zsh)

```
ftp https://github.com/devops-works/binenv/releases/download/v0.19.11/binenv_openbsd_amd64
ftp https://github.com/devops-works/binenv/releases/download/v0.19.11/checksums.txt
cksum -a sha256 -C checksums.txt binenv_openbsd_amd64
mv binenv_openbsd_amd64 binenv
chmod +x binenv
./binenv update
./binenv install binenv
rm binenv
if [[ -n $BASH ]]; then ZESHELL=bash; fi
if [[ -n $ZSH_NAME ]]; then ZESHELL=zsh; fi
echo $ZESHELL
echo -e '\nexport PATH=~/.binenv:$PATH' >> ~/.${ZESHELL}rc
echo "source <(binenv completion ${ZESHELL})" >> ~/.${ZESHELL}rc
exec $SHELL
```

If you are using a different shell, skip adding completion to your `.${SHELL}rc` file.

## Install

### User install

- download a suitable `binenv` (yes, but wait !) for your architecture/OS at
http://github.com/devops-works/binenv/releases.

```
wget -q https://github.com/devops-works/binenv/releases/download/v0.19.11/binenv_<OS>_<ARCH>
```

- rename it

```
mv binaryname binenv
```

- make it executable

```
chmod +x binenv
```

- execute an update

```
./binenv update
```

- now install `binenv` with `binenv` (so meta)

```
./binenv install binenv <version>
```

- you can now remove the downloaded file

```
rm binenv
```

- prepend `~/.binenv` to your path in your `~/.bashrc` or `~/.zshrc` or ...

```
export PATH=~/.binenv:$PATH
```

- while you're at it, install the completion (replace `bash` with your shell)

```
source <(binenv completion bash)
```

- "restart" your shell

```
exec $SHELL
```

See a walkthough on asciinema.org:
[![asciicast](https://asciinema.org/a/LmYClC9sVgNs24QZjKIFccFh4.svg)](https://asciinema.org/a/LmYClC9sVgNs24QZjKIFccFh4)

## Updating binenv

Just run `binenv install binenv`

This is the whole point.

## Supported "distributions"

For the whole list of supported binaries (a.k.a. distributions), see
[DISTRIBUTIONS.md](DISTRIBUTIONS.md).

The always up-to-date list is
[here](https://github.com/devops-works/binenv/blob/master/distributions/distributions.yaml).

The list can be generated as markdown using `make distributions`.

Open an issue (or send a PR) if you need one that is not in the list.

## Usage

### Updating available distributions versions

In order to update the list of installable version for distributions, you need
to update the version list (usually located in `$XDG_CONFIG/cache.json` or
`~/.config/binenv/cache.json`).

This is done automatically when invoking `binenv update`.

Without arguments, it will fetch the cache from this repo. This cache is
generated automatically daily.

Using the `-f` argument, `binenv` will retrieve available versions for _all_
distributions (watch out for Github API rate limits, [but see
below](#updating-versions-from-generated-cache)).

With a distribution passed as an argument (e.g. `binenv update kubectl`), it
will only update installable versions for `kubectl`.

When updating the cache, you can control fetch concurrency using the `-c` flag.
It defaults to 8 which is already pretty high. Do go crazy. This setting is
mainly used to set a lower concurrency and be nice to GitHub.

Note that Github enforces rate limits (e.g. 60 unauthenticated API requests per
hours). So you should update all distributions (e.g. `binenv update -f`) with
caution. `binenv` will stop updating distributions when you only have 4
unauthenticated API requests left.

[GitHub tokens](#updating-versions-using-a-token) are also supported to avoid
being rate-limited and fetch releases from their respective sources.

#### Updating versions using a token

To avoid being rate limited, you can also use a personal access token.

- go to [Settings/Personal Access Tokens/New personal access token](https://github.com/settings/tokens/new?description=GitHunt&scopes=public_repo)
- click "Generate token"

To use the token, just export it in the GITHUB_TOKEN environment variable:

```bash
export GITHUB_TOKEN=aaa...bbb
```

#### Update available distributions

Distributions are maintained in this
[file](https://github.com/devops-works/binenv/blob/master/distributions/distributions.yaml).

To benefit from new additions, you need to update the distribution list from
time to time.

This list is usually located in your home directory under
`$XDG_CONFIG/distributions.yaml` (often `~/.config/binenv/distribution.yaml`).

To update only distributions:

```bash
binenv update --distributions # or -d
```

To update distributions **and** their versions:

```bash
binenv update --all # or -a
```

##### Using custom distributions file (and private GitLab repos)

If you want to use a custom distributions file, you can add a `.yaml` file in
the `$XDG_CONFIG` directory (often `~/.config/binenv/`).

This file will be merged with the default distributions file.

Note that files are evaluated in lexicographical order, so if you want to
override a default, you should name your file accordingly.

You can use this mechanism to install binaries from private GitLab repositories
(GitHub not supported right now). If you need to pass a `PRIVATE-TOKEN` in the
headers, you need to set the `token_env` key in the `list` and `fetch`
sections. This key should contain the name of the environment variable that is
set with the token.

Here is an example file:

```yaml
$ cat ~/.config/binenv/distributions-custom.yaml
---
sources:
  foo:
    description: This tool let's you foo database tables
    url: https://gitlab.exemple.org/infrastructure/tools/foo
    list:
      type: gitlab-releases
      url: https://gitlab.example.org/api/v4/projects/42/releases
      token_env: FOO_PRIVATE_TOKEN
    fetch:
      url: https://gitlab.example.org/api/v4/projects/42/packages/generic/foo/{{ .Version }}/foo-{{.OS }}-{{ .Arch }}-{{ .Version }}.gz
      token_env: FOO_PRIVATE_TOKEN
    install:
      type: gzip
      binaries:
        - "foo-{{.OS }}-{{ .Arch }}-{{ .Version }}.gz"
```

You will have to `export FOO_PRIVATE_TOKEN=your_token` before running `binenv`
to make the token available.

#### Examples

- `binenv update`: update available versions for all distributions from github
  cache
- `binenv update -f`: update available versions for all distributions from all
  releases
- `binenv update -d`: update available distributions
- `binenv update kubectl helm`: update available versions for `kubectl` and
  `helm`

### Searching distributions

The `search` command lets you search a distribution by name or description:

```bash
$ binenv search kube
binenv: One binary to rule them all. Manage all those pesky binaries (kubectl, helm, terraform, ...) easily.
helm: The Kubernetes Package Manager
helmfile: Deploy Kubernetes Helm Charts
k9s: Kubernetes CLI To Manage Your Clusters In Style!
ketall: Like `kubectl get all`, but get really all resources
... (lots of things with "kube" in it)
```

### Installing new versions

After updating the list, you might want to install a shiny new version. No
 problem,`binenv install` has you covered.

If you want the latest non-prerelease version for something, just run:

`binenv install something`

If you want a specific version:

`binenv install something 1.2.3`

Note that completion works, so don't be afraid to use it.

You can also install several distribution versions at the same time:

`binenv install something 1.2.3 somethingelse 4.5.6`

Using the `--dry-run` flag (a.k.a `-n`) will show what would be installed.

#### Examples

- `binenv install kubectl`: install latest non-prerelease `kubectl version`
- `binenv install kubectl 1.18.8`: install `kubectl` version 1.18.8
- `binenv install kubectl 1.18.8 helm 3.3.0`: install `kubectl` version 1.18.8
  and `helm` 3.3.0

### Listing versions

You can list available, installed and activated distribution versions using
`binenv versions`.

When invoked without arguments, all version of all distributions will be printed.

With distributions as arguments, only versions for those distributions will be
printed.

In the output, versions printed in reverse mode are the currently selected
(a.k.a. active) versions (see [Selecting versions](#selecting-versions) below.

Versions in **bold** are installed.

All other versions are available to be installed.

#### Examples

```
$ binenv versions
terraform: 0.13.1 (/home/you/some/dir) 0.13.0 0.13.0-rc1 0.13.0-beta3 0.13.0-beta2 0.13.0-beta1 0.12.29 0.12.28 0.12.27 0.12.26 0.12.25 0.12.24 0.12.23 0.12.22 0.12.21 0.12.20 0.12.19 0.12.18 0.12.17 0.12.16 0.12.15 0.12.14 0.12.13 0.12.12 0.12.11 0.12.10 0.12.9 0.12.8 0.12.7 0.12.6
terragrunt: 0.23.38 0.23.37 0.23.36 0.23.35 0.23.34 0.23.33 0.23.32 0.23.31 0.23.30 0.23.29 0.23.28 0.23.27 0.23.26 0.23.25 0.23.24 0.23.23 0.23.22 0.23.21 0.23.20 0.23.19 0.23.18 0.23.17 0.23.16 0.23.15 0.23.14 0.23.13 0.23.12 0.23.11 0.23.10 0.23.9
toji: 0.2.4 (default) 0.2.2
vault: 1.5.3 1.5.2 1.5.1 1.5.0 1.5.0-rc 1.4.6 1.4.5 1.4.4 1.4.3 1.4.2 1.4.1 1.4.0 1.4.0-rc1 1.4.0-beta1 1.3.10 1.3.9 1.3.8 1.3.7 1.3.6 1.3.5 1.3.4 1.3.3 1.3.2 1.3.1 1.3.0 1.3.0-beta1 1.2.7 1.2.6 1.2.5 1.2.4
...
```

(the output above does not show bold or reverse terminal output)

#### Freezing versions

When the `versions` command is invoked with the `--freeze` option, it will
write a `.binenv.lock` style file on stdout.

This way you can "lock" the dependencies for your project just by issuing:

```
cd myproject
binenv versions --freeze > .binenv.lock
```

You can the commit this file to your project so everyone will use the same
distributions versions when in this repository. See [Selecting Versions](#selecting-versions) for more information on this file.

Note that currently selected versions for _all_ distributions will be
outputted. You might want to trim stuff you do not use from the file.

### Uninstalling versions

If you need to clean up a bit, you can uninstall a specific version, or all
versions for a distribution. In the latter case, a confirmation will be asked.

The command accepts:
- a single argument (remove all versions for distributions)
- an even count of arguments (distribution / version pairs)

#### Examples

- `binenv uninstall kubectl 1.18.8 helm 3.3.0`: uninstall `kubectl` version
  1.18.8 and `helm` 3.3.0
- `binenv uninstall kubectl 1.18.8 kubectl 1.16.15`: uninstall `kubectl` versions
  1.18.8 and 1.16.15
- `binenv uninstall kubectl`: removes all `kubectl` versions

### Completion

Install completion for your shell. See `binenv help completion` for in-depth
info.

### Expanding binary absolute path

To get the absolute path of the binary installed by a distribution you need to
invoke the command `expand`.

This can be useful when you need to use binenv in conjunction with other tools
like `sudo`.

#### Example

```bash
$ binenv install yq
2022-02-16T14:24:56-03:00 WRN version for "yq" not specified; using "4.18.1"
fetching yq version 4.18.1 100% |████████████████████████████| (9.1/9.1 MB, 4.858 MB/s)
2022-02-16T14:24:59-03:00 INF "yq" (4.18.1) installed
$ binenv expand yq
/Users/local-user/.binenv/binaries/yq/4.18.1
$ sudo $(binenv expand yq) --version
yq (https://github.com/mikefarah/yq/) version 4.18.1
```

### Upgrading all installed distributions

To upgrade all installed distributions to the last known version invoke the
command `upgrade`

This command will always select the last version available and will ignore any
version selection previously made by the user.

## Selecting versions

To specify which version to use, you have to create a `.binenv.lock` file in
the directory. Note that only **semver** is supported.

This file has the following structure:

```
<distributionA><constraintA>
<distributionB><constraintB>
...
```

For instance:

```
kubectl=1.18.8
terraform>0.12
terragrunt~>0.23.0
```

You can then commit the file in your project to ensure everyone in your team is
on the same page.

The constraint operators are:

- `=`:  version must match exactly
- `!=`: version must not match
- `>`:  version must be strictly higher
- `<`:  version must be strictly lower
- `>=`: version must be at least
- `<=`: version must be at most
- `~>`: version must be at least this one in the same but match the same minor
  versions

### Version selection process

When you execute a distribution (e.g. you run `kubectl`), `binenv` runs it
under the hood. Before running it, it will check which version it should use.
For this, it will check for a `.binenv.lock` file in the current directory.

If none is found, it will check in the parent folder. No lock file ? Check in
parent folder again. this process continues until `binenv` reaches your home
directory (or `/` if run in global mode).

If no version requirements are found at this point, `binenv` will use the last
non-prerelease version installed.

### Install versions form .binenv.lock

Install versions specified in `.binenv.lock` file, you can use the `--lock`
(a.k.a. `-l`) flag.

#### Example

```bash
$ cat .binenv.lock
terraform>0.13.0
helmfile<0.125.0
hadolint<1.17.0
$ binenv install -l
2020-08-29T11:39:18+02:00 WRN installing "terraform" (0.13.1) to satisfy constraint "terraform>0.13.0"
fetching terraform version 0.13.1 100% |█████████████████████████████████████████████████████████████████████████████████████████████| (33/33 MB, 3.274 MB/s) [10s:0s]
2020-08-29T11:39:29+02:00 WRN installing "helmfile" (0.124.0) to satisfy constraint "helmfile<0.125.0"
fetching helmfile version 0.124.0 100% |█████████████████████████████████████████████████████████████████████████████████████████████| (45/45 MB, 1.404 MB/s) [31s:0s]
2020-08-29T11:40:02+02:00 WRN installing "hadolint" (1.16.3) to satisfy constraint "hadolint<1.17.0"
fetching hadolint version 1.16.3 100% |███████████████████████████████████████████████████████████████████████████████████████████| (3.5/3.5 MB, 431.886 kB/s) [8s:0s]
$
```

### Adding versions to `.binenv.lock`

To populate the `.binenv.lock` file in the current directory, you can use the
`local` command with the distributions and versions you want to add.

For instance:

```bash
binenv local kubectl 1.30.0 helmfile 0.126.0
```

Note that this will update the `.binenv.lock` file and not replace it, so the
command above is equivalent to:

```bash
binenv local kubectl 1.30.0
binenv local helmfile 0.126.0
```

and produce the following `.binenv.lock` file:

```bash
kubectl=1.30.0


### Selecting versions using environment variables

_Introduced in v0.17.0_

In addition to using the .binenv.lock file, it is possible to define the
distribution version using an environment variable of the form
`BINENV_<DISTRIBUTION>_VERSION=<CONSTRAINT>`.

When an environment variable with this name exists, binenv will use the `=`
operator to look for an exact match for that constraint and will ignore the
contents of the `.binenv.lock` file if it exists.

#### Example

```bash
$ cat .binenv.lock
helm=3.7.2

$ helm version
version.BuildInfo{Version:"v3.7.2", GitCommit:"663a896f4a815053445eec4153677ddc24a0a361", GitTreeState:"clean", GoVersion:"go1.16.10"}

$ BINENV_HELM_VERSION=3.6.3 helm version
version.BuildInfo{Version:"v3.6.3", GitCommit:"d506314abfb5d21419df8c7e7e68012379db2354", GitTreeState:"clean", GoVersion:"go1.16.5"}
```

## Environment variables

Other environment variables exists to control `binenv` behavior:

- `BINENV_GLOBAL`: forces `binenv` to run un global mode (same as `-g`); see
  [SYSTEM.md](./SYSTEM.md) for more information on this mode.
- `BINENV_VERBOSE`: same as `-v`
- `BASH_COMP_DEBUG_FILE`: if set, will write debug information for bash
  completion to this file

## Removing binenv stuff

`binenv` stores

- downloaded binaries by default in `~/.binenv/binaries`
- the versions cache in `~/.cache/binenv/` (or wherever your `XDG_CACHE_HOME` variable points to)
- the list of known distributions in `~/.config/binenv/` (or wherever your `XDG_CONFIG_HOME` variable points to).

To wipe everything clean:

```bash
rm -rfi ~/.binenv ~/.config/binenv ~/.cache/binenv
```

Don't forget to remove the `PATH` and the completion you might have changed in
your shell rc file.

## Status

This is really _super alpha_ and has only be tested on Linux & MacOS. YMMV on
other platforms.

There are **no tests**. I will probably go to hell for this.

## FAQ

### I installed a binary but is still see the system (or wrong) version

Try to rehash your binaries (`hash -r` in bash or `rehash` in Zsh).

### After installing a distribution, I get a "shim: no such file or directory"

If you see something like:

```
2020-11-10T09:01:20+01:00 ERR unable to install "kubectl" (1.19.3) error="unable to find shim file: stat /Users/foo/.binenv/shim: no such file or directory"
```

you probably did not follow the [installation instructions](#quick-start).

Running `./binenv update binenv && ./binenv install binenv` should correct the
problem.

### It does not work with sudo !

Yes, for not we'restuckon this one. You still can reference thereal binary
directly:

```bash
sudo ~/.binenv/binaries/termshark/2.2.0
```

### I don't like binenv, are there alternatives ?

Sorry to hear that. Don't hesitate opening an issue or sending a PR is
something does not fit your use case

A nice alternative exists:

- https://asdf-vm.com/

## Distributions file format

[distributions.yaml](https://github.com/devops-works/binenv/blob/develop/distributions/distributions.yaml)
contains all the distributions supported by `binenv`, and how to fetch them. It
is written in YAML and is defined by the scheme below.

### Distributions file reference

```yaml
sources:

  # Name of the distribution
  <string>:

    # Description provided by the binary author(s).
    description: <string>

    # URL for binary (usually homepage or repository).
    url: <url>

    # Post install message shown after successful installation
    # Use `post_install_message: |` for multi-line messages
    post_install_message: <string>

    # map creates aliases between architectures known by binenv and those
    # expected by the original author(s).
    # Check `bat` distribution for a more meaningful example.
    [map: <map_config>]

    # list contains the kind of releases and where to fetch their
    # history.
    list:

      # Type of the releases.
      # One of "static", "github-releases", "gitlab-releases"
      type: <string>

      # Where to fetch the releases.
      # I.e. https://github.com/devops-works/binenv/releases
      url: <string>

    # fetch holds the URL from where the binaries can be downloaded.
    fetch:

      # Templatised URL to the binary. Values to templatise can be:
      # Host architecture with {{ .Arch }}, operating system with {{ .OS }},
      # version with {{ .Version }}, sometimes .exe with {{ .ExeExtension}}.
      url: <string>

    # Defines how to install the binary.
    install:

      # Type of installation. Can be :
      # "direct" if after download the binary is executable as is;
      # "tgz" if it needs to be uncompressed using tar and gzip;
      # "zip" if it needs to be unzipped;
      # "tarx" if it needs to be uncompressed with tar;
      type: <string>

      # Name of the binar(y|ies) that will be downloaded
      [binaries: <binaries_config>]

    # Supported platforms
    [supported_platforms: <supported_platforms>]
```

`map_config`:

```yaml
# Alias to amd64 arch
[amd64: <string>]

# Alias to i386 arch
[i386: <string>]

# Alias to darwin arch
[darwin: <string>]

# Alias to linux arch
[linux: <string>]

# Alias to windows arch
[windows: <string>]
```

`binaries_config`:

```yaml
# Array of binaries names that will be installed.
# The string provided is treated as a regexp.
# This regexp is compared to the filenames found in packages.
# Note that filenames contains their path in the package with the top level
# directory removed, e.g.:
# software-13.0.0-x86_64-unknown-linux-musl/foo/bar/zebinary
# becomes
# foo/bar/zebinary
# Also note that, since all binaries will be installed as the distribution
# entry name, only one (the latest match) will survive for now.
# The list is just here to allow alternate names, not real multiple binaries
# installation.
 - <regexp>
```

`supported_platforms`:

```yaml
# Array of OS/Arch parirs
# See: https://pkg.go.dev/runtime#pkg-constants
- os: <string>
  arch: <string>
```

### Distributions file example

```yaml
sources:
  popeye:
    description: A Kubernetes cluster resource sanitizer
    url: https://github.com/derailed/popeye
    map:
      amd64: x86_64
      darwin: Darwin
      linux: Linux
      windows: Windows
    list:
      type: github-releases
      url: https://api.github.com/repos/derailed/popeye/releases
    fetch:
      url: https://github.com/derailed/popeye/releases/download/v{{ .Version }/popeye_{{ .OS }}_{{ .Arch }}.tar.gz
    install:
      type: tgz
      binaries:
        - popeye
    supported_platforms:
      - os: linux
        arch: amd64
      - os: windows
        arch: amd64
      - os: darwin
        arch: amd64
```

The `distributions.yaml` file used by default by `binenv` is located [here](https://github.com/devops-works/binenv/blob/develop/distributions/distributions.yaml), don't hesitate to have a look on it's structure.

## Caveats

Since `binenv` uses your PATH and HOME to find binaries and layout it's
configuration files, using sudo with binenv-installed binaries is not very
straightforward. You can either install binenv as the root user (so it can find
it's config), or pass those two environment variables when invoking sudo, like
so:

```
sudo env "PATH=$PATH" "HOME=$HOME" binary_installed_with_binenv ...
```

## Contributions

Welcomed !

Thanks to all contributors:

- @alexanderbe-commit
- @alenzen
- @alex-bes
- @angrox
- @axgkl
- @cleming
- @Dazix
- @deknos
- @DnR-iData
- @dundee
- @eagafonov
- @earzur
- @eze-kiel
- @gwenall
- @harleypig
- @iainelder
- @jakubvokoun
- @kenni-shin
- @mpepping
- @patsevanton
- @pichouk
- @pklejch
- @semoac
- @shr-project
- @Sierra1011
- @tm-drtina
- @xx4h

## Licence

MIT
