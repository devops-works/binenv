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
wget -q https://github.com/devops-works/binenv/releases/latest/download/binenv_linux_amd64 -O binenv
chmod +x binenv
./binenv update
./binenv install binenv
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
wget -q https://github.com/devops-works/binenv/releases/latest/download/binenv_darwin_amd64 -O binenv
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
[![asciicast](https://asciinema.org/a/LmYClC9sVgNs24QZjKIFccFh4.svg)](https://asciinema.org/a/LmYClC9sVgNs24QZjKIFccFh4)

## Updating binenv

Just run `binenv install binenv`

This is the whole point.

## Supported "distributions"

"Distributions" are installable binaries. We just had to find a name ¯\\_(ツ)_/¯.

Currently supported distributions are:

- [ali](https://github.com/nakabonne/ali)
- [asciinema-edit](https://github.com/cirocosta/asciinema-edit)
- [awless](https://github.com/wallix/awless)
- [awstaghelper](https://github.com/mpostument/awstaghelper)
- [binenv](https://github.com/devops-works/binenv)
- [cli53](https://github.com/barnybug/cli53)
- [consul](https://www.consul.io/)
- [dive](https://github.com/wagoodman/dive/)
- [docker-compose](https://docs.docker.com/compose/)
- [duf](https://github.com/muesli/duf)
- [dw-query-digest](https://github.com/devops-works/dw-query-digest)
- [gitjacker](https://github.com/liamg/gitjacker)
- [goreleaser](https://github.com/goreleaser/goreleaser)
- [grizzly](https://github.com/grafana/grizzly)
- [gws](https://github.com/StreakyCobra/gws)
- [hadolint](https://github.com/hadolint/hadolint)
- [helm](https://helm.sh/)
- [helmfile](https://github.com/roboll/helmfile)
- [hey](https://github.com/rakyll/hey)
- [hugo](https://gohugo.io/)
- [k9s](https://k9scli.io/)
- [ketall](https://github.com/corneliusweig/ketall)
- [kind](https://github.com/kubernetes-sigs/kind)
- [krew](https://github.com/kubernetes-sigs/krew)
- [kops](https://kops.sigs.k8s.io/)
- [kube-bench](https://github.com/aquasecurity/kube-bench)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/)
- [kubectx](https://github.com/ahmetb/kubectx)
- [kubens](https://github.com/ahmetb/kubectx)
- [kustomize](https://github.com/kubernetes-sigs/kustomize)
- [logcli](https://github.com/grafana/loki/)
- [loki](https://github.com/grafana/loki/)
- [minikube](https://github.com/kubernetes/minikube)
- [packer](https://www.packer.io)
- [pluto](https://github.com/FairwindsOps/pluto)
- [pomerium](https://github.com/pomerium/pomerium)
- [pomerium-cli](https://github.com/pomerium/pomerium)
- [popeye](https://github.com/derailed/popeye)
- [promtail](https://github.com/grafana/loki/)
- [rancher](https://rancher.com/docs/rancher/v1.6/en/)[^1]
- [tanka](https://github.com/grafana/tanka)
- [terraform](https://www.hashicorp.com/products/terraform)
- [terragrunt](https://terragrunt.gruntwork.io/)
- [tflint](https://github.com/terraform-linters/tflint/)
- [toji](https://github.com/leucos/toji/)
- [trivy](https://github.com/aquasecurity/trivy)
- [vault](https://www.hashicorp.com/products/vault)
- [vmctl](https://github.com/VictoriaMetrics/vmctl)
- [yh](https://github.com/andreazorzetto/yh)

The always up-to-date list is
[here](https://github.com/devops-works/binenv/blob/master/distributions/distributions.yaml).

Open an issue (or send a PR) if you need one that is not in the list.

[^1]: cli for rancher 1.x

## Usage

### Updating available distributions versions

In order to update the list of installable version for distributions, you need
to update the version list (usually located in `$XDG_CONFIG/cache.json` or
`~/.config/binenv/cache.json`).

This is done automatically when invoking `binenv update`.

Without arguments, il will check for available versions for _all_ distributions
(watch out for Github API rate limits, [but see below](#updating-versions-from-generated-cache)).

With a distribution passed as an argument (e.g. `binenv update kubectl`), it
will only update installable versions for `kubectl`.

Note that Github enforces rate limits (e.g. 60 unauthenticated API requests per
hours). So you should update all distributions (e.g. `binenv update`) with
caution.

`binenv` will stop updating distributions when you only have 4 unauthenticated
API requests left.

There is currently no support for Github tokens.

#### Updating versions from generated cache

To avoid hitting the Github unauthenticated rate limit, `binenv` can fetch a nightly generated versions cache from github using the `--cache` flag:

```bash
binenv update --cache # or -c
```

#### Update available distributions

Distributions are maintained in this
[file](https://github.com/devops-works/binenv/blob/master/distributions/distributions.yaml).

To benefit from new additions, you need to update the distribution list from
time to time.

This list is usually located in your home directory under
`$XDG_CONFIG/distributions.yaml` or `~/.config/binenv/distribution.yaml`).

To update only distributions:

```bash
binenv update --distributions # or -d
```

To update distributions **and** their versions:

```bash
binenv update --all # or -a
```

#### Examples

- `binenv update`: update available versions for all distributions
- `binenv update -d`: update available distributions
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

In the output, versions printed in reverse mode are the currenly selected
(a.k.a. active) versions (see [Selecting versions](#selecting-versions) below.

Versions in **bold** are installed.

All other versions are available to be installed.

#### Examples

```
$ binenv versions
terraform: 0.13.1 (/home/leucos/dev/perso/devops.works/github-repos/binenv) 0.13.0 0.13.0-rc1 0.13.0-beta3 0.13.0-beta2 0.13.0-beta1 0.12.29 0.12.28 0.12.27 0.12.26 0.12.25 0.12.24 0.12.23 0.12.22 0.12.21 0.12.20 0.12.19 0.12.18 0.12.17 0.12.16 0.12.15 0.12.14 0.12.13 0.12.12 0.12.11 0.12.10 0.12.9 0.12.8 0.12.7 0.12.6 
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

### Install versions form .binenv.lock

Install versions sepcified in `.binenv.lock` file, you can use the `--lock`
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

## Removing binenv stuff

`binenv` stores downloaded binaries in `~/.binenv/binaries`, and a cache in
`~/.config/binenv/` (or whatever your `XDG_CONFIG` variable points to).

To wipe everything clean:

```bash
rm -rfi ~/.binenv ~/.config/binenv/
```

Don't forget to remove the `PATH` and the completion you might have changed in
your shell rc file.

## Status

This is really _super alpha_ and has only be tested on Linux & MacOS. YMMV on
other platforms.

There are **no tests**. I will probably go to hell for this.

## Contributions

Welcomed !

We will need other installation mechanisms (see
https://github.com/devops-works/binenv/tree/master/internal/install).

## Licence

MIT
