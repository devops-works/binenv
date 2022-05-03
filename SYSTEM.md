# System-wide installation

**This feature is currently in alpha state.**

`binenv` can be installed system-wide, so binaries are managed for the entire
system.

In this mode, `binenv` is installed in several system-wide places:

- `cachedir`, holding the distributions cache, in `/var/cache/binenv`
- `confdir`, holding the distributions list, in `/var/lib/binenv`
- `linkdir`, holding symlinks to the shim, in `/usr/local/bin/`
- `bindir`, holding shim and distributions binaries

Users in the system will be able to use any binary installed this way. However,
only `root` will be able to manage binaries.

To use `binenv` in this global mode, the `-g` flag is required before each
command, e.g.:

```bash
sudo binenv -g install terraform
```

To run in global mode, you can either install from [scratch](#installation), or
[add](#add-system-wide-installation) a system wide installation.

Note that operations that mutate state (i.e. `install`, `uninstall`, `update`,
`upgrade`) will have to be executed as `root`.

**Note that it is mandatory to set `BINENV_GLOBAL=true` to use binaries
installed in system-wide mode (see [Caveat](#caveat)).**

## Installation

Note: instructions for BSD systems have not been tested.

### Linux (bash/zsh)

```
wget -q https://github.com/devops-works/binenv/releases/latest/download/binenv_linux_amd64
wget -q https://github.com/devops-works/binenv/releases/latest/download/checksums.txt
sha256sum  --check --ignore-missing checksums.txt
mv binenv_linux_amd64 binenv
chmod +x ./binenv
sudo ./binenv -g update
sudo ./binenv -g install binenv 0.19.0
rm ./binenv 
if [[ -n $BASH ]]; then ZESHELL=bash; fi
if [[ -n $ZSH_NAME ]]; then ZESHELL=zsh; fi
echo $ZESHELL
echo "source <(binenv completion ${ZESHELL})" >> ~/.${ZESHELL}rc
exec $SHELL
```

### MacOS (with bash)

```
wget -q https://github.com/devops-works/binenv/releases/latest/download/binenv_darwin_amd64
wget -q https://github.com/devops-works/binenv/releases/latest/download/checksums.txt
sha256sum  --check --ignore-missing checksums.txt
mv binenv_darwin_amd64 binenv
chmod +x binenv
sudo ./binenv -g update
sudo ./binenv -g install binenv 0.19.0
rm ./binenv 
echo 'source <(binenv completion bash)' >> ~/.bashrc
exec $SHELL
```

### Windows

binenv does not support windows.

### FreeBSD (bash/zsh)

```
fetch https://github.com/devops-works/binenv/releases/latest/download/binenv_freebsd_amd64
fetch https://github.com/devops-works/binenv/releases/latest/download/checksums.txt
shasum --ignore-missing -a 512 -c checksums.txt
mv binenv_freebsd_amd64 binenv
chmod +x binenv
sudo ./binenv -g update
sudo ./binenv -g install binenv 0.19.0
rm ./binenv 
if [[ -n $BASH ]]; then ZESHELL=bash; fi
if [[ -n $ZSH_NAME ]]; then ZESHELL=zsh; fi
echo $ZESHELL
echo "source <(binenv completion ${ZESHELL})" >> ~/.${ZESHELL}rc
exec $SHELL
```

If you are using a different shell, skip adding completion to your `.${SHELL}rc` file.

To be able to verify checksums, you have to install the `p5-Digest-SHA` package.

### OpenBSD (bash/zsh)

```
ftp https://github.com/devops-works/binenv/releases/latest/download/binenv_openbsd_amd64
ftp https://github.com/devops-works/binenv/releases/latest/download/checksums.txt
cksum -a sha256 -C checksums.txt binenv_openbsd_amd64
mv binenv_openbsd_amd64 binenv
chmod +x binenv
sudo ./binenv -g update
sudo ./binenv -g install binenv 0.19.0
rm ./binenv 
if [[ -n $BASH ]]; then ZESHELL=bash; fi
if [[ -n $ZSH_NAME ]]; then ZESHELL=zsh; fi
echo $ZESHELL
echo "source <(binenv completion ${ZESHELL})" >> ~/.${ZESHELL}rc
exec $SHELL
```

If you are using a different shell, skip adding completion to your `.${SHELL}rc` file.

## Add system-wide installation

- Ensure you are using at least `binenv` version v0.19.0.

```bash
binenv version
```

- install `binenv` system-wide

```bash
binenv -g install binenv 0.19.0-rc4
```

- adjust path in your rcfiles according to your preferences

## Caveat

To use system-wide installed binaries, you have to set the `BINENV_GLOBAL=true`
environment variable.

This is required since binenv shim will look in user paths by default if this
is not set.

## FAQ

### Why do I have to use `-g` ? Can't binenv autodetect global mode ?

Yes it could. But what if you want to use a user installation along a global
one ? Say, for instance, when your sysadmin does not want to install a specific
binary.

The best way to achieve this is to `export BINENV_GLOBAL=true` in your shell
rc-file so binenv will always run in global mode.
