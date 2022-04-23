# System-wide installation

**This feature is currently in alpha state.**

`binenv` can be installed system-wide, so binaries are managed for the entire
system.

In this mode, `binenv` is installed in `/opt/binenv`, followinf the
[FHS](https://en.wikipedia.org/wiki/Filesystem_Hierarchy_Standard).

Users in the system will be able to use any binary installed this way. However,
only `root` will be able to manage binaries.

To use `binenv` in this global mode, the `-g` flag is required before each command, e.g.:

```bash
sudo binenv -g install terraform
```

To run in global mode, you can either install from [scratch](#installation), or
[move](#migrate-user-installation-to-system-wide) an existing user
installation.

Note that operations that mutate state (i.e. `install`, `uninstall`, `update`,
`upgrade`) will have to be executed as `root`.

## Installation

Note: instructions for BSD systems have not been tested.

### Linux (bash/zsh)

```
wget -q https://github.com/devops-works/binenv/releases/latest/download/binenv_linux_amd64
wget -q https://github.com/devops-works/binenv/releases/latest/download/checksums.txt
sha256sum  --check --ignore-missing checksums.txt
chmod +x ./binenv_linux_amd64
sudo ./binenv_linux_amd64 -g update
sudo ./binenv_linux_amd64 -g install binenv
rm ./binenv_linux_amd64 
if [[ -n $BASH ]]; then ZESHELL=bash; fi
if [[ -n $ZSH_NAME ]]; then ZESHELL=zsh; fi
echo $ZESHELL
echo -e '\nexport PATH=/opt/binenv/:$PATH' >> ~/.${ZESHELL}rc
echo -e '\nexport PATH=/opt/binenv/:$PATH' | sudo tee -a /root/.${ZESHELL}rc
echo "source <(binenv completion ${ZESHELL})" >> ~/.${ZESHELL}rc
exec $SHELL
```

### MacOS (with bash)

```
wget -q https://github.com/devops-works/binenv/releases/latest/download/binenv_darwin_amd64
wget -q https://github.com/devops-works/binenv/releases/latest/download/checksums.txt
sha256sum  --check --ignore-missing checksums.txt
chmod +x binenv_darwin_amd64
sudo ./binenv_darwin_amd64 -g update
sudo ./binenv_darwin_amd64 -g install binenv
rm ./binenv_darwin_amd64 
echo -e '\nexport PATH=/opt/binenv/:$PATH' >> ~/.bashrc
echo -e '\nexport PATH=/opt/binenv/:$PATH' | sudo tee -a /root/.${ZESHELL}rc
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
chmod +x binenv_freebsd_amd64
sudo ./binenv_freebsd_amd64 -g update
sudo ./binenv_freebsd_amd64 -g install binenv
rm ./binenv_freebsd_amd64 
if [[ -n $BASH ]]; then ZESHELL=bash; fi
if [[ -n $ZSH_NAME ]]; then ZESHELL=zsh; fi
echo $ZESHELL
echo -e '\nexport PATH=/opt/binenv/:$PATH' >> ~/.${ZESHELL}rc
echo -e '\nexport PATH=/opt/binenv/:$PATH' | sudo tee -a /root/.${ZESHELL}rc
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
chmod +x binenv_openbsd_amd64
sudo ./binenv_openbsd_amd64 -g update
sudo ./binenv_openbsd_amd64 -g install binenv
rm ./binenv_openbsd_amd64 
if [[ -n $BASH ]]; then ZESHELL=bash; fi
if [[ -n $ZSH_NAME ]]; then ZESHELL=zsh; fi
echo $ZESHELL
echo -e '\nexport PATH=/opt/binenv/:$PATH' >> ~/.${ZESHELL}rc
echo -e '\nexport PATH=/opt/binenv/:$PATH' | sudo tee -a /root/.${ZESHELL}rc
echo "source <(binenv completion ${ZESHELL})" >> ~/.${ZESHELL}rc
exec $SHELL
```

If you are using a different shell, skip adding completion to your `.${SHELL}rc` file.

## Migrate user installation to system wide

This is not really recommended. But if you have bandwidth constraints and do
not want to re-download everything, you can do it like so:

- replace `~/.binenv` with `/opt/binenv` in your shell rcfile (`~/.bashrc` or
  `~/.zshrc`)
- create `/opt/binenv` with proper permissions

```bash
sudo mkdir -m 755 /opt/binenv
```

- copy cache and config

```bash
sudo cp -aR ~/.config/binenv/ /opt/binenv/config
sudo cp -aR ~/.cache/binenv/ /opt/binenv/cache
```

- copy binaries in `/opt/binenv`

```bash
sudo cp -aR ~/.binenv/* /opt/binenv/
```

- recreate symlinks

```bash
for i in $(find /opt/binenv/ -type l); do 
    sudo ln -sf /opt/binenv/shim $i
done
```

## Caveats

On Debian systems, PATH is reset by sudo. So calling `sudo binenv` (or other
installed binaries with `sudo`) will not work. To fix this, you have 3 options
(spoiler: none of them is perfect):

### change `secure_path` in `/etc/sudoers`

Replace the line:

`Defaults	secure_path="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin"`

with:

`Defaults	secure_path="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin:/opt/binenv"`

### pass env variables when calling sudo

```bash
sudo env "PATH=$PATH binenv" -g install foo
```

You can also do something like `alias sudo='sudo env "PATH=$PATH binenv"'` if
you're lazy, but this is not really recommended.

### use the full path

```bash
sudo /opt/binenv/binenv -g install yq
```

## FAQ

### Why not install in `/usr/local/bin` ?

`binenv` makes quite a mess and is quite opinionated in the way it handles
things. Having the bin directory in `/usr/local/bin/` is a recipe for a
disaster.

### Why do I have to use `-g` ? Can't binenv autodetect global mode ?

Yes it could. But what if you want to use a user installation along a global
one ? Say, for instance, when your sysadmin does not want to install a specific
binary.

If you point to `/opt/binenv/binenv`, you will have to use separated `binenv`
binaries to manage your own install. Probably a PITA.

The best way to achieve this is to `export BINENV_GLOBAL=true` in your shell
rc-file so binenv will always run in global mode.