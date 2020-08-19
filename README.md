# binenv

The last binary you'll ever install.

## What

`binenv` will help you download, install and manage the binaries programs
(a.k.a. distributions) you need in you everyday DevOps life (e.g. kubectl,
helm, ...).

Think of it as a `tfenv` + `tgenv` + `helmenv` + ...

## Quick start

### Linux (bash/zsh)

```
wget -q https://github.com/devops-works/binenv/releases/latest/download/binenv-linux-amd64
mv binenv-linux-amd64 binenv
chmod +x binenv
./binenv update
./binenv install binenv 0.0.5
rm binenv
if [[ -n $BASH ]]; then ZESHELL=bash; fi
if [[ -n $ZSH_NAME ]]; then ZESHELL=zsh; fi
echo $ZESHELL
echo 'export PATH=~/.binenv:$PATH' >> ~/.${ZESHELL}rc
echo "source <(binenv completion ${ZESHELL})" >> ~/.${ZESHELL}rc
exec $SHELL
```

### MacOS (with bash)

```
wget -q https://github.com/devops-works/binenv/releases/latest/download/binenv-darwin-amd64
mv binenv-linux-darwin binenv
chmod +x binenv
./binenv update
./binenv install binenv 
rm binenv
echo 'export PATH=~/.binenv:$PATH' >> ~/.bashrc
echo 'source <(binenv completion bash)' >> ~/.bashrc
exec $SHELL
```

### Windows

TBD

## Install

- download a suitable `binenv` (yes, but wait !) for your architecture/OS at
http://github.com/devops-works/binenv/releases.

```
wget -q https://github.com/devops-works/binenv/releases/latest/download/binenv-<OS>-<ARCH>
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
[![asciicast](https://asciinema.org/a/2SY7RRnQahtkfNSiX99gxzZLW.svg)](https://asciinema.org/a/2SY7RRnQahtkfNSiX99gxzZLW)

## Updating binenv

Just run `binenv install binenv`

This is the whole point.

## Supported "distributions"

"Distributions" are installable binaries. We just had to find a name ¯\\_(ツ)_/¯.

Currently supported distributions are:

- awless
- binenv
- cli53
- consul
- docker-compose
- hadolint
- helm
- helmfile
- hugo 
- k9s
- krew
- kops
- kube-bench
- kubectl
- kubectx
- kubens
- rancher[^1] (soon [#11](https://github.com/devops-works/binenv/issues/11))
- terraform
- terragrunt
- vault

The always up-to-date list is
[here](https://github.com/devops-works/binenv/blob/master/distributions/distributions.yaml).

Open an issue (or send a PR) if you need one that is not in the list.

[^1]: cli for rancher 1.x

## Usage

### Updating available distributions

In order to update the list of installable version for distributions, you need
to update the version list (usually located in `$XDG_CONFIG/cache.json` or
`~/.config/binenv/cache.json`).

This is done automatically when invoking `binenv update`.

Without arguments, il will check for available versions for _all_ distributions
(watch out for Github API rate limits).

With a distribution passed as an argument (e.g. `binenv update kubectl`), it
will only update installable versions for `kubectl`.

Note that Github enforces rate limits (e.g. 60 unauthenticated API requests per
hours). So you should update all distributions (e.g. `binenv update`) with
caution.

There is currently no support for tokens.

#### Examples

- `binenv update`: update available versions for all distributions
- `binenv update kubectl helm`: update available versions for `kubectl` and
  `helm`

### Installing new versions

After updating the list, you might want to install a shiny new version. No
 problem,`binenv install` has you covered.
 
If you want the latest non-prerelease version for something, just run:

`binenv install something`

If you want a specific version:

`binenv install something 1.2.3`

Note that completion works, so don't be afraid to use it.

You can also install serveral distribution versions at the same time:

`binenv install something 1.2.3 somethingelse 4.5.6`

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

In the output, versions having a `*` next to them are the currenly selected
versions (see [Selecting versions](#selecting-versions) below.

Versions having a `+` next to them are installed. 

All other versions are available to be installed.

#### Examples

```
$ binenv versions
kubectl:
        ...
        1.19.0-alpha.3
        1.19.0-alpha.2
        1.18.8* (from default)
        1.18.6
        1.18.5
        1.18.4
        1.18.3
        1.18.2
        1.17.11
        1.17.9+
        1.17.8
        1.17.8-rc.1
        1.17.7
        1.17.6
        1.17.5+
        1.16.14
        1.16.13
        1.16.12
        1.16.12-rc.1

terraform:
        0.13.0
        0.13.0-rc1
        0.13.0-beta3
...
```

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

### Update available distributions

Distributions are mainained in this
[repo](https://github.com/devops-works/binenv/blob/master/distributions/distributions.yaml).
To benefit from new additions, you need to update the distribution list from
time to time.

This list is usually located in your home directory under
`$XDG_CONFIG/distributions.yaml` or `~/.config/binenv/distribution.yaml`).

Usage is `binenv distributions`.

### Completion

Install completion for your shell. See `binenv help completion` for in-depth
info.

## Selecting versions

To specify which version to use, you have to create a `.binenv.lock` file in
the directory. Note that only **semver** is supported.

This file has the follosing structure:

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
- `~>`: version must be at least this one in the same but match teh same minor
  versions

### Version selection process

When you execute a distribution (e.g. you run `kubectl`), `binenv` runs it
under the hood. Before running it, it will check which version it should use.
For this, it will check for a `.binenv.lock` file in the current directory.

If none is found, it will check in the parent folder. No lock file ? Check in
parent folder again. this process continues until `binenv` reaches your home
directory.

If no version requirements are found at this point, `binenv` will use the last
non-prerelease version installed.

## Removing binenv stuff

`binenv` stores downloaded binaries in `~/.binenv/binaries`, and a cache in
`~/.config/binenv/` (or whatever your `XDG_CONFIG` variables points to).

To wipe everything clean:

```bash
rm -rf ~/.binenv ~/.config/binenv/
```

Don't forget to remove the `PATH` and the completion you might have changed in
your shell rc file.

## Status

This is really _super alpha_ and has only be tested on Linux & MacOS. YMMV on
other platforms.

There are **no tests**. I will probably go to hell for this.

## Contributions

Welcomed !

We will need other instalaltion mechanisms (see
https://github.com/devops-works/binenv/tree/master/internal/install).

## Licence

MIT
